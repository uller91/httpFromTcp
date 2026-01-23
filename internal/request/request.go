package request

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/uller91/httpFromTcp/internal/headers"
	"io"
	"strconv"
	"strings"
)

type Request struct {
	RequestLine  RequestLine
	RequestState RequestState
	Headers      headers.Headers
	Body         []byte
	readBody     int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const crlf = "\r\n"

const bufferSize = 8

type RequestState string

const (
	RequestStateInitialized    RequestState = "initialized"
	RequestStateDone           RequestState = "done"
	RequestStateParsingHeaders RequestState = "parsing headers"
	RequestStateParsingBody    RequestState = "parsing body"
)

func (r *Request) parse(data []byte) (int, error) {
	totalSoFar := 0
	for r.RequestState != RequestStateDone {
		read, err := r.parseOne(data[totalSoFar:])
		if err != nil {
			return 0, err
		}
		totalSoFar += read
		if read == 0 {
			break
		}
	}
	return totalSoFar, nil
}

func (r *Request) parseOne(data []byte) (int, error) {
	switch r.RequestState {
	case RequestStateInitialized:
		reqLine, read, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if read == 0 {
			return read, nil
		}
		r.RequestLine = reqLine
		r.RequestState = RequestStateParsingHeaders
		return read, nil
	case RequestStateParsingHeaders:
		read, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.RequestState = RequestStateParsingBody
			return read, nil
		}
		return read, nil
	case RequestStateParsingBody:
		value, exist := r.Headers.Get("Content-Length")
		if !exist {
			r.RequestState = RequestStateDone
			return 0, nil
		}
		length, err := strconv.Atoi(value)
		if err != nil {
			return 0, errors.New("error: non-int valuse in Content-Length header")
		}
		r.Body = append(r.Body, data...)
		r.readBody += len(data)
		if r.readBody > length {
			return len(data), errors.New("error: body is greater than reported by Content-Length header; reported")
		} else if r.readBody == length {
			r.RequestState = RequestStateDone
		}
		return len(data), nil

	case RequestStateDone:
		return 0, errors.New("error: trying to read data in a done state")
	default:
		return 0, errors.New("error: unknown state")
	}
}

func parseRequestLine(data []byte) (RequestLine, int, error) {
	if !bytes.Contains(data, []byte(crlf)) {
		return RequestLine{}, 0, nil
	}

	idx := bytes.Index(data, []byte(crlf))
	line := string(data[:idx])
	if line == "" {
		return RequestLine{}, 0, errors.New("request line is empty")
	}

	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return RequestLine{}, 0, errors.New("malformed request line: not enough parts")
	}

	if parts[0] != strings.ToUpper(parts[0]) {
		return RequestLine{}, 0, errors.New("improper method in the request line")
	}

	if !strings.Contains(parts[2], "HTTP") {
		return RequestLine{}, 0, errors.New("melformed request line: shuffled parts")
	}

	version := strings.Split(parts[2], "/")
	if len(version) != 2 {
		return RequestLine{}, 0, errors.New("malformed http version in therequest line")
	}
	if version[1] != "1.1" {
		return RequestLine{}, 0, errors.New("improper http version in therequest line")
	}

	reqLine := RequestLine{
		HttpVersion:   version[1],
		RequestTarget: parts[1],
		Method:        parts[0],
	}
	return reqLine, idx + 2, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize, bufferSize)
	readToIndex := 0
	req := Request{RequestState: RequestStateInitialized, Headers: headers.NewHeaders(), Body: make([]byte, 0)}

	for req.RequestState != RequestStateDone {
		if readToIndex == cap(buf) {
			buf_new := make([]byte, cap(buf)*2, cap(buf)*2)
			_ = copy(buf_new, buf)
			buf = buf_new
		}

		read, err := reader.Read(buf[readToIndex:])
		if err == io.EOF {
			if req.RequestState != RequestStateDone {
				fmt.Println(req.RequestState)
				return nil, errors.New("incomplete or malformed request")
			}
			break
		}
		if err != nil {
			return nil, err
		}
		readToIndex += read

		read, err = req.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}
		//buf_new := make([]byte, cap(buf), cap(buf))
		//_ = copy(buf_new, buf[read:])
		//buf = buf_new
		copy(buf, buf[read:])
		readToIndex -= read
	}

	return &req, nil
}

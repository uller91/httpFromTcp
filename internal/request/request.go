package request

import (
	"bytes"
	"errors"
	"io"
	"strings"
)

type Request struct {
	RequestLine  RequestLine
	RequestState RequestState
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
	Initialized RequestState = "initialized"
	Done        RequestState = "done"
)

func (r *Request) parse(data []byte) (int, error) {
	switch r.RequestState {
	case Initialized:
		reqLine, read, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if read == 0 {
			return 0, nil
		}
		r.RequestLine = reqLine
		r.RequestState = Done
		return read, nil
	case Done:
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
	return reqLine, idx, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize, bufferSize)
	readToIndex := 0
	req := Request{RequestState: Initialized}

	for req.RequestState != Done {
		if readToIndex == cap(buf) {
			buf_new := make([]byte, cap(buf)*2, cap(buf)*2)
			_ = copy(buf_new, buf)
			buf = buf_new
		}

		read, err := reader.Read(buf[readToIndex:])
		if err == io.EOF {
			req.RequestState = Done
			break
		}
		if err != nil {
			return nil, err
		}
		readToIndex += read

		read, err = req.parse(buf)
		if err != nil {
			return nil, err
		}
		buf_new := make([]byte, cap(buf), cap(buf))
		_ = copy(buf_new, buf[read:])
		buf = buf_new
		readToIndex -= read
	}

	return &req, nil
}

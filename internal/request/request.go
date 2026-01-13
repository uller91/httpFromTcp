package request

import (
	"bytes"
	"errors"
	"fmt"
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

type RequestState int

const (
	Initialized RequestState = iota
	Done
)

func parseRequestLine(bytes []byte) (RequestLine, int, error) {
	if !bytes.Contains(data, []byte(crlf)) {
		return RequestLine{}, 0, nil
	}

	idx := bytes.Index(data, []byte(crlf))
	line := string(data[:idx])
	if line == "" {
		return nil, 0, errors.New("request line is empty")
	}

	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return RequestLine{}, 0, errors.New("malformed request line: not enough parts")
	}

	if parts[0] != strings.ToUpper(parts[0]) {
		return RequestLine{}, 0, errors.New("improper method in the request line")
	}

	if !strings.Contains(parts[2], "HTTP") {
		return RequestLine{}, errors.New("melformed request line: shuffled parts")
	}

	version := strings.Split(parts[2], "/")
	if len(version) != 2 {
		return RequestLine{}, errors.New("malformed http version in therequest line")
	}
	if version[1] != "1.1" {
		return RequestLine{}, errors.New("improper http version in therequest line")
	}

	reqLine := RequestLine{
		HttpVersion:   version[1],
		RequestTarget: parts[1],
		Method:        parts[0],
	}
	return reqLine, idx, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	bytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	/*
		reqStr := string(bytes)
		reqStrs := strings.Split(reqStr, "\r\n")
		fmt.Println(reqStrs[0])
		if reqStrs[0] == "" {
			return nil, errors.New("requesti line is empty")
		}
	*/

	reqLine, read, err := parseRequestLine(bytes)
	if err != nil {
		return nil, err
	}
	req := Request{RequestLine: reqLine}

	return &req, nil
}

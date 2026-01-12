package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func parseRequestLine(line string) (RequestLine, error) {
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return RequestLine{}, errors.New("malformed request line: not enough parts")
	}

	if parts[0] != strings.ToUpper(parts[0]) {
		return RequestLine{}, errors.New("improper method in the request line")
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
	return reqLine, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	bts, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	reqStr := string(bts)
	reqStrs := strings.Split(reqStr, "\r\n")
	fmt.Println(reqStrs[0])
	if reqStrs[0] == "" {
		return nil, errors.New("requesti line is empty")
	}

	reqLine, err := parseRequestLine(reqStrs[0])
	if err != nil {
		return nil, err
	}
	req := Request{RequestLine: reqLine}

	return &req, nil
}

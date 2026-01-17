package headers

import (
	"bytes"
	"strings"
	"errors"
)

const crlf = "\r\n"

type Headers map[string]string

func NewHeaders() Headers {
	return map[string]string{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	if !bytes.Contains(data, []byte(crlf)) {
		return 0, false, nil
	}

	idx := bytes.Index(data, []byte(crlf))
	if idx == 0 {
		return 0, true, nil
	}
	line := string(data[:idx])

	lineClean := strings.TrimSpace(line)

	parts := strings.SplitN(lineClean, ":", 2)

	if parts[0] != strings.TrimSpace(parts[0]) {
		return 0, false, errors.New("malformed header string")
	}
	
	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	h[key] = value

	return idx+2, false, nil //need to remove \r\n as well
}

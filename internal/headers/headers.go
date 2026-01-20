package headers

import (
	"bytes"
	"errors"
	"strings"
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
	key := strings.TrimSpace(parts[0])

	if parts[0] != key {
		return 0, false, errors.New("malformed header string")
	}

	if !validKey(key) {
		return 0, false, errors.New("invalid character in the key")
	}
	/*
		if strings.Contains(parts[0], "@") {
			return 0, false, errors.New("invalid character in the key")
		} */

	h.Set(key, parts[1])

	return idx + 2, false, nil //need to remove \r\n as well
}

func (h Headers) Set(k, v string) {
	key := strings.ToLower(k)
	value := strings.TrimSpace(v)
	oldValue, exist := h[key]
	if exist {
		value = oldValue + ", " + value
	}
	h[key] = value
}

func validKey(key string) bool {
	for _, c := range key {
		if !isCorrect(string(c)) {
			return false
		}
	}
	return true
}

var specialChars = []string{"!", "#", "$", "%", "&", "'", "*", "+", "-", ".", "^", "_", "`", "|", "~"}

func isCorrect(c string) bool {
	if c >= "A" && c <= "Z" ||
		c >= "a" && c <= "z" ||
		c >= "0" && c <= "9" {
		return true
	}

	for _, sp := range specialChars {
		if sp == c {
			return true
		}
	}

	return false
}

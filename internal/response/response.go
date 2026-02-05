package response

import (
	"github.com/uller91/httpFromTcp/internal/headers"
	"io"
	"strconv"
)

type StatusCode string

const (
	OK                  StatusCode = "200"
	BadRequest          StatusCode = "400"
	InternalServerError StatusCode = "500"
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	var statusLine string
	switch statusCode {
	case OK:
		statusLine = "HTTP/1.1 200 OK"
	case BadRequest:
		statusLine = "HTTP/1.1 400 Bad Request"
	case InternalServerError:
		statusLine = "HTTP/1.1 500 Internal Server Error"
	default:
		statusLine = "HTTP/1.1 " + string(statusCode) + " "
	}
	statusLine += "\r\n"
	_, err := w.Write([]byte(statusLine))
	if err != nil {
		return err
	}

	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()
	headers["Content-Length"] = strconv.Itoa(contentLen)
	headers["Connection"] = "close"
	headers["Content-Type"] = "text/plain"

	return headers
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	var headersLines string
	for key, value := range headers {
		headersLines += key + ": " + value + "\r\n"
	}
	headersLines += "\r\n"

	_, err := w.Write([]byte(headersLines))
	if err != nil {
		return err
	}

	return nil
}
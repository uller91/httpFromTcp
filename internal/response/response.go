package response

import (
	"errors"
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

type Writer struct {
	Writer          io.Writer
	OptionalHeaders map[string]string
	WriterState     WriterState
}

type WriterState string

const (
	Initialized        WriterState = "initialized"
	Done               WriterState = "done"
	FinishedStatusLine WriterState = "finished status line"
	FinishedHeaders    WriterState = "finished headers"
)

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.WriterState == Initialized {
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
		_, err := w.Writer.Write([]byte(statusLine))
		if err != nil {
			w.WriterState = Done
			return err
		}
		w.WriterState = FinishedStatusLine

		return nil
	} else {
		return errors.New("Writer is not initialized!")
	}
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()
	headers["Content-Length"] = strconv.Itoa(contentLen)
	headers["Connection"] = "close"
	headers["Content-Type"] = "text/plain"

	return headers
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.WriterState == FinishedStatusLine {
		var headersLines string
		for key, value := range headers {
			headersLines += key + ": " + value + "\r\n"
		}
		headersLines += "\r\n"

		_, err := w.Writer.Write([]byte(headersLines))
		if err != nil {
			w.WriterState = Done
			return err
		}
		w.WriterState = FinishedHeaders

		return nil
	} else {
		return errors.New("Status line hasn't been sent!")
	}
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.WriterState == FinishedHeaders {
		n, err := w.Writer.Write(p)
		if err != nil {
			w.WriterState = Done
			return 0, err
		}
		w.WriterState = Done

		return n, nil
	} else {
		return 0, errors.New("Headers hasn't been sent!")
	}
}

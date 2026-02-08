package response

import (
	"errors"
	"fmt"
	"github.com/uller91/httpFromTcp/internal/headers"
	"io"
	"strconv"
)

const crlf = "\r\n"

type StatusCode string

const (
	OK                  StatusCode = "200"
	BadRequest          StatusCode = "400"
	InternalServerError StatusCode = "500"
)

type Writer struct {
	Writer      io.Writer
	WriterState WriterState
}

type WriterState string

const (
	Initialized    WriterState = "initialized"
	Done           WriterState = "done"
	WritingHeaders WriterState = "writing headers"
	WritingBody    WriterState = "writing body"
)

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.WriterState != WritingBody {
		return 0, errors.New("Headers hasn't been sent!")
	}

	n := 0

	hexLenP := fmt.Sprintf("%x\r\n", len(p))
	write, err := w.Writer.Write([]byte(hexLenP))
	if err != nil {
		return 0, err
	}
	n += write

	write, err = w.Writer.Write(p)
	if err != nil {
		return 0, err
	}
	n += write

	write, err = w.Writer.Write([]byte(crlf))
	if err != nil {
		return 0, err
	}
	n += write

	return n, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.WriterState != WritingBody {
		return 0, errors.New("Headers hasn't been sent!")
	}

	payload := []byte(fmt.Sprintf("%x\r\n\r\n", 0))

	n, err := w.Writer.Write(payload)
	if err != nil {
		return 0, err
	}
	w.WriterState = Done

	return n, nil
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.WriterState != Initialized {
		return errors.New("Writer is not initialized!")
	}

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
		return err
	}
	w.WriterState = WritingHeaders

	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()
	headers.Set("Content-Length", strconv.Itoa(contentLen))
	headers.Set("Connection", "close")
	headers.Set("Content-Type", "text/plain")

	return headers
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.WriterState != WritingHeaders {
		return errors.New("Status line hasn't been sent!")
	}

	var headersLines string
	for key, value := range headers {
		headersLines += key + ": " + value + "\r\n"
	}
	headersLines += "\r\n"

	_, err := w.Writer.Write([]byte(headersLines))
	if err != nil {
		return err
	}
	w.WriterState = WritingBody

	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.WriterState != WritingBody {
		return 0, errors.New("Headers hasn't been sent!")
	}

	n, err := w.Writer.Write(p)
	if err != nil {
		return 0, err
	}
	w.WriterState = Done

	return n, nil
}

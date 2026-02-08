package server

import (
	//"bytes"
	"fmt"
	"github.com/uller91/httpFromTcp/internal/request"
	"github.com/uller91/httpFromTcp/internal/response"
	//"io"
	"net"
	"sync/atomic"
)

type Server struct {
	Listener net.Listener
	Handler  Handler
	Running  *atomic.Bool
}

type Handler func(w response.Writer, req *request.Request) //*HandlerError

/*
type HandlerError struct {
	StatusCode   response.StatusCode
	ErrorMessage string
}
*/

/*
func HandleError(w io.Writer, handlErr *HandlerError) {
	err := response.WriteStatusLine(w, handlErr.StatusCode)
	if err != nil {
		fmt.Errorf("Error sending status line: %v", err)
		return
	}

	headers := response.GetDefaultHeaders(len([]byte(handlErr.ErrorMessage)))

	err = response.WriteHeaders(w, headers)
	if err != nil {
		fmt.Errorf("Error sending headers: %v", err)
		return
	}

	//fmt.Println(string([]byte(handlErr.ErrorMessage)))

	_, err = w.Write([]byte(handlErr.ErrorMessage))
	if err != nil {
		fmt.Errorf("Error sending the body: %v", err)
		return
	}

	return
}
*/

func Serve(port int, handler Handler) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("Error when creating TCP listener: %v", err)
	}

	var running atomic.Bool
	running.Store(true)
	s := Server{Listener: l, Handler: handler, Running: &running}
	go s.listen()

	return &s, nil
}

func (s *Server) Close() error {
	err := s.Listener.Close()
	if err != nil {
		return fmt.Errorf("Error when closing TCP listener: %v", err)
	}
	s.Running.Store(false)

	return nil
}

func (s *Server) listen() {
	for {
		// Wait for a connection.
		conn, err := s.Listener.Accept()
		if err != nil {
			running := s.Running.Load()
			if running {
				fmt.Errorf("Error creating connectionl: %v", err)
				continue
			} else {
				return
			}
		}
		fmt.Println("TCP connection has been accepted!")

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Errorf("Error processing the request: %v", err)
		return
	}

	//b := bytes.NewBuffer([]byte{})

	resW := response.Writer{
		Writer:      conn,
		WriterState: response.Initialized,
	}

	s.Handler(resW, req)

	/*
		body := b.Bytes()
		//fmt.Println(string(body))

		//headers := response.GetDefaultHeaders(len(body))

		err = response.WriteStatusLine(conn, response.OK)
		if err != nil {
			fmt.Errorf("Error sending status line: %v", err)
			return
		}

		err = response.WriteHeaders(conn, headers)
		if err != nil {
			fmt.Errorf("Error sending headers: %v", err)
			return
		}

		_, err = conn.Write(body)
		if err != nil {
			fmt.Errorf("Error sending the body: %v", err)
			return
		}
	*/

	return
}

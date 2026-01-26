package server

import (
	"sync/atomic"
	"fmt"
	"net"
	"strconv"
)

type Server struct {
	Listener net.Listener
	Running   *atomic.Bool
}

func Serve(port int) (*Server, error) {
	strPort := strconv.Itoa(port)
	adr := ":" + strPort
	l, err := net.Listen("tcp", adr)
	if err != nil {
		return nil, fmt.Errorf("Error when creating TCP listener: %v", err)
	}

	var running atomic.Bool
	running.Store(true)
	s := Server{Listener: l, Running: &running}
	go s.listen()

	return &s, nil
}

func (s *Server) Close() error {
	err := s.Listener.Close()
	if err != nil {
		fmt.Errorf("Error when closing TCP listener: %v", err)
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
			} else {
				return
			}
		}
		fmt.Println("TCP connection has been accepted!")

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	response := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n" +
		"Content-Length: 13\r\n" +
		"\r\n" +
		"Hello World!\n"

	_, err := conn.Write([]byte(response))
	if err != nil {
		fmt.Errorf("Error sending the response: %v", err)
	}

	err = conn.Close()
	if err != nil {
		fmt.Errorf("Error when closing TCP connection: %v", err)
	}
}

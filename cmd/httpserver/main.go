package main

import (
	"github.com/uller91/httpFromTcp/internal/request"
	"github.com/uller91/httpFromTcp/internal/response"
	"github.com/uller91/httpFromTcp/internal/server"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"fmt"
)

const port = 42069

func main() {
	server, err := server.Serve(port, MainHandler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func MainHandler(w io.Writer, req *request.Request) {
	resW := response.Writer{
		Writer:          w,
		OptionalHeaders: map[string]string{"Content-Type": "text/html"},
		WriterState:     response.Initialized,
	}

	var statusCode response.StatusCode
	var body []byte

	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		statusCode = response.BadRequest
		body = []byte(`<html>
  			<head>
   				<title>400 Bad Request</title>
  			</head>
			<body>
				<h1>Bad Request</h1>
				<p>Your request honestly kinda sucked.</p>
			</body>
			</html>`)
	case "/myproblem":
		statusCode = response.InternalServerError
		body = []byte(`<html>
  			<head>
   				<title>500 Internal Server Error</title>
  			</head>
			<body>
				<h1>Internal Server Error</h1>
				<p>Okay, you know what? This one is on me.</p>
			</body>
			</html>`)
	default:
		statusCode = response.OK
		body = []byte(`<html>
  			<head>
   				<title>200 OK</title>
  			</head>
			<body>
				<h1>Success!</h1>
				<p>Your request was an absolute banger.</p>
			</body>
			</html>`)
	}

	headers := response.GetDefaultHeaders(len(body))
	headers["Content-Type"] = "text/html"

	err := resW.WriteStatusLine(statusCode)
	if err != nil {
		fmt.Errorf("Error sending status line: %v", err)
		return
	}

	err = resW.WriteHeaders(headers)
	if err != nil {
		fmt.Errorf("Error sending headers: %v", err)
		return
	}
	_, err = resW.WriteBody(body)
	if err != nil {
		fmt.Errorf("Error sending the body: %v", err)
		return
	}

	return
}

package main

import (
	"crypto/sha256"
	"fmt"
	"github.com/uller91/httpFromTcp/internal/request"
	"github.com/uller91/httpFromTcp/internal/response"
	"github.com/uller91/httpFromTcp/internal/server"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const port = 42069

func main() {
	server, err := server.Serve(port, Handler)
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

func VideoHandler(w response.Writer, req *request.Request) {
	fmt.Println("Serving the video from assets...")

	statusCode := response.OK

	videoBody, err := os.ReadFile("assets/vim.mp4")
	if err != nil {
		fmt.Errorf("Error reading video file: %v", err)
		return
	}

	headers := response.GetDefaultHeaders(len(videoBody))
	headers.Change("Content-Type", "video/mp4")

	err = w.WriteStatusLine(statusCode)
	if err != nil {
		fmt.Errorf("Error sending status line: %v", err)
		return
	}

	err = w.WriteHeaders(headers)
	if err != nil {
		fmt.Errorf("Error sending headers: %v", err)
		return
	}
	_, err = w.WriteBody(videoBody)
	if err != nil {
		fmt.Errorf("Error sending the body: %v", err)
		return
	}

	return
}

func ProxyHandler(w response.Writer, req *request.Request) {
	target := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
	targetUrl := "https://httpbin.org/" + target

	fmt.Printf("Servers is proxying to %s\n", targetUrl)

	res, err := http.Get(targetUrl)
	if err != nil {
		handler500(w, req)
		return
	}
	defer res.Body.Close()

	statusCode := response.OK

	headers := response.GetDefaultHeaders(0)
	headers.Change("Transfer-Encoding", "chunked")
	headers.Set("Trailer", "X-Content-Sha256")
	headers.Set("Trailer", "X-Content-Length")
	headers.Delete("Content-Length")

	err = w.WriteStatusLine(statusCode)
	if err != nil {
		fmt.Errorf("Error sending status line: %v", err)
		return
	}

	err = w.WriteHeaders(headers)
	if err != nil {
		fmt.Errorf("Error sending headers: %v", err)
		return
	}

	const bufSize = 1024
	buf := make([]byte, bufSize)

	chunkedBodyLen := 0
	var resBody []byte
	for {
		read, err := res.Body.Read(buf)
		fmt.Printf("Read %v bytes\n", read)

		if read > 0 {
			resBody = append(resBody, buf[:read]...)
			_, err := w.WriteChunkedBody(buf[:read])
			if err != nil {
				fmt.Errorf("Error writing chunked body done: %v", err)
				break
			}
		}
		chunkedBodyLen += read

		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Errorf("Error reading response body: %v", err)
			break
		}
	}

	_, err = w.WriteChunkedBodyDone()
	if err != nil {
		fmt.Errorf("Error writing chunked body done: %v", err)
		return
	}

	sha256 := sha256.Sum256(resBody)
	trailers := response.GetDefaultTrailers(fmt.Sprintf("%x", sha256), chunkedBodyLen)

	fmt.Printf("%#v\n", trailers)
	err = w.WriteTrailers(trailers)
	if err != nil {
		fmt.Errorf("Error sending trailers: %v", err)
		return
	}
}

func handler200(w response.Writer, req *request.Request) {
	statusCode := response.OK

	body := []byte(`<html>
		<head>
			<title>200 OK</title>
		</head>
		<body>
			<h1>Success!</h1>
			<p>Your request was an absolute banger.</p>
		</body>
		</html>`)

	headers := response.GetDefaultHeaders(len(body))
	headers.Change("Content-Type", "text/html")

	err := w.WriteStatusLine(statusCode)
	if err != nil {
		fmt.Errorf("Error sending status line: %v", err)
		return
	}

	err = w.WriteHeaders(headers)
	if err != nil {
		fmt.Errorf("Error sending headers: %v", err)
		return
	}
	_, err = w.WriteBody(body)
	if err != nil {
		fmt.Errorf("Error sending the body: %v", err)
		return
	}

	return
}

func handler400(w response.Writer, req *request.Request) {
	statusCode := response.BadRequest

	body := []byte(`<html>
		<head>
			<title>400 Bad Request</title>
		</head>
		<body>
			<h1>Bad Request</h1>
			<p>Your request honestly kinda sucked.</p>
		</body>
		</html>`)

	headers := response.GetDefaultHeaders(len(body))
	headers.Change("Content-Type", "text/html")

	err := w.WriteStatusLine(statusCode)
	if err != nil {
		fmt.Errorf("Error sending status line: %v", err)
		return
	}

	err = w.WriteHeaders(headers)
	if err != nil {
		fmt.Errorf("Error sending headers: %v", err)
		return
	}
	_, err = w.WriteBody(body)
	if err != nil {
		fmt.Errorf("Error sending the body: %v", err)
		return
	}

	return
}

func handler500(w response.Writer, req *request.Request) {
	statusCode := response.InternalServerError

	body := []byte(`<html>
		<head>
			<title>500 Internal Server Error</title>
		</head>
		<body>
			<h1>Internal Server Error</h1>
			<p>Okay, you know what? This one is on me.</p>
		</body>
		</html>`)

	headers := response.GetDefaultHeaders(len(body))
	headers.Change("Content-Type", "text/html")

	err := w.WriteStatusLine(statusCode)
	if err != nil {
		fmt.Errorf("Error sending status line: %v", err)
		return
	}

	err = w.WriteHeaders(headers)
	if err != nil {
		fmt.Errorf("Error sending headers: %v", err)
		return
	}
	_, err = w.WriteBody(body)
	if err != nil {
		fmt.Errorf("Error sending the body: %v", err)
		return
	}

	return
}

func Handler(w response.Writer, req *request.Request) {
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
		ProxyHandler(w, req)
		return
	}

	if strings.HasPrefix(req.RequestLine.RequestTarget, "/video") {
		VideoHandler(w, req)
		return
	}

	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		handler400(w, req)
		return
	case "/myproblem":
		handler500(w, req)
		return
	default:
		handler200(w, req)
		return
	}

	return
}

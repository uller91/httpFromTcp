package main

import (
	"fmt"
	"github.com/uller91/httpFromTcp/internal/request"
	"log"
	"net"
)

const addressString = ":42069"

func main() {
	l, err := net.Listen("tcp", addressString)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("TCP connection has been accepted!")

		request, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", request.RequestLine.Method)
		fmt.Printf("- Target: %s\n", request.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", request.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		for key, value := range request.Headers {
			fmt.Printf("- %s: %s\n", key, value)
		}
		fmt.Println("Body:")
		fmt.Println(string(request.Body))

		fmt.Println("TCP connection has been closed!")
	}

}

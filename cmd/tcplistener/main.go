package main

import (
	"fmt"
	"io"
	"log"
	"strings"
	"net"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	line := ""
	go func() {
		defer close(ch)
		defer f.Close()
		for {
			buffer := make([]byte, 8)
			read, err := f.Read(buffer)
			if err == io.EOF {
				ch <- line
				break
			}
			part := string(buffer[:read])
			if strings.Contains(part, "\n") {
				parts := strings.Split(part, "\n")
				line += parts[0]
				ch <- line
				line = parts[1]
			} else {
				line += part
			}
		}
	}()

	return ch
}

func main() {
	/*
		file, err := os.Open("messages.txt")
		if err != nil {
			log.Fatal(err)
		}
	*/

	const addressString = ":42069"

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

		lines := getLinesChannel(conn)
		for line := range lines {
			fmt.Println(line)
		}

		fmt.Println("TCP connection has been closed!")
	}


}

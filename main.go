package main

import (
	"fmt"
	"os"
	"log"
	"io"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	line := ""
	go func() {
		defer close(ch)
		defer f.Close()
		for ;; {
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
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal(err)
	}
	
	lines := getLinesChannel(file)
	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}
	/*
	line := ""
	for ;; {
		buffer := make([]byte, 8)
		read, err := file.Read(buffer)
		if err == io.EOF {
			fmt.Printf("read: %s\n", line)
			os.Exit(0)
		}
		part := string(buffer[:read])
		if strings.Contains(part, "\n") {
			parts := strings.Split(part, "\n")
			line += parts[0]
			fmt.Printf("read: %s\n", line)
			line = parts[1]
		} else {
			line += part
		}
	} */

}
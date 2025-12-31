package main

import (
	"fmt"
	"os"
	"log"
	"io"
	"strings"
)

func main() {
	fmt.Println("I hope I get the job!")
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

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
	}

}
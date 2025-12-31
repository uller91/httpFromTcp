package main

import (
	"fmt"
	"os"
	"log"
	"io"
)

func main() {
	fmt.Println("I hope I get the job!")
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	buffer := make([]byte, 8)
	for ;; {
		read, err := file.Read(buffer)
		if err == io.EOF {
			//fmt.Println("Exiting")
			os.Exit(0)
		}
		str := string(buffer[:read])
		fmt.Printf("read: %s\n", str)
	}
}
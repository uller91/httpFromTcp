package main

import (
	"net"
	"bufio"
	"os"
	"log"
	"fmt"
)

func main() {
	const addressString = "localhost:42069"

	addr, err := net.ResolveUDPAddr("udp", addressString)
	if err != nil {
		log.Fatal(err)
	}
	
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for ;; {
		fmt.Printf("> ")
		str, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		_, err = conn.Write([]byte(str))
		if err != nil {
			log.Fatal(err)
		}
	}
}
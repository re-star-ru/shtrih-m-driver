package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	log.Println("Launching server...")

	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := ln.Accept()
	if err != nil {
		log.Fatal(err)
	}

	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("MessageRecieved:", string(message))
	}
}

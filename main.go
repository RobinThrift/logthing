package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

const (
	timeout = time.Second * 2
	port    = ":5960"
)

func main() {
	ln, err := net.Listen("tcp", port)
	defer ln.Close()

	if err != nil {
		fmt.Printf("Can't listen on tcp port %s\n", port)
		return
	}

	fmt.Printf("Started server on port %s\n", port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection")
			continue
		}
		conn.SetDeadline(time.Now().Add(timeout))
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	msg, err := bufio.NewReader(conn).ReadString('\n')

	if err != nil {
		fmt.Printf("Error receiving message")
		conn.Close()
		return
	}

	message := strings.TrimSpace(string(msg))

	fmt.Println("Message received: ", message)

	conn.Write([]byte(msg))
	conn.Close()
}

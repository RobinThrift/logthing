package main

import (
	"bufio"
	"net"
	"testing"
)


func TestHandleConnection(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	msg := "Hello World\n"

	go handleConnection(client)

	server.Write([]byte(msg))

	retMsgBytes, err := bufio.NewReader(server).ReadString('\n')
	retMsg := string(retMsgBytes)

	if err != nil {
		t.Fatalf(err.Error())
	}

	if msg != retMsg {
		t.Fatalf("Returned message was not the same. Got %s", retMsg)
	}
}

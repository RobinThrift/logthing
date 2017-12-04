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

	s := &service{
		Interval:  0,
		Buffer:    NewRingBuffer(1),
		Recipient: "test@example.com",
	}

	msg := "Hello World\n"

	go handleConnection(client, s)

	server.Write([]byte(msg))

	retMsgBytes, err := bufio.NewReader(server).ReadString('\n')
	retMsg := string(retMsgBytes)

	s.Buffer.Do(func(v interface{}) {
		value, ok := v.(string)
		if !ok {
			t.Fatalf("Value is not `string`. Got %T", v)
		}

		if value+"\n" != retMsg {
			t.Fatalf("Value is not '%s'. Got %s", retMsg, value)
		}
	})

	if err != nil {
		t.Fatalf(err.Error())
	}

	if msg != retMsg {
		t.Fatalf("Returned message was not the same. Got %s", retMsg)
	}
}

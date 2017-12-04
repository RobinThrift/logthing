package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	timeout = time.Second * 2
	port    = ":5960"
)

func main() {
	interval, err := time.ParseDuration(os.Getenv("LT_INTERVAL"))

	if err != nil {
		fmt.Println("Error parsing interval.")
		os.Exit(1)
	}

	bufferSize, err := strconv.Atoi(os.Getenv("LT_BUFFER"))

	if err != nil {
		fmt.Println("Error parsing buffer size.")
		os.Exit(1)
	}

	s := &service{
		Interval:  interval,
		Buffer:    NewRingBuffer(bufferSize),
		Recipient: os.Getenv("LT_RECIPIENT"),
		Sender:    os.Getenv("LT_SENDER"),
	}

	ln, err := net.Listen("tcp", port)
	defer ln.Close()

	if err != nil {
		fmt.Printf("Can't listen on tcp port %s\n", port)
		return
	}

	fmt.Printf("Started server on port %s\n", port)

	scheduleStart(s, concatBufferValues, sendGridSender)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection")
			continue
		}
		conn.SetDeadline(time.Now().Add(timeout))
		go handleConnection(conn, s)
	}
}

type service struct {
	Interval  time.Duration
	Buffer    *RingBuffer
	Recipient string
	Sender    string
}

func handleConnection(conn net.Conn, srvc *service) {
	msg, err := bufio.NewReader(conn).ReadString('\n')

	if err != nil {
		fmt.Printf("Error receiving message\n")
		conn.Close()
		return
	}

	message := strings.TrimSpace(string(msg))

	srvc.Buffer.Insert(message)

	conn.Write([]byte(msg + "\n"))
	conn.Close()
}

type makeBodyFunc = func(*RingBuffer) string
type sendFunc = func(string, string, string) error

func scheduleStart(srvc *service, format makeBodyFunc, send sendFunc) {
	schedule(func() {
		body := format(srvc.Buffer)
		err := send(srvc.Sender, srvc.Recipient, body)
		if err != nil {
			fmt.Println("Error while trying to send", err)
			return
		}

		srvc.Buffer.Clear()
	}, srvc.Interval)
}

func schedule(do func(), interval time.Duration) {
	ticker := time.NewTicker(interval)

	go func() {
		for {
			<-ticker.C
			do()
		}
	}()
}

func concatBufferValues(b *RingBuffer) string {
	str := ""
	b.Do(func(v interface{}) {
		if v == nil {
			return
		}

		value, ok := v.(string)
		if !ok {
			fmt.Printf("Value is not `string`. Got %T\n", v)
			return
		}
		str = str + value + "\n"
	})
	return str
}

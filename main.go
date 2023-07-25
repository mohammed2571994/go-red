package main

import (
	"fmt"
	"net"
	"os"
)

const (
	ping = "ping"
	echo = "echo"
	set  = "set"
	get  = "get"
)

const (
	simpleMessage = "+"
	bulkString    = "$"
	errorMessage  = "-"
	arrayMessage  = "*"
)

func main() {
	l, err := net.Listen("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379", err)
		os.Exit(1)
	}

	defer l.Close()
	fmt.Println("waiting for connections >>")
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	respRequest := NewRespRequest(conn)

	for {
		command, args, err := respRequest.GetRequestData()
		if err != nil {
			fmt.Println("Error reading:", err)
			return
		}

		switch command {
		case ping:
			err = HandlePing(args, conn)
		case echo:
			err = HandleEcho(args, conn)
		default:
			err = HandleUnknownCommand(args, conn)
		}

		if err != nil {
			fmt.Println("Error writing:", err)
			return
		}
	}
}

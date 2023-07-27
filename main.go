package main

import (
	"fmt"
	"net"
	"os"

	"go-red/commands"
	"go-red/parser"
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
	respRequest := parser.NewParser(conn)

	commands.InitCommands()

	for {

		command, args, err := respRequest.Parse()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		err = command.ExecuteCommand(args, conn)
		if err != nil {
			fmt.Println("Error :", err)
			return
		}

	}
}

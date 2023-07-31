package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"

	"go-red/commands"
	"go-red/config"
	"go-red/parser"
	"go-red/storage"
)

func main() {
	// init config and commands
	commands.InitCommands()
	config.InitConfig(true, "6379")

	err := loadAof()
	if err != nil {
		fmt.Println("Failed ", err)
		os.Exit(1)
	}

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
	respRequest := parser.NewParser(bufio.NewReader(conn))

	for {

		command, args, rawData, err := respRequest.Parse()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		msg, err := command.ExecuteCommand(args, rawData)
		if err != nil {
			fmt.Println("Error while executing the command :", err)
			return
		}

		_, err = conn.Write([]byte(msg))
		if err != nil {
			fmt.Println("Error while writting to the connection:", err)
			return
		}

	}
}

func loadAof() error {
	currentShouldPersist := config.ServerConfig.ShouldPersist
	config.ServerConfig.ShouldPersist = false

	file, err := storage.LoadAof()
	if err != nil {
		fmt.Println("Failed to open AOF file", err)
		os.Exit(1)
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	parser := parser.NewParser(reader)

	for {
		command, args, rawData, err := parser.Parse()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			break
		}

		fmt.Println(args)

		_, err = command.ExecuteCommand(args, rawData)
		if err != nil {
			fmt.Println("Error while executing the command:", err)
			return err
		}
	}

	config.ServerConfig.ShouldPersist = currentShouldPersist

	return err
}

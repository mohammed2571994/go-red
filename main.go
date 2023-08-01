package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"mini-go-redis/commands"
	"mini-go-redis/config"
	"mini-go-redis/parser"
	"mini-go-redis/storage"
)

func main() {
	// init config and commands
	commands.InitCommands()
	serverConfig := config.InitConfig(true, "127.0.0.1:6379", "storage.aof")

	//storage
	storage, err := storage.NewStorage(serverConfig.AofPath)
	if err != nil {
		fmt.Println("Failed during sotrage creation ", err)
		os.Exit(1)
	}

	err = loadAof(storage)
	if err != nil {
		fmt.Println("Failed ", err)
		os.Exit(1)
	}

	l, err := net.Listen("tcp", serverConfig.Address)
	if err != nil {
		fmt.Printf("Failed to bind to address: %s, error: %v", serverConfig.Address, err)
		os.Exit(1)
	}
	defer l.Close()

	fmt.Println("waiting for connections >>")
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err)
			os.Exit(1)
		}

		handleConnection(conn, storage)
	}
}

func handleConnection(conn net.Conn, storage *storage.Storage) {
	defer conn.Close()
	respRequest := parser.NewParser(bufio.NewReader(conn))

	for {

		command, args, rawData, err := respRequest.Parse()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		msg, err := command.ExecuteCommand(args, rawData, storage)
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

func loadAof(storage *storage.Storage) (err error) {
	currentShouldPersist := config.ServerConfig.ShouldPersist
	config.ServerConfig.ShouldPersist = false

	parser := parser.NewParser(storage.Reader)

	for {
		command, args, rawData, err := parser.Parse()
		if err != nil {
			break
		}

		_, err = command.ExecuteCommand(args, rawData, storage)
		if err != nil {
			fmt.Println("Error while executing the command:", err)
			return err
		}
	}

	config.ServerConfig.ShouldPersist = currentShouldPersist

	return
}

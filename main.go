package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	ping = "ping"
	echo = "echo"
	set  = "set"
	get  = "get"
)

func main() {
	l, err := net.Listen("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379", err)
		os.Exit(1)
	}

	defer l.Close()

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
	pong := "+PONG\r\n"

	defer conn.Close()

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading:", err)
			return
		}

		receivedMessage := string(buffer[:n])
		reader := strings.NewReader(receivedMessage)

		//read the first byte to skip the *
		_, err = reader.ReadByte()
		if err != nil {
			fmt.Println("Error reading: ", err.Error())
			os.Exit(1)
		}

		//read the second byte to get the number of lements
		b, err := reader.ReadByte()
		if err != nil {
			fmt.Println("Error reading: ", err.Error())
			os.Exit(1)
		}

		numberOfElements, err := strconv.Atoi(string(b))
		if err != nil {
			fmt.Println("Error converting byte to Atoi: ", err.Error())
			os.Exit(1)
		}

		//read the next two bytes to skip /r /n
		_, err = reader.ReadByte()
		_, err = reader.ReadByte()

		command := ""

		for i := 0; i < numberOfElements; i++ {
			// skip the $
			_, err = reader.ReadByte()

			// get the number of bytes in this element
			b, err := reader.ReadByte()
			if err != nil {
				fmt.Println("Error reading numberOfBytes: ", err.Error())
				os.Exit(1)
			}

			// convert the number of bytes to int
			numberOfBytes, err := strconv.Atoi(string(b))
			if err != nil {
				fmt.Println("Error converting byte to Atoi: ", err.Error())
				os.Exit(1)
			}

			// skip /r /n
			_, err = reader.ReadByte()
			_, err = reader.ReadByte()

			// read the element
			buffer := make([]byte, numberOfBytes)
			_, err = io.ReadFull(reader, buffer)
			if err != nil {
				fmt.Println("Error :", err)
				return
			}

			// the command is in the first iteration
			if i == 0 {
				command = string(buffer)
			}

			// skip /r /n
			_, err = reader.ReadByte()
			_, err = reader.ReadByte()
		}

		fmt.Printf(">> recieved command is : %s \n", command)

		_, err = conn.Write([]byte(pong))
		if err != nil {
			fmt.Println("Error writing:", err)
			return
		}
	}
}

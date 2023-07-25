package main

import (
	"fmt"
	"net"
)

func HandlePing(args []string, conn net.Conn) (err error) {
	msg := ""

	if len(args) != 0 && len(args) != 1 {
		msg = prepareResponseMessage("wrong number of arguments for command 'ping'", errorMessage)
	} else if len(args) == 0 {
		msg = prepareResponseMessage("PONG", bulkString)
	} else {
		msg = prepareResponseMessage(args[0], bulkString)
	}

	_, err = conn.Write([]byte(msg))
	return
}

func HandleEcho(args []string, conn net.Conn) (err error) {
	msg := ""
	if len(args) != 1 {
		msg = prepareResponseMessage("wrong number of arguments for command 'echo'", errorMessage)
	} else {
		msg = prepareResponseMessage(args[0], bulkString)
	}

	_, err = conn.Write([]byte(msg))
	return
}

func HandleUnknownCommand(args []string, conn net.Conn) (err error) {

	msg := prepareResponseMessage("unknown command", errorMessage)
	_, err = conn.Write([]byte(msg))

	return
}

func prepareResponseMessage(msg string, msgType string) string {
	switch msgType {
	case simpleMessage:
		return fmt.Sprintf("+%s\r\n", msg)
	case errorMessage:
		return fmt.Sprintf("-ERR %s \r\n", msg)
	case arrayMessage:
		// TODO: add the correct format
		return fmt.Sprintf("*%d\r\n%s\r\n", len(msg), msg)
	case bulkString:
		return fmt.Sprintf("$%d\r\n%s\r\n", len(msg), msg)
	}

	return fmt.Sprintf("$%d\r\n%s\r\n", len(msg), msg)
}

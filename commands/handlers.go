package commands

import (
	"fmt"
	"go-red/storage"
	"net"
	"strconv"
)

const (
	simpleMessage  = "+"
	integerMessage = ":"
	bulkMessage    = "$"
	errorMessage   = "-"
	arrayMessage   = "*"
	nullMessage    = "-1"
)

func handlePing(args []string, conn net.Conn) (err error) {
	msg := ""

	if len(args) == 0 {
		msg = marshalResponse("PONG", bulkMessage)
	} else {
		msg = marshalResponse(args[0], bulkMessage)
	}

	_, err = conn.Write([]byte(msg))
	return
}

func handleEcho(args []string, conn net.Conn) (err error) {
	msg := marshalResponse(args[0], bulkMessage)

	_, err = conn.Write([]byte(msg))
	return
}

func handleUnknownCommand(args []string, conn net.Conn) (err error) {
	msg := marshalResponse("unknown command", errorMessage)
	_, err = conn.Write([]byte(msg))

	return
}

func handleSet(args []string, conn net.Conn) (err error) {
	storage.Set(args[0], args[1])
	msg := marshalResponse("OK", bulkMessage)

	_, err = conn.Write([]byte(msg))

	return err
}

func handleGet(args []string, conn net.Conn) (err error) {
	msg := ""

	// TODO: abstract storage
	if val, ok := storage.Get(args[0]); ok {
		msg = marshalResponse(val, bulkMessage)
	} else {
		msg = marshalResponse("", nullMessage)
	}

	_, err = conn.Write([]byte(msg))
	return
}

func handleDelete(args []string, conn net.Conn) (err error) {
	numberOfDeleteItems := 0
	for _, arg := range args {
		if _, ok := storage.Get(arg); ok {
			numberOfDeleteItems++
		}

		storage.Delete(arg)
	}

	msg := marshalResponse(fmt.Sprint(numberOfDeleteItems), integerMessage)

	_, err = conn.Write([]byte(msg))
	return
}

func handleIncrement(args []string, conn net.Conn) (err error) {
	storedValue, exsists := storage.Get(args[0])
	var convertedValue int64

	if exsists {
		convertedValue, _ = strconv.ParseInt(storedValue, 10, 64)
		fmt.Println(convertedValue)
	}

	convertedValue++

	storage.Set(args[0], fmt.Sprint(convertedValue))
	msg := marshalResponse(fmt.Sprint(convertedValue), integerMessage)

	_, err = conn.Write([]byte(msg))

	return err
}

func marshalResponse(msg string, msgType string) string {
	switch msgType {
	case simpleMessage:
		return fmt.Sprintf("+%s\r\n", msg)
	case integerMessage:
		return fmt.Sprintf(":%s\r\n", msg)
	case errorMessage:
		return fmt.Sprintf("-ERR %s \r\n", msg)
	case bulkMessage:
		return fmt.Sprintf("$%d\r\n%s\r\n", len(msg), msg)
	// TODO: add array message
	default:
		return "$-1\r\n"
	}
}

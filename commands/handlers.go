package commands

import (
	"fmt"
	"go-red/storage"
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

func handlePing(args []string, rawData []byte) (msg string, err error) {
	if len(args) == 0 {
		msg = marshalResponse("PONG", bulkMessage)
	} else {
		msg = marshalResponse(args[0], bulkMessage)
	}

	return
}

func handleEcho(args []string, rawData []byte) (msg string, err error) {
	msg = marshalResponse(args[0], bulkMessage)
	return
}

func handleUnknownCommand(args []string, rawData []byte) (msg string, err error) {
	msg = marshalResponse("unknown command", errorMessage)
	return
}

func handleSet(args []string, rawData []byte) (msg string, err error) {
	msg = ""
	err = storage.Set(args[0], args[1], rawData)
	if err != nil {
		msg = marshalResponse("something went wrong", errorMessage)
	} else {
		msg = marshalResponse("OK", bulkMessage)
	}

	return
}

func handleGet(args []string, rawData []byte) (msg string, err error) {
	if val, ok := storage.Get(args[0]); ok {
		msg = marshalResponse(val, bulkMessage)
	} else {
		msg = marshalResponse("", nullMessage)
	}

	return
}

func handleDelete(args []string, rawData []byte) (msg string, err error) {
	numberOfDeleteItems := 0
	for _, arg := range args {
		if _, ok := storage.Get(arg); ok {
			numberOfDeleteItems++
		}

		// TODO: add persistence
		storage.Delete(arg)
	}

	msg = marshalResponse(fmt.Sprint(numberOfDeleteItems), integerMessage)
	return
}

func handleIncrement(args []string, rawData []byte) (msg string, err error) {
	storedValue, exsists := storage.Get(args[0])
	var convertedValue int64 = 0 //by defaule its zero

	if exsists {
		convertedValue, _ = strconv.ParseInt(storedValue, 10, 64)
	}

	convertedValue++

	err = storage.Set(args[0], fmt.Sprint(convertedValue), rawData)
	if err != nil {
		msg = marshalResponse("something went wrong", errorMessage)
	} else {
		msg = marshalResponse(fmt.Sprint(convertedValue), integerMessage)
	}

	return
}

func handleDecrement(args []string, rawData []byte) (msg string, err error) {
	storedValue, exsists := storage.Get(args[0])
	var convertedValue int64

	if exsists {
		convertedValue, _ = strconv.ParseInt(storedValue, 10, 64)
	}

	convertedValue--

	storage.Set(args[0], fmt.Sprint(convertedValue), rawData)
	if err != nil {
		msg = marshalResponse("something went wrong", errorMessage)
	} else {
		msg = marshalResponse(fmt.Sprint(convertedValue), integerMessage)
	}

	return
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

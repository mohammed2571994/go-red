package commands

import (
	"fmt"
	"go-red/storage"
	"strconv"
)

type Command struct {
	Name             string
	minArguments     int
	maxArguments     int
	handler          func(args []string, rawData []byte) (string, error)
	specialValidator func(args []string) string
}

func (command Command) ExecuteCommand(args []string, rawData []byte) (msg string, err error) {
	msg = command.validate(args)

	if msg != "" {
		msg = marshalResponse(msg, errorMessage)
		return
	}

	return command.handler(args, rawData)
}

func (command Command) validate(args []string) string {
	if (command.maxArguments != -1 && len(args) > command.maxArguments) || len(args) < command.minArguments {
		return fmt.Sprintf("wrong number of arguments for the command %s ", command.Name)
	}

	return command.specialValidator(args)
}

var pingCommand = Command{
	Name:         "ping",
	minArguments: 0,
	maxArguments: 1,
	handler:      handlePing,
	specialValidator: func(args []string) string {
		return ""
	},
}

var setCommand Command = Command{
	Name:         "set",
	minArguments: 2,
	// TODO: add support for options like expiry date
	maxArguments: 2,
	handler:      handleSet,
	specialValidator: func(args []string) string {
		return ""
	},
}

var getCommand Command = Command{
	Name:         "get",
	minArguments: 1,
	maxArguments: 1,
	handler:      handleGet,
	specialValidator: func(args []string) string {
		return ""
	},
}

var echoCommand Command = Command{
	Name:         "echo",
	minArguments: 1,
	maxArguments: 1,
	handler:      handleEcho,
	specialValidator: func(args []string) string {
		return ""
	},
}

var deleteCommand Command = Command{
	Name:         "del",
	minArguments: 1,
	maxArguments: -1,
	handler:      handleDelete,
	specialValidator: func(args []string) string {
		return ""
	},
}

var incrementCommand Command = Command{
	Name:         "incr",
	minArguments: 1,
	maxArguments: 1,
	handler:      handleIncrement,
	specialValidator: func(args []string) string {
		storedValue, exsists := storage.Get(args[0])

		if exsists {
			_, err := strconv.ParseInt(storedValue, 10, 64)
			if err != nil {
				return "value is not an integer or out of range"
			}
		}

		return ""
	},
}

var decrementCommand Command = Command{
	Name:         "decr",
	minArguments: 1,
	maxArguments: 1,
	handler:      handleDecrement,
	specialValidator: func(args []string) string {
		storedValue, exsists := storage.Get(args[0])

		if exsists {
			_, err := strconv.ParseInt(storedValue, 10, 64)
			if err != nil {
				return "value is not an integer or out of range"
			}
		}

		return ""
	},
}

var unknownCommand Command = Command{
	Name:         "unknown",
	minArguments: 0,
	maxArguments: -1,
	handler:      handleUnknownCommand,
	specialValidator: func(args []string) string {
		return ""
	},
}

var commandsMap = make(map[string]Command)

func InitCommands() {
	commandsMap[pingCommand.Name] = pingCommand
	commandsMap[echoCommand.Name] = echoCommand
	commandsMap[getCommand.Name] = getCommand
	commandsMap[setCommand.Name] = setCommand
	commandsMap[deleteCommand.Name] = deleteCommand
	commandsMap[incrementCommand.Name] = incrementCommand
	commandsMap[decrementCommand.Name] = decrementCommand
	commandsMap[unknownCommand.Name] = unknownCommand
}

func GetCommand(name string) Command {
	command, exsists := commandsMap[name]
	if exsists {
		return command
	}

	return commandsMap["unknown"]
}

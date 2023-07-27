package commands

import (
	"fmt"
	"net"
)

type Command struct {
	name             string
	minArguments     int
	maxArguments     int
	handler          func(args []string, conn net.Conn) error
	specialValidator func(args []string) string
}

func (command Command) ExecuteCommand(args []string, conn net.Conn) error {
	msg := command.validate(args)

	if msg != "" {
		msg = marshalResponse(msg, errorMessage)

		_, err := conn.Write([]byte(msg))
		return err
	}

	return command.handler(args, conn)

}

func (command Command) validate(args []string) string {
	if (command.maxArguments != -1 && len(args) > command.maxArguments) || len(args) < command.minArguments {
		return fmt.Sprintf("wrong number of arguments for the command %s ", command.name)
	}

	return command.specialValidator(args)
}

var pingCommand = Command{
	name:         "ping",
	minArguments: 0,
	maxArguments: 1,
	handler:      handlePing,
	specialValidator: func(args []string) string {
		return ""
	},
}

var setCommand Command = Command{
	name:         "set",
	minArguments: 2,
	// TODO: add support for options like expiry date
	maxArguments: 2,
	handler:      handleSet,
	specialValidator: func(args []string) string {
		return ""
	},
}

var getCommand Command = Command{
	name:         "get",
	minArguments: 1,
	maxArguments: 1,
	handler:      handleGet,
	specialValidator: func(args []string) string {
		return ""
	},
}

var echoCommand Command = Command{
	name:         "echo",
	minArguments: 1,
	maxArguments: 1,
	handler:      handleEcho,
	specialValidator: func(args []string) string {
		return ""
	},
}

var deleteCommand Command = Command{
	name:         "del",
	minArguments: 1,
	maxArguments: -1,
	handler:      handleDelete,
	specialValidator: func(args []string) string {
		return ""
	},
}

var unknownCommand Command = Command{
	name:         "unknown",
	minArguments: 0,
	maxArguments: -1,
	handler:      handleUnknownCommand,
	specialValidator: func(args []string) string {
		return ""
	},
}

var commandsMap = make(map[string]Command)

func InitCommands() {
	commandsMap["ping"] = pingCommand
	commandsMap["echo"] = echoCommand
	commandsMap["get"] = getCommand
	commandsMap["set"] = setCommand
	commandsMap["del"] = deleteCommand
	commandsMap["unknown"] = unknownCommand
}

func GetCommand(name string) Command {
	command, exsists := commandsMap[name]
	if exsists {
		return command
	}

	return commandsMap["unknown"]
}

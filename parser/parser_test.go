package parser

import (
	"bufio"
	"go-red/commands"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	commands.InitCommands()

	testCases := []struct {
		name          string
		checkResponse func(t *testing.T)
	}{
		{
			name: "No args Ping",
			checkResponse: func(t *testing.T) {
				expecteArgsLength := 0
				recievedData := "*1\r\n$4\r\nPING\r\n"
				expectedCommandName := "ping"

				parser := NewParser(bufio.NewReader(strings.NewReader(recievedData)))

				command, args, _, err := parser.Parse()

				if err != nil {
					t.Errorf("returned error =  %v, expected nil", err)
				}

				if len(args) != expecteArgsLength {
					t.Errorf("args length is =  %d, expected is = %d", len(args), expecteArgsLength)
				}

				if command.Name != expectedCommandName {
					t.Errorf("command is =  %s, expected is = %s", command.Name, expectedCommandName)
				}
			},
		},

		{
			name: "One arg Ping",
			checkResponse: func(t *testing.T) {
				expecteArgsLength := 1
				recievedData := "*2\r\n$4\r\nPING\r\n$2\r\nHi\r\n"
				expectedCommandName := "ping"

				parser := NewParser(bufio.NewReader(strings.NewReader(recievedData)))

				command, args, _, err := parser.Parse()

				if err != nil {
					t.Errorf("returned error =  %v, expected nil", err)
				}

				if len(args) != expecteArgsLength {
					t.Errorf("args length is =  %d, expected is = %d", len(args), expecteArgsLength)
				}

				if command.Name != expectedCommandName {
					t.Errorf("command is =  %s, expected is = %s", command.Name, expectedCommandName)
				}
			},
		},

		{
			name: "Two args Ping",
			checkResponse: func(t *testing.T) {
				expecteArgsLength := 2
				recievedData := "*3\r\n$4\r\nPING\r\n$5\r\nHello\r\n$5\r\nWorld\r\n"
				expectedCommandName := "ping"

				parser := NewParser(bufio.NewReader(strings.NewReader(recievedData)))

				command, args, _, err := parser.Parse()

				if err != nil {
					t.Errorf("returned error =  %v, expected nil", err)
				}

				if len(args) != expecteArgsLength {
					t.Errorf("args length is =  %d, expected is = %d", len(args), expecteArgsLength)
				}

				if command.Name != expectedCommandName {
					t.Errorf("command is =  %s, expected is = %s", command.Name, expectedCommandName)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.checkResponse(t)
		})
	}

}

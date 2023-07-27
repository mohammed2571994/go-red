package parser

import (
	"bufio"
	"net"
	"strconv"
	"strings"

	"go-red/commands"
)

type Parser struct {
	reader *bufio.Reader
}

func NewParser(conn net.Conn) *Parser {
	return &Parser{reader: bufio.NewReader(conn)}
}

func (resp *Parser) Parse() (command commands.Command, args []string, err error) {
	numberOfElements, err := resp.readInteger()
	if err != nil {
		return
	}

	numberOfBytes, err := resp.readInteger()
	if err != nil {
		return
	}

	bytes := make([]byte, numberOfBytes)

	_, err = resp.reader.Read(bytes)
	if err != nil {
		return
	}

	// skip /r /n
	resp.skipBytes(2)

	args, err = resp.readArgs(numberOfElements)
	if err != nil {
		return
	}

	command = commands.GetCommand(strings.ToLower(string(bytes)))
	return command, args, nil
}

func (resp *Parser) readInteger() (n int, err error) {
	bytes, err := resp.readLine()
	if err != nil {
		return 0, err
	}

	numberOfBytes, err := strconv.Atoi(string(bytes[1:]))
	if err != nil {
		return 0, err
	}

	return numberOfBytes, nil
}

func (resp *Parser) readLine() (line []byte, err error) {
	line, err = resp.reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	return line[:len(line)-2], nil
}

func (resp *Parser) readArgs(numberOfElements int) (args []string, err error) {
	for i := 0; i < numberOfElements-1; i++ {
		numberOfBytes := 0
		numberOfBytes, err = resp.readInteger()
		if err != nil {
			return
		}

		bytes := make([]byte, numberOfBytes)

		_, err = resp.reader.Read(bytes)
		if err != nil {
			return
		}

		args = append(args, string(bytes))

		// skip /r /n
		err = resp.skipBytes(2)
		if err != nil {
			return
		}

	}

	return args, nil
}

func (resp *Parser) skipBytes(n int) (err error) {
	for i := 0; i < n; i++ {
		_, err = resp.reader.ReadByte()
		if err != nil {
			return
		}
	}

	return
}

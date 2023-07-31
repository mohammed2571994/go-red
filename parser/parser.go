package parser

import (
	"bufio"
	"strconv"
	"strings"

	"go-red/commands"
)

type Parser struct {
	reader *bufio.Reader
}

// TODO: find a better way than using a global variable
var RawData []byte

func NewParser(reader *bufio.Reader) *Parser {
	return &Parser{reader: reader}
}

func (parser *Parser) Parse() (command commands.Command, args []string, rawData []byte, err error) {
	RawData = []byte{}

	numberOfElements, err := parser.readInteger()
	if err != nil {
		return
	}

	numberOfBytes, err := parser.readInteger()
	if err != nil {
		return
	}

	bytes := make([]byte, numberOfBytes)

	_, err = parser.reader.Read(bytes)
	if err != nil {
		return
	}
	RawData = append(RawData, bytes...)

	// skip /r /n
	parser.skipBytes(2)

	args, err = parser.readArgs(numberOfElements)
	if err != nil {
		return
	}

	command = commands.GetCommand(strings.ToLower(string(bytes)))
	return command, args, RawData, nil
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

	RawData = append(RawData, line...)

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
		RawData = append(RawData, bytes...)

		args = append(args, string(bytes))

		// skip /r /n
		err = resp.skipBytes(2)
		if err != nil {
			return
		}

	}

	return args, nil
}

func (resp *Parser) skipBytes(n int) error {
	for i := 0; i < n; i++ {
		b, err := resp.reader.ReadByte()
		if err != nil {
			return err
		}

		RawData = append(RawData, b)
	}

	return nil
}

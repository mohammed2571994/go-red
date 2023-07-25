package main

import (
	"bufio"
	"net"
	"strconv"
)

type RespRequest struct {
	reader *bufio.Reader
}

func NewRespRequest(conn net.Conn) *RespRequest {
	return &RespRequest{reader: bufio.NewReader(conn)}
}

func (resp *RespRequest) GetRequestData() (command string, args []string, err error) {
	numberOfElements, err := resp.readNumberOfBytes()
	if err != nil {
		return
	}

	numberOfBytes, err := resp.readNumberOfBytes()
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

	args, err = resp.getArgs(numberOfElements)
	if err != nil {
		return
	}

	return string(bytes), args, nil
}

func (resp *RespRequest) readNumberOfBytes() (n int, err error) {
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

func (resp *RespRequest) readLine() (line []byte, err error) {

	for {
		b, err := resp.reader.ReadByte()

		if err != nil {
			return nil, err
		}

		if b == '\r' {
			//skip \n
			err = resp.skipBytes(1)
			if err != nil {
				return nil, err
			}

			break
		}

		line = append(line, b)
	}

	return line, nil
}

func (resp *RespRequest) getArgs(numberOfElements int) (args []string, err error) {
	var numberOfBytes int

	for i := 0; i < numberOfElements-1; i++ {
		numberOfBytes, err = resp.readNumberOfBytes()
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

func (resp *RespRequest) skipBytes(n int) (err error) {
	for i := 0; i < n; i++ {
		_, err = resp.reader.ReadByte()
		if err != nil {
			return
		}
	}

	return
}

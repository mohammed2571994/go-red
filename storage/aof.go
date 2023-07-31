package storage

import (
	"bufio"
	"fmt"
	"os"
)

type Storage struct {
	File   *os.File
	Reader *bufio.Reader
	path   string
	m      map[string]string
}

func NewStorage(path string) (*Storage, error) {
	file, err := os.OpenFile(path, os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(file)

	return &Storage{
		path:   path,
		File:   file,
		Reader: reader,
		m:      make(map[string]string),
	}, nil
}

func writeToFile(rawCommand []byte) error {
	file, err := os.OpenFile("db.txt", os.O_CREATE|os.O_APPEND, os.ModeAppend)
	if err != nil {
		fmt.Println("error while opening the file", err)
	}
	defer file.Close()

	_, err = file.Write(rawCommand)
	if err != nil {
		return err
	}

	return nil
}

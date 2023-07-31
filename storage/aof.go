package storage

import (
	"fmt"
	"os"
)

type Storage struct {
	file *os.File
	path string
}

func NewAof(path string, flag int) (*Storage, error) {
	file, err := os.OpenFile(path, flag, os.ModeAppend)
	if err != nil {
		return nil, err
	}

	return &Storage{
		path: path,
		file: file,
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

func LoadAof() (file *os.File, err error) {
	file, err = os.OpenFile("db.txt", os.O_RDONLY, os.ModeAppend)
	if err != nil {
		return
	}

	return
}

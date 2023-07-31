package storage

import (
	"bufio"
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

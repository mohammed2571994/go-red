package storage

import "go-red/config"

var m = make(map[string]string)

func (storage *Storage) Set(key string, value string, rawData []byte) (err error) {
	m[key] = value

	if config.ServerConfig.ShouldPersist {
		err = writeToFile(rawData)
		if err != nil {
			delete(m, key)
		}
	}

	return
}

func (storage *Storage) Get(key string) (value string, exists bool) {
	value, exists = m[key]
	return
}

func (storage *Storage) Delete(key string) (err error) {
	delete(m, key)
	return
}

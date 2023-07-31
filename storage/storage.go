package storage

import "go-red/config"

var m = make(map[string]string)

func Set(key string, value string, rawData []byte) (err error) {
	m[key] = value

	if config.ServerConfig.ShouldPersist {
		err = writeToFile(rawData)
		if err != nil {
			delete(m, key)
		}
	}

	return
}

func Get(key string) (value string, exists bool) {
	value, exists = m[key]
	return
}

func Delete(key string) (err error) {
	delete(m, key)
	return
}

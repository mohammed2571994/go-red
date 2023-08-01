package storage

import "mini-go-redis/config"

func (storage *Storage) Set(key string, value string, rawData []byte) (err error) {
	storage.m[key] = value

	if config.ServerConfig.ShouldPersist {
		err = storage.writeToAof(rawData)
		if err != nil {
			delete(storage.m, key)
		}
	}

	return
}

func (storage *Storage) Get(key string) (value string, exists bool) {
	value, exists = storage.m[key]
	return
}

func (storage *Storage) Delete(key string, rawData []byte) (err error) {
	val := storage.m[key]

	delete(storage.m, key)
	if config.ServerConfig.ShouldPersist {
		err = storage.writeToAof(rawData)
		if err != nil {
			storage.m[key] = val
		}
	}

	return
}

func (storage *Storage) writeToAof(rawCommand []byte) error {
	_, err := storage.File.Write(rawCommand)
	if err != nil {
		return err
	}

	return nil
}

package storage

var m = make(map[string]string)

func Set(key string, value string) {
	m[key] = value
}

func Get(key string) (value string, exists bool) {
	value, exists = m[key]
	return
}

func Delete(key string) {
	delete(m, key)
}

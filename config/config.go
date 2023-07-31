package config

type Config struct {
	ShouldPersist bool
	Port          string
}

var ServerConfig = &Config{}

func InitConfig(shouldPersist bool, port string) {
	ServerConfig.Port = port
	ServerConfig.ShouldPersist = shouldPersist
}

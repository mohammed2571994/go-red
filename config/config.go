package config

type Config struct {
	ShouldPersist bool
	Address       string
	AofPath       string
}

var ServerConfig = &Config{}

func InitConfig(shouldPersist bool, address string, aofPath string) *Config {
	ServerConfig.Address = address
	ServerConfig.ShouldPersist = shouldPersist
	ServerConfig.AofPath = aofPath

	return ServerConfig
}

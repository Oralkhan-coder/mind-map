package config

type DbConfig struct {
	Host     string `env:"MONGO_HOST"`
	Port     uint16 `env:"MONGO_PORT"`
	Username string `env:"MONGO_USERNAME"`
	Password string `env:"MONGO_PASSWORD"`
	Database string `env:"MONGO_DATABASE"`
}

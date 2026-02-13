package config

import (
	"log"
	"os"

	db "github.com/Oralkhan-coder/mind-map/pkg"
	"github.com/joho/godotenv"
)

type Config struct {
	Secret *SecretConfig
	Mongo  *db.DbConfig
}

func InitConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	secretConfig := NewSecretConfig(os.Getenv("JWT_SECRET"), os.Getenv("GOOGLE_CLIENT_ID"))
	dbConfig := &db.DbConfig{
		Host:     os.Getenv("MONGO_HOST"),
		Database: os.Getenv("MONGO_DATABASE"),
		Username: os.Getenv("MONGO_USERNAME"),
		Password: os.Getenv("MONGO_PASSWORD"),
		Port:     27017,
	}

	return &Config{
		Mongo:  dbConfig,
		Secret: secretConfig,
	}
}

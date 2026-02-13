package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Secret *SecretConfig
	Mongo  *DbConfig
	Email  *EmailConfig
}

func InitConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	secretConfig := NewSecretConfig(os.Getenv("JWT_SECRET"), os.Getenv("GOOGLE_CLIENT_ID"))
	dbConfig := &DbConfig{
		Host:     os.Getenv("MONGO_HOST"),
		Database: os.Getenv("MONGO_DATABASE"),
		Username: os.Getenv("MONGO_USERNAME"),
		Password: os.Getenv("MONGO_PASSWORD"),
		Port:     27017,
	}
	emailConfig := &EmailConfig{
		From:     os.Getenv("EMAIL_USER"),
		SmtpHost: os.Getenv("SMTP_SERVER"),
		SmtpPort: os.Getenv("SMTP_PORT"),
		Password: os.Getenv("EMAIL_PASSWORD"),
	}

	return &Config{
		Secret: secretConfig,
		Mongo:  dbConfig,
		Email:  emailConfig,
	}
}

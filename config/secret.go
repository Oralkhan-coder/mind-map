package config

import "log"

type SecretConfig struct {
	JwtSecret      string
	GoogleClientID string
}

func NewSecretConfig(jwtSecret, googleClientID string) *SecretConfig {
	if jwtSecret == "" {
		log.Fatal("JWT Secret can't be empty")
	}

	return &SecretConfig{
		JwtSecret:      jwtSecret,
		GoogleClientID: googleClientID,
	}
}

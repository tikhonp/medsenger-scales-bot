// Package util provides configuration loading, echo http helpers and other functionality for the application.
package util

import (
	"os"
	"strconv"
)

type Config struct {
	Server *Server

	DB *Database

	// Sentry configutation URL.
	SentryDSN string
}

type Database struct {
	User string

	Password string

	Dbname string

	Host string
}

type Server struct {
	// The hostname of this application.
	Host string

	// The port to listen on.
	Port uint16

	// Medsenger Agent secret key.
	MedsengerAgentKey string

	// Sets server to debug mode.
	Debug bool
}

func LoadConfigFromEnv() *Config {
	serverPort, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		panic(err)
	}
	return &Config{
		Server: &Server{
			Host:              os.Getenv("SERVER_HOST"),
			Port:              uint16(serverPort),
			MedsengerAgentKey: os.Getenv("SCALES_KEY"),
			Debug:             os.Getenv("DEBUG") == "true",
		},
		DB: &Database{
			User:     os.Getenv("DB_LOGIN"),
			Password: os.Getenv("DB_PASSWORD"),
			Dbname:   os.Getenv("DB_DATABASE"),
			Host:     os.Getenv("DB_HOST"),
		},
		SentryDSN:          os.Getenv("SENTRY_DSN"),
	}
}


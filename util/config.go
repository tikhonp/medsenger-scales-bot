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
	Host string

	Port int

	User string

	Password string

	Dbname string
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
	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
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
			Host:     os.Getenv("DB_HOST"),
			Port:     dbPort,
			User:     os.Getenv("DB_LOGIN"),
			Password: os.Getenv("DB_PASSWORD"),
			Dbname:   os.Getenv("DB_DATABASE"),
		},
		SentryDSN: os.Getenv("SENTRY_DSN"),
	}
}

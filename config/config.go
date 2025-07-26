package config

import (
	"flag"
	"os"
	"time"
)

type token struct {
	Secret string
	Exp    time.Time
}

type Config struct {
	Port  string `json:"port"`
	DSN   string
	Token token
}

func NewConfig() *Config {
	port := flag.String("p", "", "port to listen on")
	dsn := flag.String("d", "", "dsn to connect to")
	flag.Parse()

	if envAddress := os.Getenv("PORT"); envAddress != "" {
		*port = envAddress
	}

	if envDSN := os.Getenv("DSN"); envDSN != "" {
		*dsn = envDSN
	}

	return &Config{
		Port: *port,
		DSN:  *dsn,
	}
}

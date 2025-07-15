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
	Address     string `json:"address"`
	EnableHTTPS bool   `json:"enable_https"`
	EnableGRPC  bool   `json:"enable_grpc"`
	DSN         string `json:"dsn"`
	Token       token
}

func NewConfig() Config {
	address := flag.String("a", "", "address to listen on")
	enableHTTPS := flag.Bool("enable-https", false, "enable https")
	enableGRPC := flag.Bool("enable-grpc", false, "enable grpc")
	dsn := flag.String("d", "", "dsn to connect to")
	flag.Parse()

	if envAddress := os.Getenv("ADDRESS"); envAddress != "" {
		*address = envAddress
	}

	if envEnableHTTPS := os.Getenv("ENABLE_HTTPS"); envEnableHTTPS != "" {
		*enableHTTPS = envEnableHTTPS == "true"
	}

	if envEnableGRPC := os.Getenv("ENABLE_GRPC"); envEnableGRPC != "" {
		*enableGRPC = envEnableGRPC == "true"
	}

	if envDSN := os.Getenv("DSN"); envDSN != "" {
		*dsn = envDSN
	}

	return Config{
		Address:     *address,
		EnableHTTPS: *enableHTTPS,
		EnableGRPC:  *enableGRPC,
		DSN:         *dsn,
	}
}

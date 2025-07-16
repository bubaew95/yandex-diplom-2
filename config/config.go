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
	Port          string `json:"port"`
	ServerAddress string `json:"server_address"`
	EnableHTTPS   bool   `json:"enable_https"`
	EnableGRPC    bool   `json:"enable_grpc"`
	DSN           string `json:"dsn"`
	Token         token
}

func NewConfig() Config {
	port := flag.String("p", "", "port to listen on")
	serverAddress := flag.String("a", "", "serverAddress to listen on")
	enableHTTPS := flag.Bool("enable-https", false, "enable https")
	enableGRPC := flag.Bool("enable-grpc", false, "enable grpc")
	dsn := flag.String("d", "", "dsn to connect to")
	flag.Parse()

	if envAddress := os.Getenv("ADDRESS"); envAddress != "" {
		*port = envAddress
	}

	if envServerAddress := os.Getenv("DOMAIN"); envServerAddress != "" {
		*serverAddress = envServerAddress
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
		Port:          *port,
		EnableHTTPS:   *enableHTTPS,
		EnableGRPC:    *enableGRPC,
		DSN:           *dsn,
		ServerAddress: *serverAddress,
	}
}

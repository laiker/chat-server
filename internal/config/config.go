package config

import (
	"github.com/joho/godotenv"
)

var ConfigPathKey = "configPathKey"

type GRPCConfig interface {
	Address() string
	Host() string
}

type PGConfig interface {
	DSN() string
}

type AuthConfig interface {
	Host() string
}

func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}

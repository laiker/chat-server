package config

import (
	"github.com/joho/godotenv"
)

var ConfigPathKey = "configPathKey"

type GRPCConfig interface {
	Address() string
}

type PGConfig interface {
	DSN() string
}

func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}

package env

import (
	"errors"
	"fmt"
	"os"

	"github.com/laiker/chat-server/internal/config"
)

const (
	dbHost = "POSTGRES_HOST"
	dbName = "POSTGRES_DB"
	dbUser = "POSTGRES_USER"
	dbPass = "POSTGRES_PASSWORD" //nolint:golint,gosec
	dbPort = "POSTGRES_PORT"
)

var _ config.PGConfig = (*pgConfig)(nil)

type pgConfig struct {
	host string
	name string
	user string
	pass string
	port string
}

func NewPGConfig() (config.PGConfig, error) {
	host := os.Getenv(dbHost)
	if len(host) == 0 {
		return nil, errors.New("pg host not found")
	}
	name := os.Getenv(dbName)
	if len(name) == 0 {
		return nil, errors.New("pg name not found")
	}
	user := os.Getenv(dbUser)
	if len(user) == 0 {
		return nil, errors.New("pg user not found")
	}
	pass := os.Getenv(dbPass)
	if len(pass) == 0 {
		return nil, errors.New("pg pass not found")
	}
	port := os.Getenv(dbPort)
	if len(port) == 0 {
		return nil, errors.New("pg port not found")
	}

	return &pgConfig{
		host: host,
		name: name,
		user: user,
		pass: pass,
		port: port,
	}, nil
}

func (cfg *pgConfig) DSN() string {
	return fmt.Sprintf("postgresql://%v:%v@%v:%v/%v", cfg.user, cfg.pass, cfg.host, cfg.port, cfg.name)
}

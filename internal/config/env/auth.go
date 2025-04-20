package env

import (
	"os"

	"github.com/laiker/chat-server/internal/config"
	"github.com/pkg/errors"
)

const (
	authHostEnvName = "AUTH_HOST"
)

var _ config.AuthConfig = (*authConfig)(nil)

type authConfig struct {
	host string
}

func NewAuthConfig() (config.AuthConfig, error) {
	host := os.Getenv(authHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("auth host not found")
	}

	return &grpcConfig{
		host: host,
	}, nil
}

func (cfg *authConfig) Host() string {
	return cfg.host
}

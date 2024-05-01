package configs

import (
	"github.com/stdyum/api-common/env"
	"github.com/stdyum/api-common/server"
)

type Model struct {
	Ports          server.PortConfig `env:"PORT"`
	AuthServiceURL string            `env:"AUTH_SERVICE_URL"`
}

var Config Model

func init() {
	err := env.Fill(&Config)
	if err != nil {
		panic("cannot fill config: " + err.Error())
	}
}

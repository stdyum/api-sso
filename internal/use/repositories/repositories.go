package repositories

import (
	"github.com/stdyum/api-sso/internal/configs"
	"github.com/stdyum/api-sso/internal/use/repositories/auth"
)

type Repositories struct {
	Auth auth.Repository
}

func New(config configs.Model, initial ...*Repositories) *Repositories {
	return &Repositories{
		Auth: auth.New(config.AuthServiceURL),
	}
}

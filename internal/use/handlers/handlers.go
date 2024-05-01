package handlers

import (
	"github.com/stdyum/api-sso/internal/configs"
	"github.com/stdyum/api-sso/internal/use/controllers"
	"github.com/stdyum/api-sso/internal/use/handlers/http"
)

type Handlers struct {
	Http http.Handler
}

func New(config configs.Model, controllers *controllers.Controllers, initial ...*Handlers) *Handlers {
	return &Handlers{
		Http: http.New(config.Ports.HTTP, controllers.Main),
	}
}

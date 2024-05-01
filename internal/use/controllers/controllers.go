package controllers

import (
	"github.com/stdyum/api-sso/internal/use/controllers/controller"
	"github.com/stdyum/api-sso/internal/use/controllers/controller/validators"
	"github.com/stdyum/api-sso/internal/use/repositories"
)

type Controllers struct {
	Main controller.Controller
}

func New(repositories *repositories.Repositories, initial ...*Controllers) *Controllers {
	return &Controllers{
		Main: controller.New(validators.New(), repositories.Auth),
	}
}

package use

import (
	"github.com/stdyum/api-sso/internal/configs"
	"github.com/stdyum/api-sso/internal/use/controllers"
	"github.com/stdyum/api-sso/internal/use/errors"
	"github.com/stdyum/api-sso/internal/use/handlers"
	"github.com/stdyum/api-sso/internal/use/repositories"
)

type Use struct {
	Config configs.Model

	Handlers     *handlers.Handlers
	Controllers  *controllers.Controllers
	Repositories *repositories.Repositories
}

func Default() (*Use, error) {
	return New()
}

func New(initial ...*Use) (*Use, error) {
	errors.Register()

	useCase := &Use{}

	useCase.Repositories = repositories.New(configs.Config, useCase.Repositories)

	useCase.Controllers = controllers.New(useCase.Repositories, useCase.Controllers)

	useCase.Handlers = handlers.New(configs.Config, useCase.Controllers)

	return useCase, nil
}

func (u *Use) Run() error {
	return u.Handlers.Http.Run()
}

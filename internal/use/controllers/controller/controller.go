package controller

import (
	"context"

	"github.com/stdyum/api-sso/internal/use/controllers/controller/dto"
	"github.com/stdyum/api-sso/internal/use/controllers/controller/validators"
	"github.com/stdyum/api-sso/internal/use/repositories/auth"
	"github.com/stdyum/api-sso/internal/use/repositories/auth/entities"
)

type Controller interface {
	Login(ctx context.Context, request dto.LoginRequest) (dto.TokenPairResponse, error)
	Update(ctx context.Context, request dto.UpdateRequest) (dto.TokenPairResponse, error)
}

type controller struct {
	validator validators.Validator
	auth      auth.Repository
}

func New(validator validators.Validator, auth auth.Repository) Controller {
	return &controller{
		validator: validator,
		auth:      auth,
	}
}

func (c *controller) Login(ctx context.Context, request dto.LoginRequest) (dto.TokenPairResponse, error) {
	authRequest := entities.LoginRequest{
		Login:    request.Login,
		Password: request.Password,
	}

	tokens, err := c.auth.Login(ctx, authRequest)
	if err != nil {
		return dto.TokenPairResponse{}, err
	}

	return dto.TokenPairResponse{
		Access:  tokens.Access,
		Refresh: tokens.Refresh,
	}, nil
}

func (c *controller) Update(ctx context.Context, request dto.UpdateRequest) (dto.TokenPairResponse, error) {
	updateRequest := entities.UpdateRequest{
		RefreshToken: request.RefreshToken,
	}

	tokens, err := c.auth.Update(ctx, updateRequest)
	if err != nil {
		return dto.TokenPairResponse{}, err
	}

	return dto.TokenPairResponse{
		Access:  tokens.Access,
		Refresh: tokens.Refresh,
	}, nil
}

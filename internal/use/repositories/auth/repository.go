package auth

import (
	"context"

	"github.com/stdyum/api-sso/internal/use/repositories/auth/entities"
)

type Repository interface {
	Login(ctx context.Context, data entities.LoginRequest) (entities.TokenPair, error)
	Update(ctx context.Context, data entities.UpdateRequest) (entities.TokenPair, error)
}

type repository struct {
	authURL string
}

func New(url string) Repository {
	return &repository{
		authURL: url,
	}
}

func (r *repository) Login(ctx context.Context, data entities.LoginRequest) (entities.TokenPair, error) {
	tokenPair := struct {
		Tokens entities.TokenPair `json:"tokens"`
	}{}

	if err := r.sendPostRequest(ctx, data, "api/v1/login", &tokenPair); err != nil {
		return entities.TokenPair{}, err
	}

	return tokenPair.Tokens, nil
}

func (r *repository) Update(ctx context.Context, data entities.UpdateRequest) (entities.TokenPair, error) {
	tokenPair := struct {
		Tokens entities.TokenPair `json:"tokens"`
	}{}
	if err := r.sendPostRequest(ctx, data, "api/v1/token/update", &tokenPair); err != nil {
		return entities.TokenPair{}, err
	}

	return tokenPair.Tokens, nil
}

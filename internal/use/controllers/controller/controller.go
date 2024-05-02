package controller

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"github.com/stdyum/api-sso/internal/use/controllers/controller/dto"
	"github.com/stdyum/api-sso/internal/use/controllers/controller/validators"
	"github.com/stdyum/api-sso/internal/use/repositories/auth"
	"github.com/stdyum/api-sso/internal/use/repositories/auth/entities"
)

type Controller interface {
	Login(ctx context.Context, request dto.LoginRequest) (dto.TokenPairResponse, error)
	Update(ctx context.Context, request dto.UpdateRequest) (dto.TokenPairResponse, error)
	Authorize(ctx context.Context, request dto.AuthorizeRequest) (dto.TokenPairResponse, error)
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

func (c *controller) Authorize(ctx context.Context, request dto.AuthorizeRequest) (dto.TokenPairResponse, error) {
	if ok := c.isJWTTokenValid(request.AccessToken); ok {
		return dto.TokenPairResponse{
			Access:  request.AccessToken,
			Refresh: request.RefreshToken,
		}, nil
	}

	updateRequest := dto.UpdateRequest{
		RefreshToken: request.RefreshToken,
	}

	return c.Update(ctx, updateRequest)
}

func (c *controller) isJWTTokenValid(token string) bool {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false
	}

	claimsBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}

	var claims struct {
		Exp int64 `json:"exp"`
	}

	if err = json.Unmarshal(claimsBytes, &claims); err != nil {
		return false
	}

	expireAt := time.UnixMilli(claims.Exp * 1000)
	return expireAt.After(time.Now())
}

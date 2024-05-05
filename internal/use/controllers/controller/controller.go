package controller

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/stdyum/api-sso/internal/use/controllers/controller/dto"
	"github.com/stdyum/api-sso/internal/use/controllers/controller/models"
	"github.com/stdyum/api-sso/internal/use/controllers/controller/validators"
	"github.com/stdyum/api-sso/internal/use/repositories/auth"
	"github.com/stdyum/api-sso/internal/use/repositories/auth/entities"
)

type Controller interface {
	Login(ctx context.Context, request dto.LoginRequest) (dto.TokenPairResponse, error)
	Update(ctx context.Context, request dto.UpdateRequest) (dto.TokenPairResponse, error)
	Authorize(ctx context.Context, request dto.AuthorizeRequest) (dto.UserWithTokensResponse, error)
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

func (c *controller) Authorize(ctx context.Context, request dto.AuthorizeRequest) (dto.UserWithTokensResponse, error) {
	claims, err := c.parseJWTToken(request.AccessToken)
	if err == nil {
		return dto.UserWithTokensResponse{
			Tokens: dto.TokenPairResponse{
				Access:  request.AccessToken,
				Refresh: request.RefreshToken,
			},
			User: dto.UserResponse{
				Id:            claims.User.Id,
				Login:         claims.User.Login,
				PictureURL:    claims.User.PictureURL,
				Email:         claims.User.Email,
				VerifiedEmail: claims.User.VerifiedEmail,
			},
		}, nil
	}

	updateRequest := dto.UpdateRequest{
		RefreshToken: request.RefreshToken,
	}

	tokens, err := c.Update(ctx, updateRequest)
	if err != nil {
		return dto.UserWithTokensResponse{}, err
	}

	claims, err = c.parseJWTToken(tokens.Access)
	if err != nil {
		return dto.UserWithTokensResponse{}, err
	}

	return dto.UserWithTokensResponse{
		Tokens: dto.TokenPairResponse{
			Access:  tokens.Access,
			Refresh: tokens.Refresh,
		},
		User: dto.UserResponse{
			Id:            claims.User.Id,
			Login:         claims.User.Login,
			PictureURL:    claims.User.PictureURL,
			Email:         claims.User.Email,
			VerifiedEmail: claims.User.VerifiedEmail,
		},
	}, nil
}

func (c *controller) parseJWTToken(token string) (claims models.Claims, err error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return models.Claims{}, errors.New("token contains an invalid number of segments")
	}

	claimsBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return models.Claims{}, err
	}

	err = json.Unmarshal(claimsBytes, &claims)
	return
}

func (c *controller) isJWTClaimsExpired(claims models.Claims) bool {
	expireAt := time.UnixMilli(claims.Exp * 1000)
	return expireAt.After(time.Now())
}

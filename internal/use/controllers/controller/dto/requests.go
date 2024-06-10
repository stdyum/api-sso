package dto

import (
	"time"
)

type LoginRequest struct {
	Login               string    `json:"login"`
	Password            string    `json:"password"`
	SessionExpirationAt time.Time `json:"sessionExpirationAt"`
}

type UpdateRequest struct {
	RefreshToken string `json:"refresh"`
}

type AuthorizeRequest struct {
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
}

type SetTokensRequest struct {
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
}

type SetDefaultEnrollmentIdRequest struct {
	Id string `json:"id"`
}

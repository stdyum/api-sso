package entities

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

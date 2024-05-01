package entities

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UpdateRequest struct {
	RefreshToken string `json:"refresh"`
}
package dto

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UpdateRequest struct {
	RefreshToken string `json:"refresh"`
}

type AuthorizeRequest struct {
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
}

type SetDefaultEnrollmentIdRequest struct {
	Id string `json:"id"`
}

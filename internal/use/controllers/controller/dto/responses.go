package dto

type TokenPairResponse struct {
	Access  string `json:"access"`
	Refresh string `json:"-"`
}

type UserResponse struct {
	Id            string `json:"id"`
	Login         string `json:"login"`
	PictureURL    string `json:"pictureURL"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verifiedEmail"`
}

type EnrollmentResponse struct {
	Id string `json:"id"`
}

type UserWithTokensAndEnrollmentResponse struct {
	User       UserResponse       `json:"user"`
	Tokens     TokenPairResponse  `json:"tokens"`
	Enrollment EnrollmentResponse `json:"enrollment"`
}

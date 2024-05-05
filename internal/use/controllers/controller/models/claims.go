package models

type Claims struct {
	Exp  int64 `json:"exp"`
	User User  `json:"claims"`
}

type User struct {
	Id            string `json:"userID"`
	Login         string `json:"login"`
	PictureURL    string `json:"pictureURL"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verifiedEmail"`
}

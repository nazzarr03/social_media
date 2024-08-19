package user

import "github.com/nazzarr03/social-media/models"

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserDTO struct {
	UserID   int           `json:"user_id"`
	Username string        `json:"username"`
	Password string        `json:"password"`
	Email    string        `json:"email"`
	ImageURL string        `json:"image_url"`
	Posts    []models.Post `json:"posts"`
	Friends  []models.User `json:"friends"`
}

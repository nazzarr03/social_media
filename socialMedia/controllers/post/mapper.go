package post

import "github.com/nazzarr03/social-media/models"

type CreatePostRequest struct {
	Content string `json:"content"`
}

type UpdatePostRequest struct {
	Content string `json:"content"`
}

type PostDTO struct {
	PostID   int              `json:"post_id"`
	Content  string           `json:"content"`
	ImageURL string           `json:"image_url"`
	UserID   int              `json:"user_id"`
	Comments []models.Comment `json:"comments"`
	Likes    []models.Like    `json:"likes"`
}

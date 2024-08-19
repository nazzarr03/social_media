package comment

import "github.com/nazzarr03/social-media/models"

type CreateCommentRequest struct {
	Content string `json:"content"`
}

type CommentDTO struct {
	CommentID       int              `json:"comment_id"`
	Content         string           `json:"content"`
	ImageURL        string           `json:"image_url"`
	PostID          int              `json:"post_id"`
	UserID          int              `json:"user_id"`
	ParentCommentID *int             `json:"parent_comment"`
	Comments        []models.Comment `json:"comments"`
	Likes           []models.Like    `json:"likes"`
}

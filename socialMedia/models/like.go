package models

import "time"

type Like struct {
	LikeID    int  `json:"like_id" gorm:"primaryKey;autoIncrement"`
	IsLiked   bool `json:"is_liked" gorm:"default:false"`
	UserID    int  `json:"user_id"`
	PostID    int  `json:"post_id"`
	CommentID *int `json:"comment_id"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

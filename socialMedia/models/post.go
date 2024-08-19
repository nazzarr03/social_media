package models

import "time"

type Post struct {
	PostID   int       `json:"post_id" gorm:"primaryKey;autoIncrement"`
	Content  string    `json:"content" validate:"required"`
	ImageURL string    `json:"image_url"`
	UserID   int       `json:"user_id"`
	Comments []Comment `json:"comments" gorm:"foreignKey:PostID"`
	Likes    []Like    `json:"likes" gorm:"foreignKey:PostID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

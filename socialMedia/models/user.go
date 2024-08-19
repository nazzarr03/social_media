package models

import "time"

type User struct {
	UserID    int    `json:"user_id" gorm:"primaryKey;autoIncrement"`
	Username  string `json:"username" gorm:"unique" validate:"required"`
	Password  string `json:"password" validate:"required"`
	Email     string `json:"email" gorm:"unique" validate:"required,email"`
	ImageURL  string `json:"image_url"`
	Posts     []Post `json:"posts" gorm:"foreignKey:UserID"`
	Friends   []User `json:"friends" gorm:"many2many:friendships;foreignKey:UserID;joinForeignKey:UserID;JoinReferences:FriendID;References:UserID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

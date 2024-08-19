package models

import "time"

type Friendship struct {
	FriendshipID int  `json:"friendship_id" gorm:"primaryKey;autoIncrement"`
	UserID       int  `json:"user_id"`
	FriendID     int  `json:"friend_id"`
	IsActive     bool `json:"is_active" gorm:"default:false"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

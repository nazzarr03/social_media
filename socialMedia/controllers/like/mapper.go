package like

type LikeDto struct {
	LikeID    int  `json:"like_id"`
	IsLiked   bool `json:"is_liked"`
	UserID    int  `json:"user_id"`
	PostID    int  `json:"post_id"`
	CommentID *int `json:"comment_id"`
}

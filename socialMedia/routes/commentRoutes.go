package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nazzarr03/social-media/controllers/comment"
)

func CommentRoutes(app *fiber.App) {
	app.Post("/posts/:postID/comments", comment.CreateCommentToPost)
	app.Post("posts/:postID/comments/:commentID", comment.CreateCommentToComment)
}

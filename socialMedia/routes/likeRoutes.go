package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nazzarr03/social-media/controllers/like"
)

func LikeRoutes(app *fiber.App) {
	app.Post("/posts/:postID/likes", like.CreateLikeToPost)
	app.Post("posts/:postID/comments/:commentID/likes", like.CreateLikeToComment)
}

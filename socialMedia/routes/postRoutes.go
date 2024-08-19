package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nazzarr03/social-media/controllers/post"
)

func PostRoutes(app *fiber.App) {
	app.Post("/users/userID", post.CreatePost)
	app.Post("/posts/:postID/image", post.CreateImageByPostID)
	app.Put("/posts/:postID", post.UpdatePost)
	app.Put("/:postID/image", post.UpdatePostImage)
	app.Delete("/posts/:postID", post.DeletePost)
}

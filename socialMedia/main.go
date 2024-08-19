package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nazzarr03/social-media/middleware"
	"github.com/nazzarr03/social-media/routes"
)

func main() {
	app := fiber.New()
	app.Use(middleware.LogMiddleware())

	routes.UserRoutes(app)
	routes.PostRoutes(app)
	routes.LikeRoutes(app)
	routes.CommentRoutes(app)

	app.Listen(":8081")
}

package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nazzarr03/social-media/controllers/user"
	"github.com/nazzarr03/social-media/middleware"
)

func UserRoutes(app *fiber.App) {
	app.Post("/register", user.Register)
	app.Post("/login", user.Login)

	app.Use(middleware.Authentication())

	app.Post("/id/profile", user.CreateProfileImage)
	app.Put("/id/profile", user.UpdateProfileImage)
}

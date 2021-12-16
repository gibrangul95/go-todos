package userRoutes

import (
	"github.com/gibrangul95/go-todos/middleware"
	"github.com/gibrangul95/go-todos/internal/handlers/user"
	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(router fiber.Router) {
	user := router.Group("/user")
	user.Post("/signup", userHandler.CreateUser)
	user.Post("/signin", userHandler.LoginUser)
	user.Get("/fetch-access-token", userHandler.GetAccessToken)

	user.Post("/signout", middleware.SecureAuth(), userHandler.Logout)
	user.Get("/user", middleware.SecureAuth(), userHandler.GetUserData)
}
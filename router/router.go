package router

import (
	"github.com/gibrangul95/go-todos/internal/routes/todo"
	"github.com/gibrangul95/go-todos/internal/routes/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func hello(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	}

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api", logger.New())
	
	api.Get("/hello", hello)

	userRoutes.SetupUserRoutes(api)
	todoRoutes.SetupTodoRoutes(api)
}
package main

import (
	"github.com/gibrangul95/go-todos/router"
	"github.com/gibrangul95/go-todos/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"encoding/json"
)

func main() {
	database.ConnectDB()

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	app.Use(cors.New())

	router.SetupRoutes(app)

	// 404 Handler
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
	})

	app.Listen(":3000")
}
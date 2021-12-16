package todoRoutes

import (
	"github.com/gibrangul95/go-todos/internal/handlers/todo"
	"github.com/gibrangul95/go-todos/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupTodoRoutes(router fiber.Router) {
	todo := router.Group("/todo")
	todo.Post("/create", middleware.SecureAuth(), todoHandler.CreateTodo)
	todo.Get("/", middleware.SecureAuth(), todoHandler.GetTodos)
	todo.Get("/:todoId", middleware.SecureAuth(), todoHandler.GetTodo)
	todo.Delete("/:todoId", middleware.SecureAuth(), todoHandler.DeleteTodo)
	todo.Patch("/:todoId", middleware.SecureAuth(), todoHandler.UpdateTodoTitle)
	todo.Patch("/:todoId/check", middleware.SecureAuth(), todoHandler.CheckTodo)
	todo.Patch("/:todoId/uncheck", middleware.SecureAuth(), todoHandler.UncheckTodo)
	todo.Patch("/:todoId/assign", middleware.SecureAuth(), todoHandler.AssignTodo)
}
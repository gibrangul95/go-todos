package todoHandler

import (
	"strings"
	"strconv"
	"github.com/google/uuid"
	"github.com/gibrangul95/go-todos/internal/model"
	"github.com/gibrangul95/go-todos/database"
	"github.com/gofiber/fiber/v2"
	valid "github.com/asaskevich/govalidator"
)

func CreateTodo(c *fiber.Ctx) error {
	userId := c.Locals("id").(string)
	todo := new(model.Todo)

	err := c.BodyParser(todo)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}

	if !valid.IsUnixTime(strconv.Itoa(todo.DueDate)) {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Enter a valid unix timestamp", "data": err})
	}

	todo.ID = uuid.New()
	todo.Owner = uuid.Must(uuid.Parse(userId))
	err = database.DB.Create(&todo).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Created todo", "data": todo})
}

func GetTodo(c *fiber.Ctx) error {
	userId := c.Locals("id").(string)
	var todo model.Todo

	id := c.Params("todoId")
	database.DB.Find(&todo, "id = ? AND owner = ?", id, userId)

	if todo.ID == uuid.Nil {
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No todo present", "data": nil})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "todo Found", "data": todo})
}

func GetTodos(c *fiber.Ctx) error {
	userId := c.Locals("id").(string)
	sortBy := c.Query("sortBy", "created_at")
	sortOrder := c.Query("sort", "asc")
	search := "%" + strings.ToLower(c.Query("search", "")) + "%"
	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil {
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Invalid offset", "data": nil})
	}
	
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil {
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Invalid limit", "data": nil})
	}

	var todo []model.Todo
	database.DB.Order(sortBy + " " + sortOrder).Limit(limit).Offset(offset).Where("LOWER(title) LIKE ?", search).Find(&todo, "owner = ?", userId)
	return c.JSON(fiber.Map{"status": "success", "message": "todos Found", "data": todo})
}

func DeleteTodo(c *fiber.Ctx) error {
	userId := c.Locals("id").(string)
	var todo model.Todo
	id := c.Params("todoId")
	database.DB.Find(&todo, "id = ? AND owner = ?", id, userId)

	if todo.ID == uuid.Nil {
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No todo present", "data": nil})
	}

	err := database.DB.Delete(&todo, "id = ? AND owner = ?", id, userId).Error

	if err != nil {
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Failed to delete todo", "data": nil})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Deleted Todo"})
}

func UpdateTodoTitle(c *fiber.Ctx) error {
	userId := c.Locals("id").(string)
	type updateTodo struct {
			Title    string `json:"title"`
	}

	var todo model.Todo

	id := c.Params("todoId")
	database.DB.Find(&todo, "id = ? AND owner = ?", id, userId)

	if todo.ID == uuid.Nil {
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No todo present", "data": nil})
	}

	var updateTodoData updateTodo
	err := c.BodyParser(&updateTodoData)
	if err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}

	if updateTodoData.Title != "" {
		todo.Title = updateTodoData.Title
	}

	database.DB.Save(&todo)
	return c.JSON(fiber.Map{"status": "success", "message": "Todo Updated", "data": todo})
}

func CheckTodo(c *fiber.Ctx) error {
	userId := c.Locals("id").(string)

	var todo model.Todo

	id := c.Params("todoId")
	database.DB.Find(&todo, "id = ? AND owner = ?", id, userId)

	if todo.ID == uuid.Nil {
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No todo present", "data": nil})
	}
	todo.Completed = true

	database.DB.Save(&todo)
	return c.JSON(fiber.Map{"status": "success", "message": "Todo Completed", "data": todo})
}

func UncheckTodo(c *fiber.Ctx) error {
	userId := c.Locals("id").(string)

	var todo model.Todo

	id := c.Params("todoId")
	database.DB.Find(&todo, "id = ? AND owner = ?", id, userId)

	if todo.ID == uuid.Nil {
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No todo present", "data": nil})
	}
	todo.Completed = false

	database.DB.Save(&todo)
	return c.JSON(fiber.Map{"status": "success", "message": "Todo Completed", "data": todo})
}

func AssignTodo(c *fiber.Ctx) error {
	userId := c.Locals("id").(string)
	type updateTodo struct {
			AssignedTo    string `json:"assignedTo"`
	}

	var todo model.Todo

	id := c.Params("todoId")
	database.DB.Find(&todo, "id = ? AND owner = ?", id, userId)

	if todo.ID == uuid.Nil {
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No todo present", "data": nil})
	}

	var updateTodoData updateTodo
	err := c.BodyParser(&updateTodoData)
	if err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}

	if updateTodoData.AssignedTo != "" {
		if !valid.IsEmail(updateTodoData.AssignedTo) {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Invalid email address ", "data": err})
		}

		if count := database.DB.Where(&model.User{Email: updateTodoData.AssignedTo}).First(new(model.User)).RowsAffected; count > 0 {
			// Send registration email logic here
		} else {
			todo.AssignedTo = updateTodoData.AssignedTo
		}
	}

	database.DB.Save(&todo)
	return c.JSON(fiber.Map{"status": "success", "message": "Todo Assigned", "data": todo})
}
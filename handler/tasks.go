package handler

import (
	"app/database"
	"app/model"
	"fmt"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func CreateList(c *fiber.Ctx) error {
	type CreateListInput struct {
		Name string `json:"name" validate:"required,min=1"`
	}

	var input CreateListInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Review your input",
			"errors":  err.Error(),
		})
	}

	// Validate input
	validate := validator.New()
	if err := validate.Struct(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Validation failed",
			"errors":  err.Error(),
		})
	}

	// Get user ID from token (assume user ID is stored in the token)
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve user ID from token",
			"data":    nil,
		})
	}

	// Create the list for the user
	db := database.DB
	list := model.TaskList{
		UserID: uint(userID),
		Name:   input.Name,
	}

	if err := db.Create(&list).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Couldn't create list",
			"errors":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "List created successfully",
		"data":    list,
	})
}

func UpdateListName(c *fiber.Ctx) error {
	type UpdateListNameInput struct {
		Name string `json:"name" validate:"required,min=1"`
	}

	var input UpdateListNameInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Review your input",
			"errors":  err.Error(),
		})
	}

	// Validate input
	validate := validator.New()
	if err := validate.Struct(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Validation failed",
			"errors":  err.Error(),
		})
	}

	// Get the list ID from the route parameter
	listID := c.Params("list_id")

	// Get user ID from token (assume user ID is stored in the token)
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve user ID from token",
			"data":    nil,
		})
	}

	// Retrieve the list and ensure it belongs to the logged-in user
	db := database.DB
	var list model.TaskList
	if err := db.First(&list, "id = ?", listID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "List not found",
			"errors":  err.Error(),
		})
	}

	if list.UserID != uint(userID) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "You are not authorized to update this list",
			"data":    nil,
		})
	}

	// Update the list name
	list.Name = input.Name
	if err := db.Save(&list).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Couldn't update list name",
			"errors":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "List name updated successfully",
		"data":    list,
	})
}

func AddTaskToList(c *fiber.Ctx) error {
	type AddTaskInput struct {
		Text string `json:"text" validate:"required,min=1"`
	}

	var input AddTaskInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Review your input",
			"errors":  err.Error(),
		})
	}

	// Validate input
	validate := validator.New()
	if err := validate.Struct(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Validation failed",
			"errors":  err.Error(),
		})
	}

	// Get the list ID from the route parameter
	listID := c.Params("list_id")

	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve user ID from token",
			"data":    nil,
		})
	}

	// Verify the list exists
	db := database.DB
	var list model.TaskList
	if err := db.First(&list, "id = ?", listID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "List not found",
			"errors":  err.Error(),
		})
	}

	//Verify the list belongs to the authenticated user
	if list.UserID != uint(userID) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "You are not authorized to add tasks to this list",
			"data":    nil,
		})
	}

	// Create the task for the specified list
	task := model.Task{
		TaskListID: list.ID,
		Text:       input.Text,
	}

	if err := db.Create(&task).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Couldn't create task",
			"errors":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Task added successfully",
		"data":    task,
	})
}

func GetListsForUser(c *fiber.Ctx) error {
	// Get user ID from token (assume user ID is stored in the token)
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve user ID from token",
			"data":    nil,
		})
	}

	// Retrieve all lists for the user, including their tasks
	db := database.DB
	var lists []model.TaskList
	if err := db.Preload("Tasks").Where("user_id = ?", uint(userID)).Find(&lists).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Couldn't retrieve lists",
			"errors":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Lists retrieved successfully",
		"data":    lists,
	})
}

func DeleteList(c *fiber.Ctx) error {
	// Get the list ID from URL parameters
	idStr := c.Params("list_id")

	fmt.Println(idStr)
	if idStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid or missing task list ID",
			"data":    nil,
		})
	}

	// Convert id to uint
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid task list ID format",
			"data":    nil,
		})
	}

	// Get user ID from token
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve user ID from token",
			"data":    nil,
		})
	}

	db := database.DB
	var list model.TaskList

	// Check if the list exists and belongs to the user
	if err := db.First(&list, "id = ?", uint(id)).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "No task list found with the provided ID",
			"errors":  err.Error(),
		})
	}

	if list.UserID != uint(userID) {
		return c.Status(403).JSON(fiber.Map{
			"status":  "error",
			"message": "You are not authorized to delete this task list",
			"data":    nil,
		})
	}

	// Delete the task list
	if err := db.Delete(&list).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Couldn't delete task list",
			"errors":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Task list successfully deleted",
		"data":    nil,
	})
}


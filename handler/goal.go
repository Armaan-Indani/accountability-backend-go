package handler

import (
	"app/database"
	"app/model"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func CreateGoal(c *fiber.Ctx) error {
	type CreateGoalInput struct {
		Name        string    `json:"name" validate:"required,min=1,max=20"`
		Deadline    time.Time `json:"deadline" validate:"required"`
		Description string    `json:"description"`
		What        string    `json:"what"`
		HowMuch     string    `json:"how_much"`
		Resources   string    `json:"resources"`
		Alignment   string    `json:"alignment"`
	}

	var input CreateGoalInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"errors":  err.Error(),
		})
	}

	validate := validator.New()
	if err := validate.Struct(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Validation failed",
			"errors":  err.Error(),
		})
	}

	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve user ID from token",
		})
	}

	goal := model.Goal{
		UserID:      uint(userID),
		Name:        input.Name,
		Deadline:    input.Deadline,
		Description: input.Description,
		What:        input.What,
		HowMuch:     input.HowMuch,
		Resources:   input.Resources,
		Alignment:   input.Alignment,
		Completed:   false,
	}

	db := database.DB
	if err := db.Create(&goal).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Couldn't create goal",
			"errors":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Goal created successfully",
		"data":    goal,
	})
}

func GetGoals(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve user ID from token",
		})
	}

	var goals []model.Goal
	db := database.DB
	if err := db.Where("user_id = ?", uint(userID)).Find(&goals).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Couldn't fetch goals",
			"errors":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Goals retrieved successfully",
		"data":    goals,
	})
}

func DeleteGoal(c *fiber.Ctx) error {
	id := c.Params("goal_id")

	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Goal ID is required",
		})
	}

	db := database.DB
	var goal model.Goal

	// Find the goal first
	if err := db.First(&goal, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Goal not found",
		})
	}

	// Delete the goal
	if err := db.Delete(&goal).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Couldn't delete goal",
			"errors":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Goal deleted successfully",
	})
}

func UpdateGoal(c *fiber.Ctx) error {
	id := c.Params("goal_id")

	var goal model.Goal
	db := database.DB
	if err := db.First(&goal, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Goal not found",
		})
	}

	type UpdateGoalInput struct {
		Name        string    `json:"name"`
		Deadline    time.Time `json:"deadline"`
		Description string    `json:"description"`
		What        string    `json:"what"`
		HowMuch     string    `json:"how_much"`
		Resources   string    `json:"resources"`
		Alignment   string    `json:"alignment"`
	}

	var input UpdateGoalInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"errors":  err.Error(),
		})
	}

	goal.Name = input.Name
	goal.Deadline = input.Deadline
	goal.Description = input.Description
	goal.What = input.What
	goal.HowMuch = input.HowMuch
	goal.Resources = input.Resources
	goal.Alignment = input.Alignment

	if err := db.Save(&goal).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Couldn't update goal",
			"errors":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Goal updated successfully",
		"data":    goal,
	})
}

func ToggleGoalCompletedStatus(c *fiber.Ctx) error {
	id := c.Params("goal_id")

	var goal model.Goal
	db := database.DB
	if err := db.First(&goal, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Goal not found",
		})
	}

	goal.Completed = !goal.Completed
	if err := db.Save(&goal).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Couldn't toggle goal completed status",
			"errors":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Goal completed status toggled",
		"data":    goal,
	})
}

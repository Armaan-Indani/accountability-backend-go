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
	type HabitInput struct {
		Name      string `json:"name"`
		Frequency string `json:"frequency"`
	}
	type CreateGoalInput struct {
		Name        string       `json:"name" validate:"required,min=1"`
		Deadline    time.Time    `json:"deadline"`
		Description string       `json:"description"`
		What        string       `json:"what"`
		HowMuch     string       `json:"how_much"`
		Resources   string       `json:"resources"`
		Alignment   string       `json:"alignment"`
		Subgoals    []string     `json:"subgoals"` // List of subgoal names
		Habits      []HabitInput `json:"habits"`   // List of habits with names and frequency
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

	db := database.DB

	// Create goal
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

	// Add subgoals
	for _, subgoalName := range input.Subgoals {
		goal.Subgoals = append(goal.Subgoals, model.Subgoal{
			Name:      subgoalName,
			Completed: false,
		})
	}

	// Add habits
	for _, habit := range input.Habits {
		goal.Habits = append(goal.Habits, model.Habit{
			Name:      habit.Name,
			Frequency: habit.Frequency,
		})
	}

	// Save to DB
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
	if err := db.Where("user_id = ?", uint(userID)).Preload("Subgoals").Preload("Habits").Find(&goals).Error; err != nil {
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
	goal_id := c.Params("goal_id")

	if goal_id == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Goal ID is required",
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
	var goal model.Goal

	// Find the goal first
	if err := db.First(&goal, goal_id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Goal not found",
		})
	}

	// Check if the goal belongs to the user
	if goal.UserID != uint(userID) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "You are not authorized to delete this goal",
			"data":    nil,
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

	var goal model.Goal
	db := database.DB
	if err := db.Preload("Subgoals").Preload("Habits").First(&goal, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Goal not found",
		})
	}

	// Check if the goal belongs to the user
	if goal.UserID != uint(userID) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "You are not authorized to update this goal",
			"data":    nil,
		})
	}

	type HabitInput struct {
		Name      string `json:"name"`
		Frequency string `json:"frequency"`
	}

	type UpdateGoalInput struct {
		Name        string       `json:"name" validate:"required,min=1"`
		Deadline    time.Time    `json:"deadline"`
		Description string       `json:"description"`
		What        string       `json:"what"`
		HowMuch     string       `json:"how_much"`
		Resources   string       `json:"resources"`
		Alignment   string       `json:"alignment"`
		Subgoals    []string     `json:"subgoals"` // List of subgoal names
		Habits      []HabitInput `json:"habits"`   // List of habits with names and frequency
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

	// Clear existing subgoals and add new ones
	if input.Subgoals != nil {
		db.Where("goal_id = ?", goal.ID).Delete(&model.Subgoal{})
		for _, subgoalName := range input.Subgoals {
			goal.Subgoals = append(goal.Subgoals, model.Subgoal{
				GoalID:    goal.ID,
				Name:      subgoalName,
				Completed: false,
			})
		}
	}

	// Clear existing habits and add new ones
	if input.Habits != nil {
		db.Where("goal_id = ?", goal.ID).Delete(&model.Habit{})
		for _, habit := range input.Habits {
			goal.Habits = append(goal.Habits, model.Habit{
				GoalID:    goal.ID,
				Name:      habit.Name,
				Frequency: habit.Frequency,
			})
		}
	}

	// Save changes
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

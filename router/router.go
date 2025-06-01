package router

import (
	"app/handler"
	"app/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// SetupRoutes setup router api
func SetupRoutes(app *fiber.App) {
	// Middleware
	api := app.Group("/api", logger.New())
	api.Get("/", handler.Hello)

	// Auth
	auth := api.Group("/auth")
	auth.Post("/login", handler.Login)
	auth.Post("/signup", handler.Register)

	// User
	user := api.Group("/user")
	user.Get("/:id", middleware.Protected(), handler.GetUser)
	// user.Post("/register", handler.CreateUser)
	user.Patch("/:id", middleware.Protected(), handler.UpdateUser)
	user.Delete("/:id", middleware.Protected(), handler.DeleteUser)

	//TaskLists
	taskList := api.Group("/tasklist")
	taskList.Get("/", middleware.Protected(), handler.GetListsForUser)
	taskList.Post("/", middleware.Protected(), handler.CreateList)
	taskList.Patch("/:list_id", middleware.Protected(), handler.UpdateListName)
	taskList.Delete("/:list_id", middleware.Protected(), handler.DeleteList)

	//Tasks
	task := api.Group("/task")
	task.Post("/:list_id", middleware.Protected(), handler.AddTaskToList)
	task.Delete("/:task_id", middleware.Protected(), handler.DeleteTask)
	// TODO: Change to PUT - backend and frontend
	task.Patch("/:task_id", middleware.Protected(), handler.UpdateTask)
	task.Patch("/:task_id/toggle", middleware.Protected(), handler.ToggleTask)

	//Goals
	goal := api.Group("/goal")
	goal.Post("/", middleware.Protected(), handler.CreateGoal)
	goal.Get("/", middleware.Protected(), handler.GetGoals)
	goal.Put("/:goal_id", middleware.Protected(), handler.UpdateGoal)
	goal.Delete("/:goal_id", middleware.Protected(), handler.DeleteGoal)
	goal.Patch("/:goal_id/toggle", middleware.Protected(), handler.ToggleGoalCompletedStatus)
	goal.Patch("/:goal_id/:subgoal_id/toggle", middleware.Protected(), handler.ToggleSubgoalCompletedStatus)
}

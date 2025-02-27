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

	// User
	user := api.Group("/user")
	user.Get("/:id", middleware.Protected(), handler.GetUser)
	user.Post("/register", handler.CreateUser)
	user.Patch("/:id", middleware.Protected(), handler.UpdateUser)
	user.Delete("/:id", middleware.Protected(), handler.DeleteUser)


	//TaskLists
	taskList := api.Group("/tasklist")
	taskList.Get("/", middleware.Protected(), handler.GetListsForUser)
	taskList.Post("/", middleware.Protected(), handler.CreateList)
	taskList.Patch("/:list_id", middleware.Protected(), handler.UpdateListName)
	taskList.Delete("/:list_id", middleware.Protected(), handler.DeleteList)
	
	task := api.Group("/task")
	task.Post("/:list_id", middleware.Protected(), handler.AddTaskToList)
	task.Delete("/:task_id", middleware.Protected(), handler.DeleteTask)
	task.Patch("/:task_id", middleware.Protected(), handler.UpdateTask)
	task.Patch("/:task_id/toggle", middleware.Protected(), handler.ToggleTask)

}

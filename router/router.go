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

	// Product
	product := api.Group("/product")
	product.Get("/", handler.GetAllProducts)
	product.Get("/:id", handler.GetProduct)
	product.Post("/", middleware.Protected(), handler.CreateProduct)
	product.Delete("/:id", middleware.Protected(), handler.DeleteProduct)

	//Tasks
	task := api.Group("/task")
	task.Get("/", middleware.Protected(), handler.GetListsForUser)
	task.Post("/", middleware.Protected(), handler.CreateList)
	task.Patch("/:list_id", middleware.Protected(), handler.UpdateListName)
	task.Delete("/:list_id", middleware.Protected(), handler.DeleteList)

	task.Post("/:list_id", middleware.Protected(), handler.AddTaskToList)
}

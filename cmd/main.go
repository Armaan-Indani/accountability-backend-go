package main

import (
	"log"

	"app/database"
	"app/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "App Name",
	})
	
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",       // Allow requests from the frontend
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS", // Allow these HTTP methods
		AllowHeaders: "Content-Type, Authorization", // Allow these headers
	}))

	database.ConnectDB()

	router.SetupRoutes(app)
	log.Fatal(app.Listen(":5000"))
}

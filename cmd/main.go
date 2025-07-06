package main

import (
	"log"
	"os"

	"app/database"
	"app/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Connect to database first, before creating multiple processes
	database.ConnectDB()

	app := fiber.New(fiber.Config{
		Prefork:       false, // Disable prefork in Docker environment
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "App Name",
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,https://accountability-project-frontend.vercel.app",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
	}))

	router.SetupRoutes(app)
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	
	log.Fatal(app.Listen(":" + port))
}

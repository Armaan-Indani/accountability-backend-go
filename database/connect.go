package database

import (
	"fmt"
	"strconv"

	"app/config"
	"app/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectDB connect to db
func ConnectDB() {
	var err error
	
	// Get database configuration from environment
	dbHost := config.Config("DB_HOST")
	dbUser := config.Config("DB_USER")
	dbPassword := config.Config("DB_PASSWORD")
	dbName := config.Config("DB_NAME")
	p := config.Config("DB_PORT")
	
	// Validate required environment variables
	if dbHost == "" || dbUser == "" || dbPassword == "" || dbName == "" || p == "" {
		panic("Missing required database environment variables (DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT)")
	}
	
	port, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		panic(fmt.Sprintf("failed to parse database port '%s': %v", p, err))
	}

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost,
		port,
		dbUser,
		dbPassword,
		dbName,
	)
	
	fmt.Printf("Connecting to database: host=%s port=%d dbname=%s user=%s\n", dbHost, port, dbName, dbUser)
	
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}

	fmt.Println("Connection Opened to Database")
	
	// Run migrations
	err = DB.AutoMigrate(&model.User{}, &model.TaskList{}, &model.Task{}, &model.Goal{}, &model.Subgoal{}, &model.Habit{})
	if err != nil {
		panic(fmt.Sprintf("failed to run database migrations: %v", err))
	}
	
	fmt.Println("Database Migrated")
}

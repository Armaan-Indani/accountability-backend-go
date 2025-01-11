package handler

import (
	"errors"
	"log"
	"net/mail"
	"time"

	"app/config"
	"app/database"
	"app/model"

	"gorm.io/gorm"

	// "github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// CheckPasswordHash compare password with hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	log.Println(hash, string(hashedPassword), "haaaash")
	return err == nil
}

func getUserByEmail(e string) (*model.User, error) {
	db := database.DB
	var user model.User
	if err := db.Where(&model.User{Email: e}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func getUserByUsername(u string) (*model.User, error) {
	db := database.DB
	var user model.User
	if err := db.Where(&model.User{Username: u}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func valid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// Login get user and password
func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Identity string `json:"identity"`
		Password string `json:"password"`
	}
	type UserData struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	input := new(LoginInput)
	var ud UserData

	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on login request", "errors": err.Error()})
	}

	identity := input.Identity
	pass := input.Password
	userModel, err := new(model.User), *new(error)

	if valid(identity) {
		userModel, err = getUserByEmail(identity)
	} else {
		userModel, err = getUserByUsername(identity)
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Internal Server Error", "data": err})
	} else if userModel == nil {
		CheckPasswordHash(pass, "")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid identity or password", "data": err})
	} else {
		ud = UserData{
			ID:       userModel.ID,
			Username: userModel.Username,
			Password: userModel.Password,
		}
	}

	if !CheckPasswordHash(pass, ud.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid identity or password", "data": nil})
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = ud.Username
	claims["user_id"] = ud.ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	t, err := token.SignedString([]byte(config.Config("SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Success login", "data": t})
}

// func Register(c *fiber.Ctx) error {
// 	type RegisterInput struct {
// 		Username string `json:"username" validate:"required,min=3,max=50"`
// 		Email    string `json:"email" validate:"required,email"`
// 		Password string `json:"password" validate:"required,min=6,max=50"`
// 		Names    string `json:"names"`
// 	}

// 	input := new(RegisterInput)
// 	if err := c.BodyParser(input); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"status":  "error",
// 			"message": "Invalid input",
// 			"errors":  err.Error(),
// 		})
// 	}

// 	// Validate input using Fiber's built-in validator or a third-party library
// 	validate := validator.New()
// 	if err := validate.Struct(input); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"status":  "error",
// 			"message": "Validation failed",
// 			"errors":  err.Error(),
// 		})
// 	}

// 	// Check if the username or email already exists
// 	var existingUser model.User
// 	if err := database.DB.Where("username = ? OR email = ?", input.Username, input.Email).First(&existingUser).Error; err == nil {
// 		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
// 			"status":  "error",
// 			"message": "Username or email already exists",
// 		})
// 	}

// 	// Hash the password
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"status":  "error",
// 			"message": "Failed to hash password",
// 			"errors":  err.Error(),
// 		})
// 	}

// 	// Create the user
// 	newUser := model.User{
// 		Username: input.Username,
// 		Email:    input.Email,
// 		Password: string(hashedPassword),
// 		Names:    input.Names,
// 	}

// 	if err := database.DB.Create(&newUser).Error; err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"status":  "error",
// 			"message": "Failed to create user",
// 			"errors":  err.Error(),
// 		})
// 	}

// 	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
// 		"status":  "success",
// 		"message": "User registered successfully",
// 		"data": fiber.Map{
// 			"id":       newUser.ID,
// 			"username": newUser.Username,
// 			"email":    newUser.Email,
// 		},
// 	})
// }

package handler

import (
	"errors"
	// "log"

	// "net/mail"
	"time"

	"app/config"
	"app/database"
	"app/model"

	"gorm.io/gorm"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// CheckPasswordHash compare password with hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
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

func getUserByUsername(username string) (*model.User, error) {
	db := database.DB
	var user model.User
	if err := db.Where(&model.User{Username: username}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Login get user and password
func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
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

	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body", "errors": err.Error()})
	}

	email := input.Email
	pass := input.Password
	userModel, err := new(model.User), *new(error)

	userModel, err = getUserByEmail(email)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Internal Server Error", "data": err})
	} else if userModel == nil {
		CheckPasswordHash(pass, "")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid identity or password", "data": err})
	} else {
		ud = UserData{
			ID:       userModel.ID,
			Password: userModel.Password,
		}
	}

	if !CheckPasswordHash(pass, ud.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid identity or password", "data": nil})
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = ud.ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	t, err := token.SignedString([]byte(config.Config("SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Success login", "data": t})
}

// Register creates a new user
func Register(c *fiber.Ctx) error {
	type RegisterInput struct {
		Username   string `json:"username" validate:"required"`
		Email      string `json:"email" validate:"required,email"`
		Password   string `json:"password" validate:"required,min=6"`
		Occupation string `json:"occupation"`
	}

	input := new(RegisterInput)

	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on register request", "errors": err.Error()})
	}

	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body", "errors": err.Error()})
	}

	db := database.DB
	existingUser, err := getUserByEmail(input.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Internal Server Error", "data": err})
	}
	if existingUser != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "error", "message": "Email already in use"})
	}

	existingUser2, err := getUserByUsername(input.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Internal Server Error", "data": err})
	}
	if existingUser2 != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "error", "message": "Username already in use"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not hash password", "data": err})
	}

	newUser := model.User{
		Username:   input.Username,
		Email:      input.Email,
		Password:   string(hashedPassword),
		Occupation: input.Occupation,
	}

	if err := db.Create(&newUser).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not create user", "data": err})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "message": "User registered successfully"})
}

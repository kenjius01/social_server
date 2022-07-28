package controllers

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/kenjius01/social-sever/database"
	"github.com/kenjius01/social-sever/models"
	"golang.org/x/crypto/bcrypt"
)

// Create new user
func Register(c *fiber.Ctx) error {
	user := models.User{}

	err := c.BodyParser(&user)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	checkExist := models.User{}
	database.DB.First(&checkExist, "username = ?", user.Username)
	if checkExist.ID != 0 {
		return c.Status(404).JSON("This username has been exist!")
	}
	password, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	user.Password = string(password)
	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(200).JSON(CreateResponseUser(user))
}

func Login(c *fiber.Ctx) error {
	data := models.User{}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	user := models.User{}

	database.DB.Where("username = ?", data.Username).First(&user)
	if user.ID == 0 {
		return c.Status(fiber.StatusPreconditionFailed).JSON(fiber.Map{
			"Error":   "User not found!",
			"success": false,
		})
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		return c.Status(402).JSON(fiber.Map{
			"Error":   "Incorrect Password!",
			"success": false,
		})

	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    user.Username,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	secretKey := os.Getenv("SECRET_KEY")

	token, err := claims.SignedString([]byte(secretKey))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Could not login!",
			"success": false,
		})

	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 99999),
		HTTPOnly: false,
	}
	c.Cookie(&cookie)
	return c.JSON(CreateResponseUser(user))

}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"success": true,
	})

}

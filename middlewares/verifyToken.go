package middlewares

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/kenjius01/social-sever/database"
	"github.com/kenjius01/social-sever/models"
)

func VerifyUser(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	SecretKey := os.Getenv("SECRET_KEY")
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if token == nil {
		return c.Status(401).JSON(fiber.Map{
			"message": "You are not authenticated!",
		})
	}

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthenticate",
		})
	}
	claims := token.Claims.(*jwt.StandardClaims)
	var user models.User

	database.DB.Where("username = ?", claims.Issuer).First(&user)
	id, _ := c.ParamsInt("id")
	userFollow := models.Follower{}
	if err := c.BodyParser(&userFollow); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if user.ID == (id) || user.IsAdmin {
		return c.Next()
	} else {
		return c.Status(403).JSON("you are not authorized!")
	}

}

func VeryfyAdmin(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	SecretKey := os.Getenv("SECRET_KEY")
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if token == nil {
		return c.Status(401).JSON(fiber.Map{
			"message": "You are not authenticated!",
		})
	}

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthenticate",
		})
	}
	claims := token.Claims.(*jwt.StandardClaims)
	var user models.User

	database.DB.Where("username = ?", claims.Issuer).First(&user)
	if user.IsAdmin {
		return c.Next()
	} else {
		return c.Status(403).JSON("You are not Admin!")
	}
}

func VerifyUserPost(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	SecretKey := os.Getenv("SECRET_KEY")
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if token == nil {
		return c.Status(401).JSON(fiber.Map{
			"message": "You are not authenticated!",
		})
	}

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthenticate",
		})
	}
	claims := token.Claims.(*jwt.StandardClaims)
	var user models.User

	database.DB.Where("username = ?", claims.Issuer).First(&user)
	post := models.Post{}
	if err := c.BodyParser(&post); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if user.IsAdmin || post.UserId == user.ID {
		return c.Next()
	} else {
		return c.Status(403).JSON("you are not authorized!")
	}

}

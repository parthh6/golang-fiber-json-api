package main

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func AuthHandlers(route fiber.Router, db *gorm.DB) {
	route.Post("/register", func(c *fiber.Ctx) error {
		user := &User{
			Username: c.FormValue("username"),
			Password: c.FormValue("password"),
		}

		if user.Username == "" || user.Password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Username and Password are required",
			})
		}
		hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		user.Password = string(hashed)
		db.Create(user)
		//Generating token from utils.go
		token, err := GenerateToken(user)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		c.Cookie(&fiber.Cookie{
			Name:     "jwt",
			Value:    "token",
			HTTPOnly: !c.IsFromLocal(),
			Secure:   !c.IsFromLocal(),
			MaxAge:   3600 * 24 * 7, //7 days
		})
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"token": token,
		})

	})
	route.Post("/login", func(c *fiber.Ctx) error {
		dbUser := new(User)
		authUser := &User{
			Username : c.FormValue("username"),
			Password : c.FormValue("password"),
		}
		if authUser.Username == "" || authUser.Password == ""{
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error" : "Username and Password are required",
			})
		}
		//Username is checked if it exists in the DB
		db.Where("username = ?",authUser.Username).First(dbUser)
		if dbUser.ID == 0{
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error" : "User not found",
			})
		}
		//Compare the hashed and entered password
		err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password),[]byte(authUser.Password))
		if err != nil{
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error" : "Invalid Credentials",
			})
		}
		token,err := GenerateToken(authUser)
		if err != nil{
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error" : err.Error(),
			})
		}
		c.Cookie(&fiber.Cookie{
			Name:     "jwt",
			Value:    "token",
			HTTPOnly: !c.IsFromLocal(),
			Secure:   !c.IsFromLocal(),
			MaxAge:   3600 * 24 * 7, //7 days
		})
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"token": token,
		})
	})
}

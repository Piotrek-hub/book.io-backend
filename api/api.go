package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	db "github.com/piotrek-hub/book.io-backend/db"
)

type Resp struct {
	UserKey string `form:"user_key"`
}

func StartApi() {
	app := fiber.New()
	app.Use(cors.New())

	app.Post("/login", func(c *fiber.Ctx) error {
		u := new(db.User)
		if err := c.BodyParser(u); err != nil {
			return err
		}

		userKey := db.Login(u.Login, u.Password)
		return c.JSON(fiber.Map{
			"succes":   200,
			"user_key": userKey,
		})
	})

	app.Post("/register", func(c *fiber.Ctx) error {
		u := new(db.User)
		if err := c.BodyParser(u); err != nil {
			return err
		}

		userKey, info := db.Register(u.Login, u.Password)
		return c.JSON(fiber.Map{
			"succes":  200,
			"userKey": userKey,
			"info":    info,
		})
	})

	app.Post("/addBook", func(c *fiber.Ctx) error {
		bookRequest := new(db.BookRequest)
		if err := c.BodyParser(bookRequest); err != nil {
			return err
		}
		if bookRequest.UserKey == "" {
			return c.SendString("Provide user key")
		}
		info := db.AddBook(*bookRequest)
		return c.JSON(fiber.Map{
			"succes": 200,
			"info":   info,
		})
	})

	app.Post("/setBookStatus", func(c *fiber.Ctx) error {
		bookRequest := new(db.BookRequest)
		if err := c.BodyParser(bookRequest); err != nil {
			return err
		}
		if bookRequest.UserKey == "" {
			return c.SendString("Provide user key")
		}

		info := db.SetBookStatus(*bookRequest)
		return c.JSON(fiber.Map{
			"succes": 200,
			"info":   info,
		})
	})

	app.Post("/deleteBook", func(c *fiber.Ctx) error {
		bookRequest := new(db.BookRequest)
		if err := c.BodyParser(bookRequest); err != nil {
			return err
		}
		if bookRequest.UserKey == "" {
			return c.SendString("Provide user key")
		}

		info := db.DeleteBook(*bookRequest)
		return c.JSON(fiber.Map{
			"succes": 200,
			"info":   info,
		})
	})

	app.Post("/getBooks", func(c *fiber.Ctx) error {
		bookRequest := new(db.BookRequest)

		if err := c.BodyParser(bookRequest); err != nil {
			fmt.Println(err)
			return c.Status(404).JSON(&fiber.Map{
				"success": false,
				"error":   "Not provided user key",
			})
		}
		fmt.Println("Resp", bookRequest)
		books := db.GetBooks(bookRequest.Username)

		return c.JSON(fiber.Map{
			"succes": 200,
			"books":  books,
		})
	})

	app.Get("/getUsers",func(c *fiber.Ctx) error {

		users := db.GetUsers()

		return c.JSON(fiber.Map{
			"succes": 200,
			"users":  users,
		})
	})

	app.Listen(":3000")
}

package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	db "github.com/piotrek-hub/book.io-backend/db"
)

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
			"status":   200,
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
			"status":  200,
			"userKey": userKey,
			"info":    info,
		})
	})

	// Add Book
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
			"status": 200,
			"info":   info,
		})
	})

	// Set Book Status
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
			"status": 200,
			"info":   info,
		})
	})

	// Delete Book
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
			"status": 200,
			"info":   info,
		})
	})

	// Fetch Books
	app.Get("/getBooks", func(c *fiber.Ctx) error {
		books := db.GetBooks()
		fmt.Println(books)
		return c.JSON(fiber.Map{
			"status": 200,
			"books":  books,
		})
	})

	app.Listen(":3000")
}

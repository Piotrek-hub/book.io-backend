package api

import (
	"github.com/gofiber/fiber/v2"
	db "github.com/piotrek-hub/book.io-backend/db"
)

// POST Log In (username, password)
// POST Register (username, password)

// POST 1. Add Book (name, date ended, pages)
// POST 2. Set Book Status (name, status)
// POST 3. Delete Book (name)
// GET 4. Fetch books

func StartApi() {
	app := fiber.New()

	// Log In
	app.Post("/login", func(c *fiber.Ctx) error {
		// Get Headers
		u := new(db.User)
		if err := c.BodyParser(u); err != nil {
			return err
		}

		// Returun userKey
		userKey := db.Login(u.Login, u.Password)
		return c.JSON(fiber.Map{
			"status":   200,
			"user_key": userKey,
		})
	})

	// Register
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
	app.Delete("/deleteBook", func(c *fiber.Ctx) error {
		bookRequest := new(db.BookRequest)
		if err := c.BodyParser(bookRequest); err != nil {
			return err
		}
		if bookRequest.UserKey == "" {
			return c.SendString("Provide user key")
		}
		return c.Send(c.Body())
	})

	// Fetch Books
	app.Get("/getBooks", func(c *fiber.Ctx) error {
		return c.SendString("Get Books Page")
	})

	app.Listen(":3000")
}

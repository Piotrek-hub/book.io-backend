package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/piotrek-hub/book.io-backend/db"
	"github.com/piotrek-hub/book.io-backend/utils"
)

type Resp struct {
	UserKey string `form:"user_key"`
}

func login(c *fiber.Ctx) error {
	u := new(db.User)
	if err := c.BodyParser(u); err != nil {
		return err
	}

	userKey := db.Login(u.Login, u.Password)
	return c.JSON(fiber.Map{
		"succes":   200,
		"user_key": userKey,
	})
}

func register(c *fiber.Ctx) error {
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
}

func addBook(c *fiber.Ctx) error {
	bookRequest := new(utils.BookRequest)
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
}

func setBookStatus(c *fiber.Ctx) error {
	bookRequest := new(utils.BookRequest)
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
}

func deleteBook(c *fiber.Ctx) error {
	bookRequest := new(utils.BookRequest)
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
}

func getBooks(c *fiber.Ctx) error {
	bookRequest := new(utils.BookRequest)

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
}

func getUsers(c *fiber.Ctx) error {
	users := db.GetUsers()

	return c.JSON(fiber.Map{
		"succes": 200,
		"users":  users,
	})
}

func StartApi() {
	app := fiber.New()
	app.Use(cors.New())

	app.Post("/login", login)
	app.Post("/register", register)
	app.Post("/addBook", addBook)
	app.Post("/setBookStatus", setBookStatus)
	app.Post("/deleteBook", deleteBook)
	app.Post("/getBooks", getBooks)
	app.Get("/getUsers",getUsers)

	app.Listen(":3000")
}

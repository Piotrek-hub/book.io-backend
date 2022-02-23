package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/piotrek-hub/book.io-backend/db"
	"github.com/piotrek-hub/book.io-backend/utils"
	"log"
)


func login(c *fiber.Ctx) error {
	u := new(db.User)
	if err := c.BodyParser(u); err != nil {
		return err
	}

	utils.LogRequest[db.User]("Login", *u)
	token, err := db.Login(u.Login, u.Password)
	if err != nil {
		return c.JSON(fiber.Map{
			"success": false,
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success":   true,
		"token": token,
	})
}

func register(c *fiber.Ctx) error {
	u := new(db.User)
	if err := c.BodyParser(u); err != nil {
		return err
	}

	utils.LogRequest[db.User]("Register", *u)
	token, err := db.Register(u.Login, u.Password)
	if err != nil {
		return c.JSON(fiber.Map{
			"success":  false,
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success":  true,
		"token": token,
	})
}

func addBook(c *fiber.Ctx) error {
	bookRequest := new(utils.BookRequest)
	if err := c.BodyParser(bookRequest); err != nil {
		return err
	}

	utils.LogRequest[utils.BookRequest]("AddBook", *bookRequest)
	err := db.AddBook(*bookRequest)
	if err != nil {
		return c.JSON(fiber.Map{
			"success": false,
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
	})
}

func setBookStatus(c *fiber.Ctx) error {
	bookRequest := new(utils.BookRequest)
	if err := c.BodyParser(bookRequest); err != nil {
		return err
	}

	utils.LogRequest[utils.BookRequest]("SetBookStatus", *bookRequest)
	err := db.SetBookStatus(*bookRequest)
	if err != nil {
		return c.JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
	})
}

func deleteBook(c *fiber.Ctx) error {
	bookRequest := new(utils.BookRequest)
	if err := c.BodyParser(bookRequest); err != nil {
		return err
	}

	utils.LogRequest[utils.BookRequest]("DeleteBook", *bookRequest)
	err := db.DeleteBook(*bookRequest)
	if err != nil {
		return c.JSON(fiber.Map{
			"success": false,
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
	})
}

func getBooks(c *fiber.Ctx) error {
	bookRequest := new(utils.BookRequest)
	if err := c.BodyParser(bookRequest); err != nil {
		return err
	}

	utils.LogRequest[utils.BookRequest]("GetBooks", *bookRequest)
	books, err := db.GetBooks(bookRequest.Username)
	if err != nil {
		return c.JSON(fiber.Map{
			"success": false,
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"books":  books,
	})
}

func getUsers(c *fiber.Ctx) error {
	users := db.GetUsers()

	utils.LogRequest[[]string]("GetUsers", users)
	return c.JSON(fiber.Map{
		"success": true,
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

	err := app.Listen(":3000")
	if err != nil {
		log.Fatal(err)
	}
}

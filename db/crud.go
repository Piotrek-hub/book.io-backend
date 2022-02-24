package db

import (
	"context"
	"errors"
	"github.com/piotrek-hub/book.io-backend/utils"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

func Register(login string, password string) (string, error) {
	ctx, client, coll := connect("users")
	defer client.Disconnect(*ctx)

	_, exists := utils.CheckIfUserExists(bson.D{{"login", login}, {"password", password}}, coll)

	if !exists {
		token, err := utils.GenerateToken(login)
		if err != nil {
			return "", err
		}
		doc := bson.D{{"login", login}, {"password", password}, {"token", token}}
		_, err = coll.InsertOne(context.TODO(), doc)

		if err != nil {
			return "", errors.New("Error during register")
		}
		return token, nil

	} else {
		return "", errors.New("User already exisits")
	}

}

func Login(login string, password string) (string, error) {
	ctx, client, coll := connect("users")
	defer client.Disconnect(*ctx)

	token, userExisits := utils.CheckIfUserExists(bson.D{{"login", login}, {"password", password}}, coll)

	if userExisits {
		return token, nil
	}
	return "", errors.New("User doesnt exists")
}

func AddBook(bookRequest utils.BookRequest) error {
	if bookRequest.Token == "" {
		return errors.New("User key not provided")
	}
	if len(bookRequest.Username) == 0 {
		return errors.New("Username not provided")
	}

	ctx, client, books := connect("books")
	defer client.Disconnect(*ctx)

	usersCtx, usersClient, users := connect("users")
	defer usersClient.Disconnect(*usersCtx)

	_, userExists := utils.CheckIfUserExists(bson.D{{"token", bookRequest.Token}}, users)
	bookExists := utils.CheckIfBookExists(bookRequest, books)

	if bookExists {
		return errors.New("Book already exists")
	}
	if !userExists  {
		return errors.New("User doesnt userExists or book already!")
	}

	doc := utils.InitBookDoc(bookRequest, bookRequest.Token, bookRequest.Username)

	_, err := books.InsertOne(context.TODO(), doc)
	if err != nil {
		return errors.New("Error while adding book to db")
	}

	return nil
}

func SetBookStatus(bookRequest utils.BookRequest) (error) {
	if bookRequest.Token == "" {
		return errors.New("Provide user key")
	}

	ctx, client, coll := connect("books")
	defer client.Disconnect(*ctx)

	bookExists := utils.CheckIfBookExists(bookRequest, coll)

	if !bookExists {
		return errors.New("Book doesnt exists")
	}

	_, err := coll.UpdateOne(
		*ctx,
		bson.M{"title": bookRequest.Title, "token": bookRequest.Token},
		bson.D{
			{"$set", bson.D{{"Status", bookRequest.Status}}},
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func DeleteBook(bookRequest utils.BookRequest) (error) {
	if bookRequest.Token == "" {
		return errors.New("User key not provided")
	}

	ctx, client, coll := connect("books")
	defer client.Disconnect(*ctx)

	if !utils.CheckIfBookExists(bookRequest, coll) {
		return errors.New("Book doesnt exists")
	}

	result, err := coll.DeleteOne(*ctx, bson.M{"title": bookRequest.Title, "token": bookRequest.Token})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(result)
	return nil
}

func GetBooks(username string) ([]Book, error) {
	var books []Book

	ctx, client, coll := connect("books")
	defer client.Disconnect(*ctx)

	result, err := coll.Find(context.TODO(), bson.M{"username": username})
	if err != nil {
		return books, err
	}

	for result.Next(*ctx) {
		var res bson.M
		err := result.Decode(&res)
		if err != nil {
			log.Fatal(err)
		} else {
			var book Book

			bsonBytes, _ := bson.Marshal(res)
			bson.Unmarshal(bsonBytes, &book)

			books = append(books, book)
		}
	}

	return books, nil
}

func GetUsers() []string {
	type Res struct {
		Login string
	}

	ctx, client, coll := connect("users")
	defer client.Disconnect(*ctx)

	result, err := coll.Find(context.TODO(), bson.M{})

	var users []string

	if err != nil {
		defer result.Close(*ctx)
	} else {
		for result.Next(*ctx) {
			var res bson.M
			err := result.Decode(&res)
			if err != nil {
				log.Fatal(err)
			} else {
				var user Res

				bsonBytes, _ := bson.Marshal(res)
				bson.Unmarshal(bsonBytes, &user)

				users = append(users, user.Login)
			}
		}
	}

	return users
}

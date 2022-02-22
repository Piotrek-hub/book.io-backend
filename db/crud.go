package db

import (
	"context"
	"fmt"
	"github.com/piotrek-hub/book.io-backend/utils"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

func Register(login string, password string) (string, string) {

	ctx, client, coll := connect("users")
	defer client.Disconnect(*ctx)

	_, exisits := utils.CheckIfUserExists(bson.D{{"Login", login}, {"Password", password}}, coll)

	if !exisits {
		userKey := utils.DeriveUserKey(login, password)
		doc := bson.D{{"Login", login}, {"Password", password}, {"UserKey", userKey}}
		_, err := coll.InsertOne(context.TODO(), doc)

		if err != nil {
			return "", "Error during register"
		}
		return userKey, "User registered successfully"

	} else {
		return "", "User already exisits"
	}

}

func Login(login string, password string) string {
	ctx, client, coll := connect("users")
	defer client.Disconnect(*ctx)

	userKey, userExisits := utils.CheckIfUserExists(bson.D{{"Login", login}, {"Password", password}}, coll)

	if userExisits {
		return userKey
	} else {
		return "User doesnt exists"
	}
	return ""
}

func AddBook(bookRequest utils.BookRequest) string {
	ctx, client, coll := connect("books")
	defer client.Disconnect(*ctx)

	usersCtx, usersClient, users := connect("users")
	defer usersClient.Disconnect(*usersCtx)

	// Check if user exists
	fmt.Println(bookRequest)
	userKey := bookRequest.UserKey
	username := bookRequest.Username
	_, exisits := utils.CheckIfUserExists(bson.D{{"UserKey", userKey}}, users)

	if !exisits {
		return "Caller doesnt exists"
	}

	if len(username) == 0 {
		return "provide username"
	}

	doc := utils.InitBookDoc(bookRequest, userKey, username)

	_, err := coll.InsertOne(context.TODO(), doc)

	if err != nil {
		return "Error while adding book to db"
	}

	return "Book added successfully"
}

func SetBookStatus(bookRequest utils.BookRequest) string {
	ctx, client, coll := connect("books")
	defer client.Disconnect(*ctx)

	bookExists := utils.CheckIfBookExists(bookRequest, coll)

	if !bookExists {
		return "Book doesnt exisits"
	}

	_, err := coll.UpdateOne(
		*ctx,
		bson.M{"Title": bookRequest.Title, "UserKey": bookRequest.UserKey},
		bson.D{
			{"$set", bson.D{{"Status", bookRequest.Status}}},
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	return "Status changed successfully"
}

func DeleteBook(bookRequest utils.BookRequest) string {
	ctx, client, coll := connect("books")
	defer client.Disconnect(*ctx)

	if !utils.CheckIfBookExists(bookRequest, coll) {
		return "Book doesnt exists"
	}

	result, err := coll.DeleteOne(*ctx, bson.M{"Title": bookRequest.Title, "UserKey": bookRequest.UserKey})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
	return "Deleted successfully"
}

func GetBooks(username string) []Book {
	ctx, client, coll := connect("books")
	defer client.Disconnect(*ctx)

	result, err := coll.Find(context.TODO(), bson.M{"Username": username})

	var books []Book

	if err != nil {
		defer result.Close(*ctx)
	} else {
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
	}
	return books
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

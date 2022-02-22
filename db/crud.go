package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/piotrek-hub/book.io-backend/utils"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

func Register(login string, password string) (string, error) {

	ctx, client, coll := connect("users")
	defer client.Disconnect(*ctx)

	_, exisits := utils.CheckIfUserExists(bson.D{{"Login", login}, {"Password", password}}, coll)

	if !exisits {
		userKey := utils.DeriveUserKey(login, password)
		doc := bson.D{{"Login", login}, {"Password", password}, {"UserKey", userKey}}
		_, err := coll.InsertOne(context.TODO(), doc)

		if err != nil {
			return "", errors.New("Error during register")
		}
		return userKey, nil

	} else {
		return "", errors.New("User already exisits")
	}

}

func Login(login string, password string) (string, error) {
	ctx, client, coll := connect("users")
	defer client.Disconnect(*ctx)

	userKey, userExisits := utils.CheckIfUserExists(bson.D{{"Login", login}, {"Password", password}}, coll)

	if userExisits {
		return userKey, nil
	}
	return "", errors.New("User doesnt exists")
}

func AddBook(bookRequest utils.BookRequest) error {
	if bookRequest.UserKey == "" {
		return errors.New("User key not provided")
	}
	if len(bookRequest.Username) == 0 {
		return errors.New("Username not provided")
	}

	ctx, client, coll := connect("books")
	defer client.Disconnect(*ctx)

	usersCtx, usersClient, users := connect("users")
	defer usersClient.Disconnect(*usersCtx)

	_, userExists := utils.CheckIfUserExists(bson.D{{"UserKey", bookRequest.UserKey}}, users)
	bookExists := utils.CheckIfBookExists(bookRequest, coll)


	if bookExists {
		return errors.New("Book already exists")
	}
	if !userExists || !bookExists {
		return errors.New("User doesnt userExists or book already!")
	}

	doc := utils.InitBookDoc(bookRequest, bookRequest.UserKey, bookRequest.Username)

	_, err := coll.InsertOne(context.TODO(), doc)
	if err != nil {
		return errors.New("Error while adding book to db")
	}

	return nil
}

func SetBookStatus(bookRequest utils.BookRequest) (error) {
	if bookRequest.UserKey == "" {
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
		bson.M{"Title": bookRequest.Title, "UserKey": bookRequest.UserKey},
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
	if bookRequest.UserKey == "" {
		return errors.New("User key not provided")
	}

	ctx, client, coll := connect("books")
	defer client.Disconnect(*ctx)

	if !utils.CheckIfBookExists(bookRequest, coll) {
		return errors.New("Book doesnt exists")
	}

	result, err := coll.DeleteOne(*ctx, bson.M{"Title": bookRequest.Title, "UserKey": bookRequest.UserKey})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
	return nil
}

func GetBooks(username string) ([]Book, error) {
	var books []Book

	ctx, client, coll := connect("books")
	defer client.Disconnect(*ctx)

	result, err := coll.Find(context.TODO(), bson.M{"Username": username})
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

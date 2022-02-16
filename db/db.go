package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// env variables

const uri = "mongodb+srv://mesi:qwer1234@cluster0.ffdei.mongodb.net/Cluster0?retryWrites=true&w=majority"

func connect(dbName string) (*context.Context, *mongo.Client, *mongo.Collection) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	coll := client.Database("bookio").Collection(dbName)
	return &ctx, client, coll
}

func Register(login string, password string) string {

	ctx, client, coll := connect("users")
	defer client.Disconnect(*ctx)

	_, userExisits := checkIfUserExisits(login, password, coll)

	if !userExisits {
		userKey := deriveUserKey(login, password)
		doc := bson.D{{"login", login}, {"password", password}, {"userKey", userKey}}
		_, err := coll.InsertOne(context.TODO(), doc)

		if err != nil {
			return "Error during register"
		}

	} else {
		return "User already exisits"
	}

	return "User registered successfully"
}

func Login(login string, password string) string {
	ctx, client, coll := connect("users")
	defer client.Disconnect(*ctx)

	user, userExisits := checkIfUserExisits(login, password, coll)

	if userExisits {
		return user.UserKey
	} else {
		return "User doesnt exists"
	}
}

func AddBook(bookRequest BookRequest) string {
	ctx, client, coll := connect("books")
	defer client.Disconnect(*ctx)

	usersCtx, usersClient, users := connect("users")
	defer usersClient.Disconnect(*usersCtx)

	// Check if user exists
	fmt.Println(bookRequest)
	userKey := bookRequest.UserKey
	_, exisits := checkIfUserExisitsByKey(userKey, users)

	if !exisits {
		return "Caller doesnt exists"
	}

	doc := initBookDoc(bookRequest, userKey)

	_, err := coll.InsertOne(context.TODO(), doc)

	if err != nil {
		return "Error while adding book to db"
	}

	return "Book added successfully"
}

func SetBookStatus(bookRequest BookRequest) string {
	ctx, client, coll := connect("books")
	defer client.Disconnect(*ctx)

	fmt.Println(bookRequest)

	bookExists := checkIfBookExists(bookRequest, coll)

	if !bookExists {
		return "Book doesnt exisits"
	}

	_, err := coll.UpdateOne(
		*ctx,
		bson.M{"Title": bookRequest.Title, "userKey": bookRequest.UserKey},
		bson.D{
			{"$set", bson.D{{"Status", bookRequest.Status}}},
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	return "Status changed successfully"
}



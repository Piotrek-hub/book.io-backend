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

func connect() (*context.Context, *mongo.Client, *mongo.Collection) {
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
	// databases, err := client.ListDatabaseNames(ctx, bson.M{})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	coll := client.Database("bookio").Collection("users")
	return &ctx, client, coll
}

func Register(login string, password string) string {

	ctx, client, coll := connect()
	defer client.Disconnect(*ctx)

	var res bson.M
	err := coll.FindOne(context.TODO(), bson.D{{"login", login}}).Decode(&res)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			var userKey string = deriveUserKey(login, password)
			doc := bson.D{{"login", login}, {"password", password}, {"userKey", userKey}}
			_, err := coll.InsertOne(context.TODO(), doc)
			if err != nil {
				return "Error during register"
			}
			return "Registered successfully"
		}
	}
	return "User already exists"

}

func Login(login string, password string) string {
	ctx, client, coll := connect()
	defer client.Disconnect(*ctx)

	var res bson.M
	err := coll.FindOne(context.TODO(), bson.D{{"login", login}, {"password", password}}).Decode(&res)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("No User Found")
			return ""
		}
	}
	var user User
	bsonBytes, _ := bson.Marshal(res)
	bson.Unmarshal(bsonBytes, &user)

	return user.UserKey
}

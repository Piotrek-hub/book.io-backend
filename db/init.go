package db

import (
	"context"
	"github.com/piotrek-hub/book.io-backend/utils"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)



func connect(dbName string) (*context.Context, *mongo.Client, *mongo.Collection) {
	uri := utils.GetDatabaseUri()

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



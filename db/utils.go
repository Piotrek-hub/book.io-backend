package db

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func deriveUserKey(login string, password string) string {
	hash := sha512.New()
	hash.Write([]byte(login + password + time.Now().String()))
	md := hash.Sum(nil)
	mdStr := hex.EncodeToString(md)
	fmt.Println(mdStr)
	return mdStr
}

// bson.D{{"login", login}, {"password", password}}
func checkIfUserExisits(filter bson.D, coll *mongo.Collection) (User, bool) {
	var result bson.M
	var user User
	err := coll.FindOne(context.TODO(), filter).Decode(&result)

	if err == mongo.ErrNoDocuments {
		return user, false
	}
	if err != nil {
		fmt.Println("Error calling FindOne():", err)
		return user, false
	}

	bsonBytes, _ := bson.Marshal(result)
	bson.Unmarshal(bsonBytes, &user)

	return user, true
}


func initBookDoc(bookRequest BookRequest, userKey string) bson.D {
	return bson.D{
		{"Title", bookRequest.Title},
		{"Author", bookRequest.Author},
		{"Pages", bookRequest.Pages},
		{"DateCompleted", bookRequest.DateCompleted},
		{"Status", bookRequest.Status},
		{"UserKey", userKey},
	}
}

func checkIfBookExists(bookRequest BookRequest, coll *mongo.Collection) bool {
	var result bson.M
	err := coll.FindOne(context.TODO(), bson.D{{"Title", bookRequest.Title}, {"UserKey", bookRequest.UserKey}}).Decode(&result)

	if err == mongo.ErrNoDocuments {
		return false
	}
	if err != nil {
		fmt.Println("Error calling FindOne():", err)
		return false
	}

	return true
}

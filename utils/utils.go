package utils

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func DeriveUserKey(login string, password string) string {
	hash := sha512.New()
	hash.Write([]byte(login + password + time.Now().String()))
	md := hash.Sum(nil)
	mdStr := hex.EncodeToString(md)
	fmt.Println(mdStr)
	return mdStr
}

func CheckIfUserExists(filter bson.D, coll *mongo.Collection) (string, bool) {
	var result bson.M
	err := coll.FindOne(context.TODO(), filter).Decode(&result)

	userKey := fmt.Sprintf("%v", result["UserKey"])
	fmt.Println(userKey)

	if err == mongo.ErrNoDocuments {
		return "", false
	}
	if err != nil {
		fmt.Println("Error calling FindOne():", err)
		return "", false
	}

	return userKey, true
}

func InitBookDoc(bookRequest BookRequest, userKey string, username string) bson.D {
	return bson.D{
		{"Title", bookRequest.Title},
		{"Author", bookRequest.Author},
		{"Pages", bookRequest.Pages},
		{"DateCompleted", bookRequest.DateCompleted},
		{"Status", bookRequest.Status},
		{"UserKey", userKey},
		{"Username", username},
	}
}

func CheckIfBookExists(bookRequest BookRequest, coll *mongo.Collection) bool {
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

func GetDatabaseUri() (string) {
	yfile, err := ioutil.ReadFile("./config/config.yaml")

	if err != nil {
		log.Fatal(err)
	}

	data := make(map[string]string)

	err = yaml.Unmarshal(yfile, &data)

	if err != nil {
		log.Fatal(err)
	}

	var values []string;
	for _, value := range data {
		values = append(values, value)
	}
	return "mongodb+srv://"+values[0]+":"+values[1]+"@cluster0.ffdei.mongodb.net/Cluster0?retryWrites=true&w=majority"
}
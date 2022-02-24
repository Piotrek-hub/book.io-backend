package utils

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var mySigningKey = []byte("mysecretphrase")

func GenerateToken(login string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user"] = login
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func CheckIfUserExists(filter bson.D, coll *mongo.Collection) (string, bool) {
	var result bson.M
	err := coll.FindOne(context.TODO(), filter).Decode(&result)

	token := fmt.Sprintf("%v", result["token"])

	if err == mongo.ErrNoDocuments {
		log.SetFlags(log.Ldate | log.Lshortfile)
		log.Println("User doesnt exists")
		return "", false
	}

	if err != nil {
		return "", false
	}

	return token, true
}

func InitBookDoc(bookRequest BookRequest, token string, username string) bson.D {
	return bson.D{
		{"title", bookRequest.Title},
		{"author", bookRequest.Author},
		{"pages", bookRequest.Pages},
		{"dateCompleted", bookRequest.DateCompleted},
		{"status", bookRequest.Status},
		{"username", username},
	}
}

func CheckIfBookExists(bookRequest BookRequest, coll *mongo.Collection) bool {
	var result bson.M
	err := coll.FindOne(context.TODO(), bson.D{{"title", bookRequest.Title}, {"username", bookRequest.Username}}).Decode(&result)

	if err == mongo.ErrNoDocuments || err != nil{
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

func LogRequest[T any](message string, request T) {
	log.Println(message, request)
}
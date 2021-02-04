package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectDB is helper function to connect MongoDB Atlas
func ConnectDB() *mongo.Collection {
	clientOptions := options.Client().ApplyURI("mongodb+srv://retr00exe:system32@backend-cluster.vy2xd.mongodb.net/crypt-v2?retryWrites=true&w=majority")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB")
	collection := client.Database("my-db").Collection("people")
	return collection
}

// ErrorResponse define the error model
type ErrorResponse struct {
	StatusCode   int    `json:"status"`
	ErrorMessage string `json:"message"`
}

// GetError define helper to debug code
func GetError(err error, w http.ResponseWriter) {
	log.Fatal(err.Error())
	var response = ErrorResponse{
		StatusCode:   http.StatusInternalServerError,
		ErrorMessage: err.Error(),
	}
	message, _ := json.Marshal(response)
	w.WriteHeader(response.StatusCode)
	w.Write(message)
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Person struct define the model for http response
type Person struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
}

var client *mongo.Client

// CreatePersonEndpoint define endpoint append person to database
func CreatePersonEndpoint(response http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	response.Header().Add("content-type", "application/json")
	fmt.Println("OK 1")
	var person Person
	_ = json.NewDecoder(request.Body).Decode(&person)
	fmt.Println("OK 2")
	collection := client.Database("my-db").Collection("people")
	fmt.Println("OK 3")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := collection.InsertOne(ctx, person)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(response).Encode(result)
}

func main() {
	fmt.Println("Connecting to MongoDB Atlas ...")
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://retr00exe:system32@backend-cluster.vy2xd.mongodb.net/crypt-v2?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)
	router := mux.NewRouter()
	router.HandleFunc("/person", CreatePersonEndpoint).Methods("POST")
	http.ListenAndServe("127.0.0.1:8080", router)
}

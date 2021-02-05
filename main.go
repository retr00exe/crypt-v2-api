package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"quickstart/helper"
	"quickstart/models"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	fmt.Println("Connecting to MongoDB Atlas ...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://retr00exe:system32@backend-cluster.vy2xd.mongodb.net/crypt-v2?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("PING success! Connected to MongoDB Atlas")
	r := mux.NewRouter()
	r.HandleFunc("/api/person", getPeople).Methods("GET")
	r.HandleFunc("/api/person/{id}", getPerson).Methods("GET")
	r.HandleFunc("/api/person", createPerson).Methods("POST")
	http.ListenAndServe("127.0.0.1:8080", r)
}

// getPeople define endpoint to get list of all users in database
func getPeople(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	var people []models.Person
	collection := helper.ConnectDB()
	result, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		helper.GetError(err, w)
		return
	}
	for result.Next(context.TODO()) {
		var person models.Person
		err := result.Decode(&person)
		if err != nil {
			log.Fatal(err)
		}
		people = append(people, person)
	}
	if err := result.Err(); err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(people)
	fmt.Println("GET /api/person 200 OK!")
}

func getPerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	var person models.Person
	var params = mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	filter := bson.M{"_id": id}
	collection := helper.ConnectDB()
	err := collection.FindOne(context.TODO(), filter).Decode(&person)
	if err != nil {
		helper.GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(person)
	fmt.Println("GET /api/person/{id} 200 OK!")
}

// createPerson define endpoint to append person into the database
func createPerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	var person models.Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	collection := helper.ConnectDB()
	result, err := collection.InsertOne(context.TODO(), person)
	if err != nil {
		helper.GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(result)
	fmt.Println("POST /api/person 200 OK!")
}

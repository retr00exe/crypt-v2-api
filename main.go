package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"quickstart/helper"
	"quickstart/models"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var personCollection *mongo.Collection = helper.ConnectDB("my-db", "people")
var hashCollection *mongo.Collection = helper.ConnectDB("crypt-v2", "hash")
var host string = "127.0.0.1:8080"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connecting to MongoDB Atlas ...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("ATLAS_URI")))
	defer client.Disconnect(ctx)
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
	r.HandleFunc("/api/person/{id}", updatePerson).Methods("PUT")
	r.HandleFunc("/api/person/{id}", deletePerson).Methods("DELETE")
	r.HandleFunc("/api/hash", getHashes).Methods("GET")
	r.HandleFunc("/api/hash", createHash).Methods("POST")
	log.Fatal(http.ListenAndServe(host, r))
}

/*
	@route 	GET /api/person
	@desc		Get all users
	@access	Public
*/
func getPeople(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	var people []models.Person
	result, err := personCollection.Find(context.TODO(), bson.M{})
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
	fmt.Println("GET	/api/person 200 OK!")
}

/*
	@route 	GET /api/person/:id
	@desc		Get specific user
	@access	Public
*/
func getPerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	var person models.Person
	var params = mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	filter := bson.M{"_id": id}
	err := personCollection.FindOne(context.TODO(), filter).Decode(&person)
	if err != nil {
		helper.GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(person)
	fmt.Println("GET	/api/person/:id 200 OK!")
}

/*
	@route 	POST /api/person
	@desc		Create a user
	@access	Public
*/
func createPerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	var person models.Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	result, err := personCollection.InsertOne(context.TODO(), person)
	if err != nil {
		helper.GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(result)
	fmt.Println("POST	/api/person 200 OK!")
}

/*
	@route 	PUT /api/person/:id
	@desc		Update a user
	@access	Public
*/
func updatePerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	var params = mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var person models.Person
	filter := bson.M{"_id": id}
	_ = json.NewDecoder(r.Body).Decode(&person)
	update := bson.M{
		"$set": bson.M{
			"firstname": person.Firstname,
			"lastname":  person.Lastname,
		},
	}
	err := personCollection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&person)
	if err != nil {
		helper.GetError(err, w)
		return
	}
	person.ID = id
	json.NewEncoder(w).Encode(person)
	fmt.Println("PUT	/api/person/:id 200 OK!")
}

/*
	@route 	DELETE /api/person/:id
  @desc		Delete a users
	@access	Public
*/
func deletePerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	var params = mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	filter := bson.M{"_id": id}
	deleteResult, err := personCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		helper.GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(deleteResult)
	fmt.Println("DELETE	/api/person/:id 200 OK!")
}

/*
	@route 	GET /api/hash
	@desc		Get all users
	@access	Public
*/
func getHashes(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	var hashes []models.Hash
	result, err := hashCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		helper.GetError(err, w)
		return
	}
	for result.Next(context.TODO()) {
		var hash models.Hash
		err := result.Decode(&hash)
		if err != nil {
			log.Fatal(err)
		}
		hashes = append(hashes, hash)
	}
	if err := result.Err(); err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(hashes)
	fmt.Println("GET	/api/hash 200 OK!")
}

/*
	@route 	GET /api/hash/:id
	@desc		Get specific user
	@access	Public
*/
func getHash(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	var hash models.Hash
	var params = mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	filter := bson.M{"_id": id}
	err := hashCollection.FindOne(context.TODO(), filter).Decode(&hash)
	if err != nil {
		helper.GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(hash)
	fmt.Println("GET	/api/hash/:id 200 OK!")
}

/*
	@route 	POST /api/hash
	@desc		Create a hash description
	@access	Public
*/
func createHash(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	var hash models.Hash
	_ = json.NewDecoder(r.Body).Decode(&hash)
	result, err := personCollection.InsertOne(context.TODO(), hash)
	if err != nil {
		helper.GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(result)
	fmt.Println("POST	/api/hash 200 OK!")
}

/*
	@route	PUT /api/hash/:id
	@desc		Update a user
	@access	Public
*/
func updateHash(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	var params = mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var hash models.Hash
	filter := bson.M{"_id": id}
	_ = json.NewDecoder(r.Body).Decode(&hash)
	update := bson.M{
		"$set": bson.M{
			"name":        hash.Name,
			"digest_size": hash.DigestSize,
			"block_size":  hash.BlockSize,
			"rounds":      hash.Rounds,
			"date":        hash.DatePublished,
			"designer":    hash.Designer,
		},
	}
	err := hashCollection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&hash)
	if err != nil {
		helper.GetError(err, w)
		return
	}
	hash.ID = id
	json.NewEncoder(w).Encode(hash)
	fmt.Println("PUT	/api/hash/:id 200 OK!")
}

/*
	@route	DELETE /api/hash/:id
  @desc		Delete a users
	@access	Public
*/
func deleteHash(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	var params = mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	filter := bson.M{"_id": id}
	deleteResult, err := hashCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		helper.GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(deleteResult)
	fmt.Println("DELETE	/api/hash/:id 200 OK!")
}

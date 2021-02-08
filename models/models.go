package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Person struct define the model for http response
type Person struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
}

// Hash is
type Hash struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name          string             `json:"name,omitempty" bson:"name,omitempty"`
	DigestSize    int16              `json:"digest_size,omitempty" bson:"digest_size,omitempty"`
	BlockSize     int16              `json:"block_size,omitempty" bson:"block_size,omitempty"`
	Rounds        int8               `json:"rounds,omitempty" bson:"rounds,omitempty"`
	DatePublished primitive.DateTime `json:"date,omitempty" bson:"date,omitempty"`
	Designer      string             `json:"designer,omitempty" bson:"designer,omitempty"`
}

// ResponseHash is
type ResponseHash struct {
	Status  int16  `json:"status"`
	Message string `json:"message"`
	Data    Hash   `json:"data" bson:"data"`
}

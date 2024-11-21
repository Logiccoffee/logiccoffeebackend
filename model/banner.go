package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Banner struct {
	ID   primitive.ObjectID `json:"id" bson:"_id"`
	Name string             `json:"name" bson:"name"`
	Photo string            `json:"photo" bson:"photo"`
}

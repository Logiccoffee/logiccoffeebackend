package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description,omitempty" json:"description"`
	Price       float64            `bson:"price" json:"price"`
	Available   bool               `bson:"available" json:"available"`
	Photo       string             `bson:"photo" json:"photo"` // URL or filename
}

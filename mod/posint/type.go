package posint

import "go.mongodb.org/mongo-driver/bson/primitive"

type Item struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Destination     string             `bson:"Destination" json:"Destination"`
	ProhibitedItems string             `bson:"Prohibited Items" json:"Prohibited Items"`
}

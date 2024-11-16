package kimseok

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Datasets struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Question string             `json:"question" bson:"question"`
	Answer   string             `json:"answer" bson:"answer"`
	Origin   string             `json:"origin" bson:"origin"`
}

type Session struct {
	ID          string    `bson:"_id,omitempty"`
	PhoneNumber string    `bson:"phonenumber"`
	CreatedAt   time.Time `bson:"createdAt"`
}

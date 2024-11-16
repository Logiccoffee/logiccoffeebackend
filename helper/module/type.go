package module

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Module struct {
	Name         string   `json:"name,omitempty" bson:"name,omitempty"`
	Keyword      []string `json:"keyword,omitempty" bson:"keyword,omitempty"`
	Phonenumbers []string `json:"phonenumbers,omitempty" bson:"phonenumbers,omitempty"`
	Group        bool     `json:"group,omitempty" bson:"group,omitempty"`
	Personal     bool     `json:"personal,omitempty" bson:"personal,omitempty"`
}
type Typo struct {
	From string `json:"from,omitempty" bson:"from,omitempty"`
	To   string `json:"to,omitempty" bson:"to,omitempty"`
}

type Session struct {
	ID          string     `bson:"_id,omitempty"`
	PhoneNumber string     `bson:"phonenumber"`
	Menulist    []MenuList `bson:"list"`
	CreatedAt   time.Time  `bson:"createdAt"`
}

type Menu struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"` // Field ID untuk MongoDB
	Keyword string             `bson:"keyword,omitempty" json:"keyword,omitempty"`
	Header  string             `bson:"header,omitempty" json:"header,omitempty"`
	List    []MenuList         `bson:"list,omitempty" json:"list,omitempty"`
	Footer  string             `bson:"footer,omitempty" json:"footer,omitempty"`
}

type MenuList struct {
	No      int    `bson:"no,omitempty" json:"no,omitempty"`
	Keyword string `bson:"keyword,omitempty" json:"keyword,omitempty"`
	Konten  string `bson:"konten,omitempty" json:"konten,omitempty"`
}

type Datasets struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Question string             `json:"question" bson:"question"`
	Answer   string             `json:"answer" bson:"answer"`
	Origin   string             `json:"origin" bson:"origin"`
}

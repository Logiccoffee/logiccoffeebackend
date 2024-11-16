package otorisasiwebstatis

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Stp struct {
	PhoneNumber  string `bson:"phonenumber,omitempty" json:"phonenumber,omitempty"`
	PasswordHash string `bson:"password,omitempty" json:"password,omitempty"`
}
type Config struct {
	PhoneNumber          string `json:"phonenumber,omitempty" bson:"phonenumber,omitempty"`
	LeaflyURL            string `json:"leaflyurl,omitempty" bson:"leaflyurl,omitempty"`
	LeaflyURLKTP         string `json:"leaflyurlktp,omitempty" bson:"leaflyurlktp,omitempty"`
	LeaflySecret         string `json:"leaflysecret,omitempty" bson:"leaflysecret,omitempty"`
	DomyikadoPresensiURL string `json:"domyikadopresensiurl,omitempty" bson:"domyikadopresensiurl,omitempty"`
	DomyikadoUserURL     string `json:"domyikadouserurl,omitempty" bson:"domyikadouserurl,omitempty"`
	DomyikadoSecret      string `json:"domyikadosecret,omitempty" bson:"domyikadosecret,omitempty"`
}

type Userdomyikado struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name                 string             `bson:"name,omitempty" json:"name,omitempty"`
	PhoneNumber          string             `bson:"phonenumber,omitempty" json:"phonenumber,omitempty"`
	Email                string             `bson:"email,omitempty" json:"email,omitempty"`
	GithubUsername       string             `bson:"githubusername,omitempty" json:"githubusername,omitempty"`
	GitlabUsername       string             `bson:"gitlabusername,omitempty" json:"gitlabusername,omitempty"`
	GitHostUsername      string             `bson:"githostusername,omitempty" json:"githostusername,omitempty"`
	Poin                 float64            `bson:"poin,omitempty" json:"poin,omitempty"`
	GoogleProfilePicture string             `bson:"googleprofilepicture,omitempty" json:"picture,omitempty"`
	PasswordHash         string             `bson:"passwordhash,omitempty" json:"passwordhash,omitempty"`
	PasswordExpiry       time.Time          `bson:"passwordexpiry,omitempty" json:"passwordexpiry,omitempty"`
}

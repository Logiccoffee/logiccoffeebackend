package config

import (
	"log"
	"os"

	"github.com/gocroot/helper/at"
	"github.com/gocroot/helper/atdb"
	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var PrivateKey string = os.Getenv("PRKEY")

var AESKey string = os.Getenv("AESKEY")

var IPPort, Net = at.GetAddress()

var PhoneNumber string = os.Getenv("PHONENUMBER")

var Profile itmodel.Profile

func SetEnv() {
	if ErrorMongoconn != nil {
		log.Println(ErrorMongoconn.Error())
	}
	Profile, err := atdb.GetOneDoc[itmodel.Profile](Mongoconn, "profile", primitive.M{"phonenumber": PhoneNumber})
	if err != nil {
		log.Println(err)
	}
	PublicKeyWhatsAuth = Profile.PublicKey
	WAAPIToken = Profile.Token
}

var GHAccessToken string = os.Getenv("GH_ACCESS_TOKEN")

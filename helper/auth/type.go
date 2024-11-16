package auth

import "time"

type GoogleCredential struct {
	Token        string    `bson:"token"`
	RefreshToken string    `bson:"refresh_token"`
	TokenURI     string    `bson:"token_uri"`
	ClientID     string    `bson:"client_id"`
	ClientSecret string    `bson:"client_secret"`
	Scopes       []string  `bson:"scopes"`
	Expiry       time.Time `bson:"expiry"`
}
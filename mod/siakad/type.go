package siakad

type Config struct {
	PhoneNumber               string `json:"phonenumber,omitempty" bson:"phonenumber,omitempty"`
	SiakadLoginURL            string `json:"siakadloginurl,omitempty" bson:"siakadloginurl,omitempty"`
	BapURL                    string `json:"bapurl,omitempty" bson:"bapurl,omitempty"`
	ApproveBimbinganURL       string `json:"approvebimbinganurl,omitempty" bson:"approvebimbinganurl,omitempty"`
	ApproveBimbinganByPoinURL string `json:"approvebimbinganbypoinurl,omitempty" bson:"approvebimbinganbypoinurl,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

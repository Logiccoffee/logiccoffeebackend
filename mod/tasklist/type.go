package tasklist

import "go.mongodb.org/mongo-driver/bson/primitive"

type TaskList struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	MeetID      primitive.ObjectID `json:"meetid,omitempty" bson:"meetid,omitempty"`
	LaporanID   primitive.ObjectID `json:"laporanid,omitempty" bson:"laporanid,omitempty"`
	UserID      primitive.ObjectID `json:"userid,omitempty" bson:"userid,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	PhoneNumber string             `json:"phonenumber,omitempty" bson:"phonenumber,omitempty"`
	Email       string             `json:"email,omitempty" bson:"email,omitempty"`
	Task        string             `json:"task,omitempty" bson:"task,omitempty"`
	ProjectID   primitive.ObjectID `json:"projectid,omitempty" bson:"projectid,omitempty"`
	ProjectName string             `json:"projectname,omitempty" bson:"projectname,omitempty"`
	IsDone      bool               `json:"isdone,omitempty" bson:"isdone,omitempty"`
}

type Config struct {
	PhoneNumber          string `json:"phonenumber,omitempty" bson:"phonenumber,omitempty"`
	LeaflyURL            string `json:"leaflyurl,omitempty" bson:"leaflyurl,omitempty"`
	LeaflySecret         string `json:"leaflysecret,omitempty" bson:"leaflysecret,omitempty"`
	DomyikadoPresensiURL string `json:"domyikadopresensiurl,omitempty" bson:"domyikadopresensiurl,omitempty"`
	DomyikadoTaskListURL string `json:"domyikadotasklisturl,omitempty" bson:"domyikadotasklisturl,omitempty"`
	DomyikadoSecret      string `json:"domyikadosecret,omitempty" bson:"domyikadosecret,omitempty"`
}

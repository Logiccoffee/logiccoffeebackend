package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Peserta struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"` // Field ID untuk MongoDB
	KodeProvinsi  string             `bson:"kodeProvinsi,omitempty" json:"kodeProvinsi,omitempty"`
	KodeKabupaten string             `bson:"kodeKabupaten,omitempty" json:"kodeKabupaten,omitempty"`
	KodeKecamatan string             `bson:"kodeKecamatan,omitempty" json:"kodeKecamatan,omitempty"`
	KodeDesa      string             `bson:"kodeDesa,omitempty" json:"kodeDesa,omitempty"`
	Provinsi      string             `bson:"provinsi,omitempty" json:"provinsi,omitempty"`
	Kab           string             `bson:"kab,omitempty" json:"kab,omitempty"`
	Kec           string             `bson:"kec,omitempty" json:"kec,omitempty"`
	Desa          string             `bson:"desa,omitempty" json:"desa,omitempty"`
	Fullname      string             `bson:"fullname,omitempty" json:"fullname,omitempty"`
	Username      string             `bson:"username,omitempty" json:"username,omitempty"`
	PhoneNumber   string             `bson:"phoneNumber,omitempty" json:"phoneNumber,omitempty"`
	Position      string             `bson:"position,omitempty" json:"position,omitempty"`
	Approved      string             `bson:"Approved,omitempty" json:"Approved,omitempty"`
	IsOnWhatsApp  bool               `bson:"isOnWhatsApp,omitempty" json:"isOnWhatsApp,omitempty"`
	Komentar      string             `json:"komentar,omitempty" bson:"komentar,omitempty"`
	Rating        int                `json:"rating,omitempty" bson:"rating,omitempty"`
}

type Testi struct {
	Isi    string `json:"isi,omitempty" bson:"isi,omitempty"`
	Nama   string `json:"nama,omitempty" bson:"nama,omitempty"`
	Daerah string `json:"daerah,omitempty" bson:"daerah,omitempty"`
}

type Depan struct {
	List []Testi `json:"list,omitempty" bson:"list,omitempty"`
}

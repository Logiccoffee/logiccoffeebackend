package model

import (
	"time"

	"github.com/gocroot/helper/gcallapi"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Laporan struct {
	ID          primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty" query:"id" url:"_id,omitempty" reqHeader:"_id"`
	MeetID      primitive.ObjectID   `json:"meetid,omitempty" bson:"meetid,omitempty"`
	MeetEvent   gcallapi.SimpleEvent `json:"meetevent,omitempty" bson:"meetevent,omitempty"`
	Project     Project              `json:"project,omitempty" bson:"project,omitempty"`
	User        Userdomyikado        `json:"user,omitempty" bson:"user,omitempty"`
	Petugas     string               `json:"petugas,omitempty" bson:"petugas,omitempty"`
	NoPetugas   string               `json:"nopetugas,omitempty" bson:"nopetugas,omitempty"`
	Kode        string               `json:"kode,omitempty" bson:"kode,omitempty"`
	Team        string               `json:"team,omitempty" bson:"team,omitempty"`
	Scope       string               `json:"scope,omitempty" bson:"scope,omitempty"`
	Nama        string               `json:"nama,omitempty" bson:"nama,omitempty"`
	Phone       string               `json:"phone,omitempty" bson:"phone,omitempty"`
	Masalah     string               `json:"masalah,omitempty" bson:"masalah,omitempty"`
	Solusi      string               `json:"solusi,omitempty" bson:"solusi,omitempty"`
	Komentar    string               `json:"komentar,omitempty" bson:"komentar,omitempty"`
	Terlayani   bool                 `json:"terlayani,omitempty" bson:"terlayani,omitempty"`
	Rating      float64              `json:"rating,omitempty" bson:"rating,omitempty"`
	RateLayanan int                  `json:"ratelayanan,omitempty" bson:"ratelayanan,omitempty"`
}

type PushReport struct {
	ProjectName string   `bson:"projectname" json:"projectname"`
	Project     Project  `bson:"project" json:"project"`
	User        MenuItem `bson:"menu,omitempty" json:"menu,omitempty"`
	Username    string   `bson:"username" json:"username"`
	Email       string   `bson:"email,omitempty" json:"email,omitempty"`
	Repo        string   `bson:"repo" json:"repo"`
	Ref         string   `bson:"ref" json:"ref"`
	Message     string   `bson:"message" json:"message"`
	Modified    string   `bson:"modified,omitempty" json:"modified,omitempty"`
	RemoteAddr  string   `bson:"remoteaddr,omitempty" json:"remoteaddr,omitempty"`
}

type Project struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Secret               string             `bson:"secret,omitempty" json:"secret,omitempty"`
	Name                 string             `bson:"name" json:"name"`
	Title                string             `bson:"title" json:"title"`
	Description          string             `bson:"description" json:"description"`
	KalimatPromosi       string             `bson:"kalimatpromosi" json:"kalimatpromosi"`
	Owner                Userdomyikado      `bson:"owner,omitempty" json:"owner,omitempty"`
	Editor               Userdomyikado      `bson:"editor,omitempty" json:"editor,omitempty"`
	Manager              Userdomyikado      `bson:"manager,omitempty" json:"manager,omitempty"`
	IsApproved           bool               `bson:"isapproved,omitempty" json:"isapproved,omitempty"`
	WAGroupID            string             `bson:"wagroupid,omitempty" json:"wagroupid,omitempty"`
	RepoOrg              string             `bson:"repoorg,omitempty" json:"repoorg,omitempty"`
	RepoLogName          string             `bson:"repologname,omitempty" json:"repologname,omitempty"`
	CoverBuku            string             `bson:"coverbuku,omitempty" json:"coverbuku,omitempty"`
	DraftBuku            string             `bson:"draftbuku,omitempty" json:"draftbuku,omitempty"`
	DraftPDFBuku         string             `bson:"draftpdfbuku,omitempty" json:"draftpdfbuku,omitempty"`
	SampulPDFBuku        string             `bson:"sampulpdfbuku,omitempty" json:"sampulpdfbuku,omitempty"`
	URLKatalog           string             `bson:"urlkatalog,omitempty" json:"urlkatalog,omitempty"`
	PATHKatalog          string             `bson:"pathkatalog,omitempty" json:"pathkatalog,omitempty"`
	SPK                  string             `bson:"spk,omitempty" json:"spk,omitempty"`
	SPI                  string             `bson:"spi,omitempty" json:"spi,omitempty"`
	ISBN                 string             `bson:"isbn,omitempty" json:"isbn,omitempty"`
	Terbit               string             `bson:"terbit,omitempty" json:"terbit,omitempty"`
	Ukuran               string             `bson:"ukuran,omitempty" json:"ukuran,omitempty"`
	JumlahHalaman        string             `bson:"jumlahhalaman,omitempty" json:"jumlahhalaman,omitempty"`
	Tebal                string             `bson:"tebal,omitempty" json:"tebal,omitempty"`
	LinkGramed           string             `bson:"linkgramed,omitempty" json:"linkgramed,omitempty"`
	LinkPlayBook         string             `bson:"linkplaybook,omitempty" json:"linkplaybook,omitempty"`
	LinkKubuku           string             `bson:"linkkubuku,omitempty" json:"linkkubuku,omitempty"`
	LinkMyedisi          string             `bson:"linkmyedisi,omitempty" json:"linkmyedisi,omitempty"`
	LinkDepositPerpusnas string             `bson:"linkdepositperpusnas,omitempty" json:"linkdepositperpusnas,omitempty"`
	LinkDepositPerpusda  string             `bson:"linkdepositperpusda,omitempty" json:"linkdepositperpusda,omitempty"`
	Menu                 []MenuItem         `bson:"menu,omitempty" json:"menu,omitempty"`
	Members              []Userdomyikado    `bson:"members,omitempty" json:"members,omitempty"`
	Closed               bool               `bson:"closed,omitempty" json:"closed,omitempty"`
}

type MenuItem struct {
	IDDatabase primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ID         string             `json:"id,omitempty" bson:"id,omitempty"`
	Name       string             `json:"name,omitempty" bson:"name,omitempty"`
	Price      int                `json:"price,omitempty" bson:"price,omitempty"`
	Image      string             `json:"image,omitempty" bson:"image,omitempty"`
}

type Userdomyikado struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name                 string             `bson:"name,omitempty" json:"name,omitempty"`
	PhoneNumber          string             `bson:"phonenumber,omitempty" json:"phonenumber,omitempty"`
	Email                string             `bson:"email,omitempty" json:"email,omitempty"`
	NIK                  string             `bson:"nik,omitempty" json:"nik,omitempty"`
	Pekerjaan            string             `bson:"pekerjaan,omitempty" json:"pekerjaan,omitempty"`
	AlamatRumah          string             `bson:"alamatrumah,omitempty" json:"alamatrumah,omitempty"`
	AlamatKantor         string             `bson:"alamatkantor,omitempty" json:"alamatkantor,omitempty"`
	GithubUsername       string             `bson:"githubusername,omitempty" json:"githubusername,omitempty"`
	GitlabUsername       string             `bson:"gitlabusername,omitempty" json:"gitlabusername,omitempty"`
	GitHostUsername      string             `bson:"githostusername,omitempty" json:"githostusername,omitempty"`
	Poin                 float64            `bson:"poin,omitempty" json:"poin,omitempty"`
	GoogleProfilePicture string             `bson:"googleprofilepicture,omitempty" json:"picture,omitempty"`
	ProfilePicture       string             `bson:"profpic,omitempty" json:"profpic,omitempty"`
	Bio                  string             `bson:"bio,omitempty" json:"bio,omitempty"`
	URLBio               string             `bson:"urlbio,omitempty" json:"urlbio,omitempty"`
	PATHBio              string             `bson:"pathbio,omitempty" json:"pathbio,omitempty"`
	Team                 string             `json:"team,omitempty" bson:"team,omitempty"`
	Scope                string             `json:"scope,omitempty" bson:"scope,omitempty"`
	Section              string             `json:"section,omitempty" bson:"section,omitempty"`
	Chapter              string             `json:"chapter,omitempty" bson:"chapter,omitempty"`
	LinkedDevice         string             `json:"linkeddevice,omitempty" bson:"linkeddevice,omitempty"`
	JumlahAntrian        int                `json:"jumlahantrian,omitempty" bson:"jumlahantrian,omitempty"`
	IsEditor             bool               `json:"iseditor,omitempty" bson:"iseditor,omitempty"`
	IsManager            bool               `json:"ismanager,omitempty" bson:"ismanager,omitempty"`
}

type Task struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ProjectID string             `bson:"projectid" json:"projectid"`
	Name      string             `bson:"name" json:"name"`
	PIC       Userdomyikado      `bson:"pic" json:"pic"`
	Done      bool               `bson:"done,omitempty" json:"done,omitempty"`
}

type LoginRequest struct {
	PhoneNumber string `json:"phonenumber"`
	Password    string `json:"password"`
}

type Stp struct {
	PhoneNumber  string    `bson:"phonenumber,omitempty" json:"phonenumber,omitempty"`
	PasswordHash string    `bson:"password,omitempty" json:"password,omitempty"`
	CreatedAt    time.Time `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
}

type VerifyRequest struct {
	PhoneNumber string `json:"phonenumber"`
	Password    string `json:"password"`
}

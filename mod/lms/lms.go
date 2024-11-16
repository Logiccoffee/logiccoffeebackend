package lms

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gocroot/helper/atapi"
	"github.com/gocroot/helper/atdb"
	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func ReplyRekapUsers(Profile itmodel.Profile, Pesan itmodel.IteungMessage, db *mongo.Database) (msg string) {
	//kasih pesan dulu biar nunggu
	msgstr := "Hai kak " + Pesan.Alias_name + " permintaannya sedang di proses nih. mohon tunggu sekitar 3 menit ya kak."
	dt := &itmodel.TextMessage{
		To:       Pesan.Chat_number,
		IsGroup:  Pesan.Is_group,
		Messages: msgstr,
	}
	go atapi.PostStructWithToken[itmodel.Response]("Token", Profile.Token, dt, Profile.URLAPIText)

	//lanjutkan rekap
	rkp, err := GetRekapPendaftaranUsers(db)
	if err != nil {
		msg = "Gagal mendapatkan rekap pendaftaran user:" + err.Error()
		return
	}
	msg = "Berikut ini rekapitulasi terupdate saat ini tentang pendaftaran user:\n"
	msg += "Total user: " + strconv.Itoa(int(rkp.Total))
	msg += "\nBelum Lengkap: " + strconv.Itoa(int(rkp.BelumLengkap))
	msg += "\nMenunggu Persetujuan: " + strconv.Itoa(int(rkp.MenungguPersetujuan))
	msg += "\nDisetujui: " + strconv.Itoa(int(rkp.Disetujui))
	msg += "\nDitolak: " + strconv.Itoa(int(rkp.Ditolak))
	return

}

func RefreshCookie(db *mongo.Database) (err error) {
	profile, err := atdb.GetOneDoc[LoginProfile](db, "lmscreds", bson.M{})
	if err != nil {
		return
	}
	newxs, newls, newbar, err := GetNewCookie(profile.Xsrf, profile.Lsession, db)
	if err != nil {
		return
	}
	profile.Bearer = newbar
	profile.Xsrf = newxs
	profile.Lsession = newls
	_, err = atdb.ReplaceOneDoc(db, "lmscreds", bson.M{"username": "madep"}, profile)
	if err != nil {
		return
	}
	return

}

func GetTotalUser(db *mongo.Database) (total int, err error) {
	profile, err := atdb.GetOneDoc[LoginProfile](db, "lmscreds", bson.M{})
	if err != nil {
		return
	}
	url := profile.URLUsers
	url = strings.ReplaceAll(url, "##PAGE##", "1")

	_, res, err := atapi.GetWithBearer[Root](profile.Bearer, url)
	if err != nil {
		err = errors.New("GetWithBearer:" + err.Error() + " " + url + " " + profile.Bearer)
		return
	}
	total = res.Data.Meta.Total
	return
}

func GetAllUser(db *mongo.Database) (users []User, err error) {
	profile, err := atdb.GetOneDoc[LoginProfile](db, "lmscreds", bson.M{})
	if err != nil {
		return
	}

	i := 1
	var res Root
	for {
		url := strings.ReplaceAll(profile.URLUsers, "##PAGE##", strconv.Itoa(i))
		_, res, err = atapi.GetWithBearer[Root](profile.Bearer, url)
		if err != nil {
			err = errors.New("GetWithBearer:" + err.Error() + profile.Bearer + " " + url)
			return
		}
		users = append(users, res.Data.Data...)
		if res.Data.Meta.LastItem == res.Data.Meta.Total {
			break
		}
		i++
		//users = res.Data.Data
	}
	return
}

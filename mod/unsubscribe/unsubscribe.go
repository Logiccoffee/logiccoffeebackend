package unsubscribe

import (
	"github.com/gocroot/helper/atdb"
	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/mongo"
)

func Unsubscribe(Pesan itmodel.IteungMessage, db *mongo.Database) (reply string) {
	_, err := atdb.InsertOneDoc(db, "unsubscribe", Pesan)
	if err != nil {
		return err.Error()
	}
	return "Terima kasih atas informasi dan kerjasamanya kak *" + Pesan.Alias_name + "*\nBerhasil untuk unsubscribe informasi dari kami.\nSalam tim helpdesk LMS Pamong Desa\npd.my.id"
}

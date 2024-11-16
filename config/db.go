package config

import (
	"os"

	"github.com/gocroot/helper/atdb"
)

var MongoString string = os.Getenv("MONGOSTRING")

var mongoinfo = atdb.DBInfo{
	DBString: MongoString,
	DBName:   "naskah",
}

var Mongoconn, ErrorMongoconn = atdb.MongoConnect(mongoinfo)

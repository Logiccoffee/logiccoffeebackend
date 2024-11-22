package config
import (
	"os"
	"github.com/gocroot/helper/atdb"
)
var MongoString string = os.Getenv("mongodb+srv://fathyafathazzra:Mongodbatlas12@cluster0.8xvps.mongodb.net/")
var mongoinfo = atdb.DBInfo{
	DBString: MongoString,
	DBName:   "logiccoffee",
}
var Mongoconn, ErrorMongoconn = atdb.MongoConnect(mongoinfo)
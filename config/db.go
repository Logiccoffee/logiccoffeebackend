package config

import (
	"context"
	"log"
	"time"

	"github.com/gocroot/helper/atdb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Mendefinisikan MongoString secara langsung
var MongoString string = "mongodb+srv://fathyafathazzra:Mongodbatlas12@cluster0.8xvps.mongodb.net/"

// Konfigurasi database dengan nama 'akuntan'
var mongoinfo = atdb.DBInfo{
	DBString: MongoString,
	DBName:   "logiccoffee",
}

var Mongoconn, ErrorMongoconn = atdb.MongoConnect(mongoinfo)


// Inisialisasi variabel untuk MongoDB client dan koleksi
var Client *mongo.Client
var CategoryCollection *mongo.Collection
var BannerCollection *mongo.Collection
var MenuCollection *mongo.Collection
var UsersCollection *mongo.Collection
var OrdersCollection *mongo.Collection
var ReviewCollection *mongo.Collection

// Fungsi untuk menginisialisasi koneksi ke MongoDB
func InitMongoDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(MongoString)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Gagal terhubung ke MongoDB: %v", err)
	}

	// Memastikan koneksi berhasil
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Gagal ping ke MongoDB: %v", err)
	}

	log.Println("MongoDB connected")

	// Menetapkan variabel global client dan collections
	Client = client
	CategoryCollection = client.Database("logiccoffee").Collection("category")
	BannerCollection = client.Database("logiccoffee").Collection("banner")
	MenuCollection = client.Database("logiccoffee").Collection("menu")
	UsersCollection = client.Database("logiccoffee").Collection("users")
	OrdersCollection = client.Database("logiccoffee").Collection("orders")
	ReviewCollection = client.Database("logiccoffee").Collection("review")

}
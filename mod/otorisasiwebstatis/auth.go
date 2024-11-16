package otorisasiwebstatis

import (
	"context"
	"crypto/rand"
	"math/big"
	"net/http"
	"regexp"
	"time"

	"github.com/gocroot/helper/atapi"
	"github.com/gocroot/helper/atdb"
	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func generateRandomPassword(length int) (string, error) {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    b := make([]byte, length)
    for i := range b {
        randomInt, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
        if err != nil {
            return "", err
        }
        b[i] = charset[randomInt.Int64()]
    }
    return string(b), nil
}

func hashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func AuthWhatsApp(Pesan itmodel.IteungMessage, db *mongo.Database) (reply string) {
    // Validate phone number
    re := regexp.MustCompile(`^\+62\d{9,15}$`)
    if !re.MatchString(Pesan.Phone_number) {
        return "Nomor telepon tidak sesuai format Indonesia"
    }

    // Retrieve config from database
    conf, err := atdb.GetOneDoc[Config](db, "config", bson.M{"phonenumber": "62895601060000"})
	if err != nil {
		return "Wah kak " + Pesan.Alias_name + " mohon maaf ada kesalahan dalam pengambilan config di database " + err.Error()
	}

    // Generate random password
    randomPassword, err := generateRandomPassword(12)
    if err != nil {
        return "Gagal menghasilkan password acak: " + err.Error()
    }

    // Hash the password
    hashedPassword, err := hashPassword(randomPassword)
    if err != nil {
        return "Gagal meng-hash password: " + err.Error()
    }

    // Set password expiry time to 1 minute
    passwordExpiry := time.Now().Add(1 * time.Minute)

    // Update or insert the user in the database
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    collection := db.Collection("userdomyikado")
    filter := bson.M{"phonenumber": Pesan.Phone_number}

    update := bson.M{
        "$set": Userdomyikado{
            PhoneNumber:  Pesan.Phone_number,
            Name:         Pesan.Alias_name,
            PasswordHash: hashedPassword, // Store hashed password in database
            PasswordExpiry: passwordExpiry, // Store password expiry time
        },
    }
    opts := options.Update().SetUpsert(true)
    _, err = collection.UpdateOne(ctx, filter, update, opts)
    if err != nil {
        return "Gagal menyimpan informasi pengguna: " + err.Error()
    }

    // Send the random password via WhatsApp
    recentUserAuth := Userdomyikado{
        PhoneNumber: Pesan.Phone_number,
        Name:        Pesan.Alias_name,
    }

    statuscode, httpresp, err := atapi.PostStructWithToken[itmodel.Response]("secret", conf.DomyikadoSecret, recentUserAuth, conf.DomyikadoUserURL)
    if err != nil {
        return "Akses ke endpoint domyikado gagal: " + err.Error()
    }
    if statuscode != http.StatusOK {
        return "Salah posting endpoint domyikado: " + httpresp.Response + "\ninfo\n" + httpresp.Info
    }

    return "Hai kak " + Pesan.Alias_name + ", password login telah dikirim ke nomor WhatsApp Anda. Password akan kedaluwarsa dalam 1 menit. Hashed Password: " + hashedPassword
}
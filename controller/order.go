package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gocroot/config"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/helper/watoken"
	"github.com/gocroot/model"
	"go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func formatrupiah(price float64) string {
	formatter := message.NewPrinter(language.Indonesian)
	return formatter.Sprintf("Rp %.2f", price)
}

// CreateOrder - Membuat order baru
// CreateOrder - Membuat order baru
func CreateOrder(respw http.ResponseWriter, req *http.Request) {
	// Dekode token WhatsAuth untuk mengambil informasi pengguna
	payload, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
	if err != nil {
		at.WriteJSON(respw, http.StatusForbidden, model.Response{
			Status:   "Error: Token Tidak Valid",
			Location: "Decode Token Error",
			Response: err.Error(),
		})
		return
	}

	// Ambil data JSON dari body request
	var order model.Order
	if err := json.NewDecoder(req.Body).Decode(&order); err != nil {
		at.WriteJSON(respw, http.StatusBadRequest, model.Response{
			Status:   "Error: Bad Request",
			Response: err.Error(),
		})
		return
	}

	// Ambil data user berdasarkan PhoneNumber yang ada di payload dari database
var user model.Userdomyikado
filter := bson.M{"phonenumber": payload.Id}

// Gunakan operator =, bukan :=, karena variabel 'user' sudah ada
user, err = atdb.GetOneDoc[model.Userdomyikado](config.Mongoconn, "user", filter)
if err != nil {
	at.WriteJSON(respw, http.StatusNotFound, model.Response{
		Status:   "Error: Data Pengguna Tidak Ditemukan",
		Response: err.Error(),
	})
	return
}



	// Membuat order baru dengan data dari request body atau default dari database
	newOrder := model.Order{
		OrderNumber:   order.OrderNumber,
		QueueNumber:   order.QueueNumber,
		OrderDate: time.Now(),
		UserID:        user.ID,
		UserInfo: model.UserInfo{
			Name:     order.UserInfo.Name,     // Jika ada data dari frontend, gunakan data tersebut
			Whatsapp: order.UserInfo.Whatsapp, // Jika ada data dari frontend, gunakan data tersebut
			Note:     order.UserInfo.Note,     // Jika ada data dari frontend, gunakan data tersebut
		},
		Orders:        order.Orders,
		Total:         order.Total,
		PaymentMethod: order.PaymentMethod,
		Status:        "terkirim",
		CreatedBy:     user.Name,
		CreatedByRole: user.Role,
		CreatedAt:     time.Now().Unix(),
	}
	// Simpan ke database MongoDB tanpa menggunakan InsertedID
	insertResult, err := atdb.InsertOneDoc(config.Mongoconn, "orders", newOrder)
	if err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, model.Response{
			Status:   "Error: Gagal Insert Database",
			Response: err.Error(),
		})
		return
	}
	newOrder.ID = insertResult
	// Membuat response data
	response := map[string]interface{}{
		"status":  "success",
		"message": "Order berhasil dibuat",
		"user":    user.Name, // Gunakan nama user dari database
		"data": map[string]interface{}{
			"order_number":   newOrder.OrderNumber,
			"queue_number":   newOrder.QueueNumber,
			"order_date":     newOrder.OrderDate.Format("15:04:05 02-01-2006"),
			"total":          formatrupiah(newOrder.Total),
			"payment_method": newOrder.PaymentMethod,
			"status":         newOrder.Status,
		},
	}

	// Kirim response ke client
	at.WriteJSON(respw, http.StatusOK, response)
}

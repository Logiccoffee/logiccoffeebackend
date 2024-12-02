package controller

import (
	"encoding/json"
	"net/http"
	"time"
	"strings"

	"github.com/gocroot/config"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/helper/watoken"
	"github.com/gocroot/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// GetAllOrder - Ambil Semua Data Order
func GetAllOrder(respw http.ResponseWriter, req *http.Request) {
	// Dekode token WhatsAuth untuk validasi
	_, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
	if err != nil {
		at.WriteJSON(respw, http.StatusForbidden, model.Response{
			Status:   "Error: Token Tidak Valid",
			Location: "Decode Token Error",
			Response: err.Error(),
		})
		return
	}

	// Ambil semua data order
	data, err := atdb.GetAllDoc[[]model.Order](config.Mongoconn, "orders", bson.M{})
	if err != nil {
		at.WriteJSON(respw, http.StatusNotFound, model.Response{
			Status:   "Error: Data order tidak ditemukan",
			Response: err.Error(),
		})
		return
	}

	if len(data) == 0 {
		at.WriteJSON(respw, http.StatusNotFound, model.Response{
			Status: "Error: Data order kosong",
		})
		return
	}

	var orders []map[string]interface{}
	for _, order := range data {
		orders = append(orders, map[string]interface{}{
			"id":             order.ID.Hex(),
			"order_number":   order.OrderNumber,
			"queue_number":   order.QueueNumber,
			"order_date":     order.OrderDate, // Karena sudah dalam format string
			"user_id":        order.UserID.Hex(),
			"user_info": map[string]interface{}{
				"name":     order.UserInfo.Name,
				"whatsapp": order.UserInfo.Whatsapp,
				"note":     order.UserInfo.Note,
			},
			"orders":          order.Orders,
			"total":           formatrupiah(order.Total),
			"payment_method":  order.PaymentMethod,
			"status":          order.Status,
			"created_by":      order.CreatedBy,
			"created_by_role": order.CreatedByRole,
			"created_at":      time.Unix(order.CreatedAt, 0).Format("15:04:05 02-01-2006"),
		})
	}

	at.WriteJSON(respw, http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Data order berhasil diambil",
		"data":    orders,
	})
}

// GetOrderByID - Ambil Order Berdasarkan ID
func GetOrderByID(respw http.ResponseWriter, req *http.Request) {
	// Dekode token WhatsAuth untuk validasi
	_, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
	if err != nil {
		at.WriteJSON(respw, http.StatusForbidden, model.Response{
			Status:   "Error: Token Tidak Valid",
			Location: "Decode Token Error",
			Response: err.Error(),
		})
		return
	}

	// Ambil ID dari URL
	pathParts := strings.Split(req.URL.Path, "/")
	orderID := pathParts[len(pathParts)-1]
	if orderID == "" {
		at.WriteJSON(respw, http.StatusBadRequest, model.Response{
			Status: "Error: ID Order tidak ditemukan di URL",
		})
		return
	}

	// Ubah ID string menjadi ObjectID
	objectID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		at.WriteJSON(respw, http.StatusBadRequest, model.Response{
			Status: "Error: ID Order tidak valid",
		})
		return
	}

	// Ambil data order berdasarkan ID
	order, err := atdb.GetOneDoc[model.Order](config.Mongoconn, "orders", bson.M{"_id": objectID})
	if err != nil {
		at.WriteJSON(respw, http.StatusNotFound, model.Response{
			Status:   "Error: Order tidak ditemukan",
			Response: err.Error(),
		})
		return
	}

	// Membuat response untuk data order lengkap
	response := map[string]interface{}{
		"status":  "success",
		"message": "Order ditemukan",
		"data": map[string]interface{}{
			"id":             order.ID.Hex(),
			"order_number":   order.OrderNumber,
			"queue_number":   order.QueueNumber,
			"order_date":     order.OrderDate, // Karena sudah dalam format string
			"user_id":        order.UserID.Hex(),
			"user_info": map[string]interface{}{
				"name":     order.UserInfo.Name,
				"whatsapp": order.UserInfo.Whatsapp,
				"note":     order.UserInfo.Note,
			},
			"orders":          order.Orders,
			"total":           formatrupiah(order.Total),
			"payment_method":  order.PaymentMethod,
			"status":          order.Status,
			"created_by":      order.CreatedBy,
			"created_by_role": order.CreatedByRole,
			"created_at":      time.Unix(order.CreatedAt, 0).Format("15:04:05 02-01-2006"),
		},
	}

	// Kirim response ke client
	at.WriteJSON(respw, http.StatusOK, response)
}

func UpdateOrder(respw http.ResponseWriter, req *http.Request) {
    // Ambil token dari header dan decode menggunakan public key
    payload, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
    if err != nil {
        var respn model.Response
        respn.Status = "Error: Token Tidak Valid"
        respn.Info = at.GetSecretFromHeader(req)
        respn.Location = "Decode Token Error"
        respn.Response = err.Error()
        at.WriteJSON(respw, http.StatusForbidden, respn)
        return
    }

    // Ambil ID order dari URL
    pathParts := strings.Split(req.URL.Path, "/")
    orderID := pathParts[len(pathParts)-1] // Ambil bagian terakhir dari URL
    if orderID == "" {
        var respn model.Response
        respn.Status = "Error: ID Order tidak ditemukan di URL"
        at.WriteJSON(respw, http.StatusBadRequest, respn)
        return
    }

    // Konversi ID order ke ObjectID MongoDB
    objectID, err := primitive.ObjectIDFromHex(orderID)
    if err != nil {
        var respn model.Response
        respn.Status = "Error: ID Order tidak valid"
        at.WriteJSON(respw, http.StatusBadRequest, respn)
        return
    }

    // Decode body langsung ke map
    var requestBody map[string]interface{}
    err = json.NewDecoder(req.Body).Decode(&requestBody)
    if err != nil {
        var respn model.Response
        respn.Status = "Error: Gagal membaca data JSON"
        respn.Response = err.Error()
        at.WriteJSON(respw, http.StatusBadRequest, respn)
        return
    }

    // Periksa apakah body request kosong
    if len(requestBody) == 0 {
        var respn model.Response
        respn.Status = "Error: Tidak ada data untuk diperbarui"
        at.WriteJSON(respw, http.StatusBadRequest, respn)
        return
    }

    // Menyiapkan data untuk update (langsung menggantikan field yang diberikan)
    updateData := bson.M{}
    if name, exists := requestBody["name"]; exists && name != "" {
        updateData["name"] = name
    }
    if phonenumber, exists := requestBody["phonenumber"]; exists && phonenumber != "" {
        updateData["phonenumber"] = phonenumber
    }
    if status, exists := requestBody["status"]; exists && status != "" {
        updateData["status"] = status
    }

    // Jika tidak ada perubahan data, beri respon error
    if len(updateData) == 0 {
        var respn model.Response
        respn.Status = "Error: Tidak ada perubahan yang dilakukan"
        at.WriteJSON(respw, http.StatusNotModified, respn)
        return
    }

    // Update data order
    result, err := atdb.UpdateOneDoc(config.Mongoconn, "orders", bson.M{"_id": objectID}, updateData)
    if err != nil {
        var respn model.Response
        respn.Status = "Error: Gagal mengupdate order"
        respn.Response = err.Error()
        at.WriteJSON(respw, http.StatusInternalServerError, respn)
        return
    }

    // Jika tidak ada dokumen yang dimodifikasi, beri respons error
    if result.ModifiedCount == 0 {
        var respn model.Response
        respn.Status = "Error: Tidak ada perubahan yang dilakukan"
        at.WriteJSON(respw, http.StatusNotModified, respn)
        return
    }

    // Respons sukses
    response := map[string]interface{}{
        "status":  "success",
        "message": "Order berhasil diupdate",
        "data": map[string]interface{}{
            "id":            objectID.Hex(),
            "updatedFields": updateData,
        },
		"updatedBy": payload.Alias,
    }
    at.WriteJSON(respw, http.StatusOK, response)
}

func DeleteOrder(respw http.ResponseWriter, req *http.Request) {
    // Ambil token dari header dan decode menggunakan public key
    payload, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
    if err != nil {
        var respn model.Response
        respn.Status = "Error: Token Tidak Valid"
        respn.Info = at.GetSecretFromHeader(req)
        respn.Location = "Decode Token Error"
        respn.Response = err.Error()
        at.WriteJSON(respw, http.StatusForbidden, respn)
        return
    }

    // Ambil ID order dari URL
    pathParts := strings.Split(req.URL.Path, "/")
    orderID := pathParts[len(pathParts)-1] // Ambil bagian terakhir dari URL
    if orderID == "" {
        var respn model.Response
        respn.Status = "Error: ID Order tidak ditemukan di URL"
        at.WriteJSON(respw, http.StatusBadRequest, respn)
        return
    }

    // Konversi ID order ke ObjectID MongoDB
    objectID, err := primitive.ObjectIDFromHex(orderID)
    if err != nil {
        var respn model.Response
        respn.Status = "Error: ID Order tidak valid"
        at.WriteJSON(respw, http.StatusBadRequest, respn)
        return
    }

    // Hapus data order berdasarkan ID
    filter := bson.M{"_id": objectID}
    deleteResult, err := atdb.DeleteOneDoc(config.Mongoconn, "orders", filter)
    if err != nil {
        var respn model.Response
        respn.Status = "Error: Gagal menghapus order"
        respn.Response = err.Error()
        at.WriteJSON(respw, http.StatusInternalServerError, respn)
        return
    }

    if deleteResult.DeletedCount == 0 {
        var respn model.Response
        respn.Status = "Error: Order tidak ditemukan"
        at.WriteJSON(respw, http.StatusNotFound, respn)
        return
    }

    // Berhasil menghapus order
    response := map[string]interface{}{
        "status":  "success",
        "message": "Order berhasil dihapus",
        "user":    payload.Alias,
        "data":    deleteResult,
    }
    at.WriteJSON(respw, http.StatusOK, response)
}

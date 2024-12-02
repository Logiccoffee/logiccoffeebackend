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

	// Dapatkan data user dari token untuk mencatat CreatedBy dan CreatedByRole
	var user model.Userdomyikado
	filter := bson.M{"phonenumber": payload.Id}
	user, err = atdb.GetOneDoc[model.Userdomyikado](config.Mongoconn, "user", filter)
	if err != nil {
		at.WriteJSON(respw, http.StatusNotFound, model.Response{
			Status:   "Error: Data Pengguna Tidak Ditemukan",
			Response: err.Error(),
		})
		return
	}

	// Pastikan UserInfo diisi manual di backend berdasarkan data request
	if order.UserInfo.Name == "" || order.UserInfo.Whatsapp == "" {
		at.WriteJSON(respw, http.StatusBadRequest, model.Response{
			Status:   "Error: Bad Request",
			Response: "UserInfo harus berisi Name dan Whatsapp",
		})
		return
	}

	// Validasi PaymentMethod
	allowedPaymentMethods := map[string]bool{"Cash": true}
	if !allowedPaymentMethods[order.PaymentMethod] {
		at.WriteJSON(respw, http.StatusBadRequest, model.Response{
			Status:   "Error: Invalid Payment Method",
			Response: "Metode pembayaran hanya diperbolehkan 'Cash'",
		})
		return
	}

	// Membuat order baru berdasarkan data dari frontend
	newOrder := model.Order{
		OrderNumber:   order.OrderNumber,
		QueueNumber:   order.QueueNumber,
		OrderDate:     time.Now(), // Gunakan waktu sekarang
		UserID:        user.ID,   // UserID hanya untuk referensi
		UserInfo: model.UserInfo{
			Name:     order.UserInfo.Name,
			Whatsapp: order.UserInfo.Whatsapp,
			Note:     order.UserInfo.Note, // Note opsional
		},
		Orders:        order.Orders,
		Total:         order.Total,
		PaymentMethod: order.PaymentMethod,
		Status:        "terkirim",
		CreatedBy:     user.Name, // Dari data user
		CreatedByRole: user.Role, // Dari data user
	}

	// Simpan order baru ke database
	insertResult, err := atdb.InsertOneDoc(config.Mongoconn, "orders", newOrder)
	if err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, model.Response{
			Status:   "Error: Gagal Insert Database",
			Response: err.Error(),
		})
		return
	}

	// Tambahkan ID yang dihasilkan ke newOrder
	newOrder.ID = insertResult

	// Format waktu Indonesia
	loc, _ := time.LoadLocation("Asia/Jakarta")
	formattedOrderDate := newOrder.OrderDate.In(loc).Format("15:04:05 02-01-2006")

	// Membuat response lengkap
	response := map[string]interface{}{
		"status":  "success",
		"message": "Order berhasil dibuat",
		"data":    newOrder,
		"formatted_date": formattedOrderDate, // Tanggal dalam format Indonesia
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
			"order_date":     order.OrderDate.Format("2006-01-02 15:04:05"), // Format Indonesia
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
			"order_date":     order.OrderDate.Format("2006-01-02 15:04:05"), // Format Indonesia
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
		},
	}

	// Kirim response ke client
	at.WriteJSON(respw, http.StatusOK, response)
}

func UpdateOrder(respw http.ResponseWriter, req *http.Request) {
    // Decode token untuk mendapatkan payload pengguna
    payload, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
    if err != nil {
        at.WriteJSON(respw, http.StatusForbidden, map[string]string{"error": "Token tidak valid"})
        return
    }

    // Ambil ID order dari URL
    pathParts := strings.Split(req.URL.Path, "/")
    orderID := pathParts[len(pathParts)-1]
    if orderID == "" {
        at.WriteJSON(respw, http.StatusBadRequest, map[string]string{"error": "ID Order tidak ditemukan"})
        return
    }

    // Konversi ID order ke ObjectID MongoDB
    objectID, err := primitive.ObjectIDFromHex(orderID)
    if err != nil {
        at.WriteJSON(respw, http.StatusBadRequest, map[string]string{"error": "ID Order tidak valid"})
        return
    }

    // Decode body request ke map
    var requestBody map[string]interface{}
    if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
        at.WriteJSON(respw, http.StatusBadRequest, map[string]string{"error": "Gagal membaca data JSON"})
        return
    }

    // Ambil data pesanan saat ini untuk validasi status
    currentOrder, err := atdb.GetOneDoc[model.Order](config.Mongoconn, "orders", bson.M{"_id": objectID})
    if err != nil {
        at.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"error": "Gagal mengambil data pesanan"})
        return
    }

    // Siapkan data untuk update
    updateData := bson.M{}

    // Validasi dan update UserInfo
    if userInfo, exists := requestBody["user_info"]; exists {
        if userInfoMap, ok := userInfo.(map[string]interface{}); ok {
            updateData["user_info"] = userInfoMap
        }
    }

    // Validasi dan update Orders
    if orders, exists := requestBody["orders"]; exists {
        if ordersArray, ok := orders.([]interface{}); ok {
            updateData["orders"] = ordersArray
        }
    }

    // Validasi dan update Status
    if status, exists := requestBody["status"]; exists {
        if statusStr, ok := status.(string); ok {
            if statusStr == "diproses" || statusStr == "selesai" {
                updateData["status"] = statusStr
            } else if statusStr == "dibatalkan" {
                // Cek apakah status saat ini memungkinkan pembatalan
                if currentOrder.Status != "terkirim" {
                    at.WriteJSON(respw, http.StatusBadRequest, map[string]string{
                        "error": "Pesanan tidak dapat dibatalkan karena status saat ini adalah '" + currentOrder.Status + "'",
                    })
                    return
                }
                updateData["status"] = "dibatalkan"
            }
        }
    }

    // Pastikan queueNumber tidak diizinkan untuk diubah
    if _, exists := requestBody["queue_number"]; exists {
        at.WriteJSON(respw, http.StatusForbidden, map[string]string{
            "error": "Queue number tidak dapat diubah",
        })
        return
    }

    // Pastikan ada data yang diupdate
    if len(updateData) == 0 {
        at.WriteJSON(respw, http.StatusBadRequest, map[string]string{"error": "Tidak ada perubahan yang dilakukan"})
        return
    }

    // Update data order di database
    result, err := atdb.UpdateOneDoc(config.Mongoconn, "orders", bson.M{"_id": objectID}, bson.M{"$set": updateData})
    if err != nil {
        at.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"error": "Gagal mengupdate order"})
        return
    }

    // Periksa apakah ada dokumen yang dimodifikasi
    if result.ModifiedCount == 0 {
        at.WriteJSON(respw, http.StatusNotModified, map[string]string{"error": "Tidak ada perubahan yang dilakukan"})
        return
    }

    // Respons sukses
    at.WriteJSON(respw, http.StatusOK, map[string]interface{}{
        "status":  "success",
        "message": "Order berhasil diupdate",
        "data": map[string]interface{}{
            "id":            objectID.Hex(),
            "updatedFields": updateData,
        },
        "updatedBy": payload.Alias,
    })
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

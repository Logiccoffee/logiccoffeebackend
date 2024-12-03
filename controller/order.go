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

// mengubah ke format rupiah
func formatrupiah(price float64) string {
	formatter := message.NewPrinter(language.Indonesian)
	return formatter.Sprintf("Rp %.2f", price)
}

// ngubah ke format rupiah buat price yang di dalam array orderitem
// TAPI GA DIPAKE DISIMPEN AJA SOALNYA KALI AJA BUTUH HAHAHA
// func SendFormattedOrder(respw http.ResponseWriter, orders []model.OrderItem) {
// 	formattedOrders := make([]map[string]interface{}, len(orders))
// 	for i, item := range orders {
// 		formattedOrders[i] = map[string]interface{}{
// 			"menu_name":       item.MenuName,
// 			"quantity":        item.Quantity,
// 			"price":           item.Price, 
// 			"price_formatted": formatrupiah(item.Price), 
// 		}
// 	}

// 	at.WriteJSON(respw, http.StatusOK, map[string]interface{}{
// 		"status": "success",
// 		"data":   formattedOrders,
// 	})
// }


// FormatToIndonesianTime - Mengonversi dan memformat waktu ke zona waktu Indonesia
func FormatToIndonesianTime(t time.Time) (string, error) {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return "", err
	}
	return t.In(loc).Format("02-01-2006 15:04:05"), nil
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

	// Ambil waktu dalam format Indonesia
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, model.Response{
			Status:   "Error: Gagal memuat zona waktu Indonesia",
			Response: err.Error(),
		})
		return
	}
	currentTimeInID := time.Now().In(location)

	// Format waktu Indonesia untuk respons
	orderDateInID, err := FormatToIndonesianTime(currentTimeInID)
	if err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, model.Response{
			Status:   "Error: Gagal memformat waktu",
			Response: err.Error(),
		})
		return
	}

	// Membuat order baru berdasarkan data dari frontend
	newOrder := model.Order{
		OrderNumber:   order.OrderNumber,
		QueueNumber:   order.QueueNumber,
		OrderDate:     currentTimeInID, // Gunakan waktu Indonesia
		UserID:        user.ID,         // UserID hanya untuk referensi
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

	// Membuat response lengkap
	response := map[string]interface{}{
		"status":         "success",
		"message":        "Order berhasil dibuat",
		"formatted_date": orderDateInID, // Tanggal dalam format Indonesia
		"data":           newOrder,
		"total":          formatRupiah(newOrder.Total),
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

	 // Deklarasi variabel orders
	 var orders []map[string]interface{}

	 for _, order := range data {
		 orderDateInID, err := FormatToIndonesianTime(order.OrderDate)
		 if err != nil {
			 at.WriteJSON(respw, http.StatusInternalServerError, model.Response{
				 Status:   "Error: Gagal memformat waktu",
				 Response: err.Error(),
			 })
			 return
		 }

		orders = append(orders, map[string]interface{}{
			"id":             order.ID.Hex(),
			"order_number":   order.OrderNumber,
			"queue_number":   order.QueueNumber,
			"order_date":     orderDateInID, // Menggunakan waktu Indonesia
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

	// Kirim respons dengan data yang sudah diformat
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

	// Format tanggal dan waktu menjadi format Indonesia
	orderDateInID, err := FormatToIndonesianTime(order.OrderDate)
	if err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, model.Response{
			Status:   "Error: Gagal memformat waktu",
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
			"order_date":     orderDateInID, // Menggunakan waktu Indonesia
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

	// Ambil user dari database menggunakan phonenumber dari payload
	var user model.Userdomyikado
	filter := bson.M{"phonenumber": payload.Id}
	user, err = atdb.GetOneDoc[model.Userdomyikado](config.Mongoconn, "user", filter)
	if err != nil {
		at.WriteJSON(respw, http.StatusNotFound, map[string]string{
			"error": "Data pengguna tidak ditemukan",
		})
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

	// Ambil data pesanan saat ini
	currentOrder, err := atdb.GetOneDoc[model.Order](config.Mongoconn, "orders", bson.M{"_id": objectID})
	if err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"error": "Gagal mengambil data pesanan"})
		return
	}

	// Validasi dan update status
	if status, exists := requestBody["status"]; exists {
		if statusStr, ok := status.(string); ok {
			switch statusStr {
			case "diproses", "selesai":
				currentOrder.Status = statusStr
			case "dibatalkan":
				if currentOrder.Status != "terkirim" {
					at.WriteJSON(respw, http.StatusBadRequest, map[string]string{
						"error": "Pesanan tidak dapat dibatalkan karena status saat ini adalah '" + currentOrder.Status + "'",
					})
					return
				}
				currentOrder.Status = "dibatalkan"
			default:
				at.WriteJSON(respw, http.StatusBadRequest, map[string]string{"error": "Status tidak valid"})
				return
			}
		}
	} else {
		at.WriteJSON(respw, http.StatusBadRequest, map[string]string{"error": "Status harus diisi"})
		return
	}

	// Tambahkan informasi siapa yang mengupdate
	currentOrder.UpdatedBy = user.Name       // Menggunakan nama dari struct Userdomyikado
	currentOrder.UpdatedByRole = user.Role   // Menggunakan role dari struct Userdomyikado

	// Gunakan waktu Indonesia untuk UpdatedAt
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"error": "Gagal memuat zona waktu Indonesia"})
		return
	}
	currentOrder.UpdatedAt = time.Now().In(location)

	// Update data order di database
	_, err = atdb.ReplaceOneDoc(config.Mongoconn, "orders", bson.M{"_id": objectID}, currentOrder)
	if err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"error": "Gagal mengupdate order"})
		return
	}

	// Respons sukses
	at.WriteJSON(respw, http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Order berhasil diupdate",
		"data": map[string]interface{}{
			"id":              objectID.Hex(),
			"updated_at":      currentOrder.UpdatedAt.Format("02 Januari 2006, 15:04 WIB"),
			"status":          currentOrder.Status,
			"updated_by":      currentOrder.UpdatedBy,
			"updated_by_role": currentOrder.UpdatedByRole,
			"orders":          currentOrder.Orders,
		},
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

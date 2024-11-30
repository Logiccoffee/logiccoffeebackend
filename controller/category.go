package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gocroot/config"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/helper/watoken"
	"github.com/gocroot/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateCategory(respw http.ResponseWriter, req *http.Request) {
	payload, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))

	if err != nil {
		payload, err = watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))

		if err != nil {
			var respn model.Response
			respn.Status = "Error: Token Tidak Valid"
			respn.Info = at.GetSecretFromHeader(req)
			respn.Location = "Decode Token Error"
			respn.Response = err.Error()
			at.WriteJSON(respw, http.StatusForbidden, respn)
			return
		}
	}

	var category model.Category
	if err := json.NewDecoder(req.Body).Decode(&category); err != nil {
		var respn model.Response
		respn.Status = "Error: Bad Request"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	newCategory := model.Category{
		Name: category.Name,
		Image: category.Image,
	}
	_, err = atdb.InsertOneDoc(config.Mongoconn, "category", newCategory)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Gagal Insert Database"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusNotModified, respn)
		return
	}
	response := map[string]interface{}{
		"status":  "success",
		"message": "Kategori berhasil ditambahkan",
		"name":    payload.Alias,
		"data":    newCategory,
	}
	at.WriteJSON(respw, http.StatusOK, response)
}

func GetAllCategory(respw http.ResponseWriter, req *http.Request) {
	// Ambil semua data kategori dari koleksi
	data, err := atdb.GetAllDoc[[]model.Category](config.Mongoconn, "category", primitive.M{})
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Data kategori tidak ditemukan"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	// Cek apakah data kosong
	if len(data) == 0 {
		var respn model.Response
		respn.Status = "Error: Data kategori kosong"
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	// Format hasil sebagai slice of map dengan ID dalam bentuk string
	var categories []map[string]interface{}
	for _, category := range data {
		categories = append(categories, map[string]interface{}{
			"id":    category.ID.Hex(), // Konversi ObjectID ke string
			"name":  category.Name,
			"image": category.Image,
		})
	}

	// Kirim data kategori dalam format JSON
	at.WriteJSON(respw, http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Data kategori berhasil diambil",
		"data":    categories,
	})
}


func GetCategoryByID(respw http.ResponseWriter, req *http.Request) {
	categoryID := req.URL.Query().Get("id")
	if categoryID == "" {
		var respn model.Response
		respn.Status = "Error: ID Category tidak ditemukan"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: ID Category tidak valid"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	 // Query ke database untuk mengambil kategori
	 filter := bson.M{"_id": objectID}
	 category, err := atdb.GetOneDoc[model.Category](config.Mongoconn, "category", filter)
	 if err != nil {
		 var respn model.Response
		 respn.Status = "Error: Category tidak ditemukan"
		 respn.Response = err.Error()
		 at.WriteJSON(respw, http.StatusNotFound, respn)
		 return
	 }

	response := map[string]interface{}{
		"status":  "success",
		"message": "Category ditemukan",
		"data":    category,
	}
	at.WriteJSON(respw, http.StatusOK, response)
}


func UpdateCategory(respw http.ResponseWriter, req *http.Request) {
    // Decode token untuk validasi
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

    // Ambil ID kategori dari query parameter
    categoryID := req.URL.Query().Get("id")
    if categoryID == "" {
        var respn model.Response
        respn.Status = "Error: ID Category tidak ditemukan"
        at.WriteJSON(respw, http.StatusBadRequest, respn)
        return
    }

    // Validasi format ID
    objectID, err := primitive.ObjectIDFromHex(categoryID)
    if err != nil {
        var respn model.Response
        respn.Status = "Error: ID Category tidak valid"
        at.WriteJSON(respw, http.StatusBadRequest, respn)
        return
    }

    // Periksa apakah kategori ada di database
    filter := bson.M{"_id": objectID}
    _, err = atdb.GetOneDoc[model.Category](config.Mongoconn, "category", filter)
    if err != nil {
        var respn model.Response
        respn.Status = "Error: Category tidak ditemukan"
        at.WriteJSON(respw, http.StatusNotFound, respn)
        return
    }

    // Decode data permintaan
    var requestBody struct {
        Name  string `json:"name"`
        Image string `json:"image"`
    }
    err = json.NewDecoder(req.Body).Decode(&requestBody)
    if err != nil {
        var respn model.Response
        respn.Status = "Error: Gagal membaca data JSON"
        at.WriteJSON(respw, http.StatusBadRequest, respn)
        return
    }

    // Siapkan data pembaruan
    updateData := bson.M{}
    if requestBody.Name != "" {
        updateData["name"] = requestBody.Name
    }
    if requestBody.Image != "" {
        updateData["image"] = requestBody.Image
    }

    // Tambahkan logika untuk memeriksa jika tidak ada data yang diupdate
    if len(updateData) == 0 {
        var respn model.Response
        respn.Status = "Error: Tidak ada data untuk diperbarui"
        at.WriteJSON(respw, http.StatusBadRequest, respn)
        return
    }

    // Lakukan update ke database
    update := bson.M{"$set": updateData}
    _, err = atdb.UpdateOneDoc(config.Mongoconn, "category", filter, update)
    if err != nil {
        var respn model.Response
        respn.Status = "Error: Gagal mengupdate category"
        respn.Response = err.Error()
        at.WriteJSON(respw, http.StatusNotModified, respn)
        return
    }

    // Respons sukses
    response := map[string]interface{}{
        "status":  "success",
        "message": "category berhasil diupdate",
        "data":    updateData,
        "name":    payload.Alias,
    }
    at.WriteJSON(respw, http.StatusOK, response)
}


func DeleteCategory(respw http.ResponseWriter, req *http.Request) {
	// Ambil token dari header
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

	// Ambil ID toko dari query parameter
	categoryID := req.URL.Query().Get("id")
	if categoryID == "" {
		var respn model.Response
		respn.Status = "Error: ID Toko tidak ditemukan"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	// Konversi categoryID dari string ke ObjectID MongoDB
	objectID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: ID Toko tidak valid"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	// Hapus data Category berdasarkan ID
	filter := bson.M{"_id": objectID}
	deleteResult, err := atdb.DeleteOneDoc(config.Mongoconn, "category", filter)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Gagal menghapus Category"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusInternalServerError, respn)
		return
	}

	if deleteResult.DeletedCount == 0 {
		var respn model.Response
		respn.Status = "Error: Category tidak ditemukan"
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Category berhasil dihapus",
		"user":    payload.Alias,
		"data":    deleteResult,
	}
	at.WriteJSON(respw, http.StatusOK, response)
}
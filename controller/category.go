package controller

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gocroot/config"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/helper/watoken"
	"github.com/gocroot/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateCategory(respw http.ResponseWriter, req *http.Request) {
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

    // Decode body untuk mendapatkan data kategori
    var category model.Category
    if err := json.NewDecoder(req.Body).Decode(&category); err != nil {
        var respn model.Response
        respn.Status = "Error: Bad Request"
        respn.Response = err.Error()
        at.WriteJSON(respw, http.StatusBadRequest, respn)
        return
    }

    // Siapkan kategori baru tanpa ID
    newCategory := model.Category{
        Name:  category.Name,
        Image: category.Image,
    }

    // Masukkan kategori ke dalam database
    insertResult, err := atdb.InsertOneDoc(config.Mongoconn, "category", newCategory)
    if err != nil {
        var respn model.Response
        respn.Status = "Error: Gagal Insert Database"
        respn.Response = err.Error()
        at.WriteJSON(respw, http.StatusNotModified, respn)
        return
    }

    // Gunakan hasil insertResult langsung (karena bertipe primitive.ObjectID)
    newCategory.ID = insertResult

    // Siapkan respons
    response := map[string]interface{}{
        "status":  "success",
        "message": "Kategori berhasil ditambahkan",
        "name":    payload.Alias,
        "data": map[string]interface{}{
            "id":    newCategory.ID.Hex(), // Convert ObjectID to string
            "name":  newCategory.Name,
            "image": newCategory.Image,
        },
    }
    at.WriteJSON(respw, http.StatusOK, response)
}

func GetAllCategory(respw http.ResponseWriter, req *http.Request) {
	// Ambil semua data kategori dari koleksi
	data, err := atdb.GetAllDoc[[]model.Category](config.Mongoconn, "category", bson.M{})
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
    // Ambil ID kategori dari URL menggunakan Split
    pathParts := strings.Split(req.URL.Path, "/")
    categoryID := pathParts[len(pathParts)-1] // Ambil bagian terakhir dari URL
    if categoryID == "" {
        var respn model.Response
        respn.Status = "Error: ID Category tidak ditemukan di URL"
        at.WriteJSON(respw, http.StatusBadRequest, respn)
        return
    }

    // Konversi ID kategori ke ObjectID MongoDB
    objectID, err := primitive.ObjectIDFromHex(categoryID)
    if err != nil {
        var respn model.Response
        respn.Status = "Error: ID Category tidak valid"
        at.WriteJSON(respw, http.StatusBadRequest, respn)
        return
    }

    // Query ke database untuk mengambil kategori berdasarkan ObjectID
    filter := bson.M{"_id": objectID}
    category, err := atdb.GetOneDoc[model.Category](config.Mongoconn, "category", filter)
    if err != nil {
        var respn model.Response
        respn.Status = "Error: Category tidak ditemukan"
        respn.Response = err.Error()  // Sertakan error untuk debugging
        at.WriteJSON(respw, http.StatusNotFound, respn)
        return
    }

    // Format response jika kategori ditemukan
    response := map[string]interface{}{
        "status":  "success",
        "message": "Category ditemukan",
        "data":    category,
    }
    at.WriteJSON(respw, http.StatusOK, response)
}

func UpdateCategory(respw http.ResponseWriter, req *http.Request) {
    // Ambil token dari header untuk validasi
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

    // Ambil ID kategori dari URL menggunakan Split
    pathParts := strings.Split(req.URL.Path, "/")
    categoryID := pathParts[len(pathParts)-1] // Ambil bagian terakhir dari URL
    if categoryID == "" {
        var respn model.Response
        respn.Status = "Error: ID Category tidak ditemukan di URL"
        at.WriteJSON(respw, http.StatusBadRequest, respn)
        return
    }

    // Konversi ID kategori ke ObjectID MongoDB
    objectID, err := primitive.ObjectIDFromHex(categoryID)
    if err != nil {
        var respn model.Response
        respn.Status = "Error: ID Category tidak valid"
        at.WriteJSON(respw, http.StatusBadRequest, respn)
        return
    }

    // Decode data kategori yang ingin diupdate
    var category model.Category
    err = json.NewDecoder(req.Body).Decode(&category)
    if err != nil {
        var respn model.Response
        respn.Status = "Error: Gagal membaca data JSON"
        respn.Response = err.Error()
        at.WriteJSON(respw, http.StatusBadRequest, respn)
        return
    }

    // Periksa apakah body request kosong
    if category.Name == "" && category.Image == "" {
        var respn model.Response
        respn.Status = "Error: Tidak ada data untuk diperbarui"
        at.WriteJSON(respw, http.StatusBadRequest, respn)
        return
    }

    // Update langsung menggunakan body request
    update := bson.M{"$set": category}
    filter := bson.M{"_id": objectID}
    result, err := atdb.UpdateOneDoc(config.Mongoconn, "category", filter, update)
    if err != nil {
        var respn model.Response
        respn.Status = "Error: Gagal mengupdate category"
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
        "message": "Category berhasil diupdate",
        "data": map[string]interface{}{
            "id":    objectID.Hex(),
            "name":  category.Name,
            "image": category.Image,
        },
        "updatedBy": payload.Alias,
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

	// Ambil ID kategori dari URL
	pathParts := strings.Split(req.URL.Path, "/")
	categoryID := pathParts[len(pathParts)-1] // Ambil bagian terakhir dari URL
	if categoryID == "" {
		var respn model.Response
		respn.Status = "Error: ID Category tidak ditemukan di URL"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	// Konversi ID kategori ke ObjectID MongoDB
	objectID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: ID Category tidak valid"
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

	// Berhasil menghapus kategori
	response := map[string]interface{}{
		"status":  "success",
		"message": "Category berhasil dihapus",
		"user":    payload.Alias,
		"data":    deleteResult,
	}
	at.WriteJSON(respw, http.StatusOK, response)
}
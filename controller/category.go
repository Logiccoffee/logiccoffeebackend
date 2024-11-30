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
    // Ambil ID kategori dari query params
    categoryID := req.URL.Query().Get("id")
    if categoryID == "" {
        var respn model.Response
        respn.Status = "Error: ID Category tidak ditemukan"
        at.WriteJSON(respw, http.StatusBadRequest, respn)
        return
    }

    // Konversi ID ke ObjectID
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
        respn.Response = err.Error()  // Sertakan error untuk debugging
        at.WriteJSON(respw, http.StatusNotFound, respn)
        return
    }

    // Format response
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

    // Decode data kategori langsung ke model.Category
    var category model.Category
    err = json.NewDecoder(req.Body).Decode(&category)
    if err != nil {
        var respn model.Response
        respn.Status = "Error: Gagal membaca data JSON"
        respn.Response = err.Error()
        at.WriteJSON(respw, http.StatusBadRequest, respn)
        return
    }

    // Validasi format ID kategori
    if !category.ID.IsZero() { // Pastikan ID sudah valid
        _, err := primitive.ObjectIDFromHex(category.ID.Hex())
        if err != nil {
            var respn model.Response
            respn.Status = "Error: ID Category tidak valid"
            at.WriteJSON(respw, http.StatusBadRequest, respn)
            return
        }
    } else {
        var respn model.Response
        respn.Status = "Error: ID Category kosong"
        at.WriteJSON(respw, http.StatusBadRequest, respn)
        return
    }

    // Periksa apakah kategori ada di database
    filter := bson.M{"_id": category.ID}
    existingCategory, err := atdb.GetOneDoc[model.Category](config.Mongoconn, "category", filter)
    if err != nil {
        var respn model.Response
        respn.Status = "Error: Category tidak ditemukan"
        respn.Response = err.Error()
        at.WriteJSON(respw, http.StatusNotFound, respn)
        return
    }

    // Preserve unmodifiable fields
    category.ID = existingCategory.ID

    // Update hanya field yang boleh dimodifikasi
    updateData := bson.M{}
    if category.Name != "" && category.Name != existingCategory.Name {
        updateData["name"] = category.Name
    }
    if category.Image != "" && category.Image != existingCategory.Image {
        updateData["image"] = category.Image
    }

    // Periksa jika tidak ada perubahan
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
        "message": "Category berhasil diupdate",
        "data":    category,
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

	// Ambil ID Category dari query parameter
	categoryID := req.URL.Query().Get("id")
	if categoryID == "" {
		var respn model.Response
		respn.Status = "Error: ID Category tidak ditemukan"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	// Konversi categoryID dari string ke ObjectID MongoDB
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

	response := map[string]interface{}{
		"status":  "success",
		"message": "Category berhasil dihapus",
		"user":    payload.Alias,
		"data":    deleteResult,
	}
	at.WriteJSON(respw, http.StatusOK, response)
}
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
	// Verifikasi Token Pengguna
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

	var category model.Category
	if err := json.NewDecoder(req.Body).Decode(&category); err != nil {
		var respn model.Response
		respn.Status = "Error: Bad Request"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	// Simpan Category ke Database
	_, err = atdb.InsertOneDoc(config.Mongoconn, "category", category)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Gagal Insert Database"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusNotModified, respn)
		return
	}

	// Response sukses
	response := map[string]interface{}{
		"status":  "success",
		"message": "Category berhasil ditambahkan",
		"name":    payload.Alias,
		"data":    category,
	}
	at.WriteJSON(respw, http.StatusOK, response)
}

func GetAllCategories(respw http.ResponseWriter, req *http.Request) {
	// Verifikasi Token Pengguna
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

	// Ambil Data Category dari Database
	data, err := atdb.GetAllDoc[[]model.Category](config.Mongoconn, "category", primitive.M{})
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Data category tidak ditemukan"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	// Cek jika Data Kosong
	if len(data) == 0 {
		var respn model.Response
		respn.Status = "Error: Data category kosong"
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	// Transformasi Data Category
	var categories []map[string]interface{}
	for _, category := range data {
		categories = append(categories, map[string]interface{}{
			"id":   category.ID,
			"name": category.Name,
		})
	}

	// Response Berhasil dengan Data Category
	response := map[string]interface{}{
		"status":  "success",
		"message": "Data category berhasil diambil",
		"user":    payload.Alias,
		"data":    categories,
	}
	at.WriteJSON(respw, http.StatusOK, response)
}

func GetCategoryByID(respw http.ResponseWriter, req *http.Request) {
	// Ambil ID dari query string
	categoryID := req.URL.Query().Get("id")
	if categoryID == "" {
		var respn model.Response
		respn.Status = "Error: ID Category tidak ditemukan"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	// Convert ID ke ObjectID
	objectID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: ID Category tidak valid"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	// Ambil Data Category dari Database
	var category model.Category
	filter := bson.M{"_id": objectID}
	_, err = atdb.GetOneDoc[model.Category](config.Mongoconn, "category", filter)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Category tidak ditemukan"
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	// Response dengan Data Category
	response := map[string]interface{}{
		"status":  "success",
		"message": "Category ditemukan",
		"data":    category,
	}
	at.WriteJSON(respw, http.StatusOK, response)
}

func UpdateCategory(respw http.ResponseWriter, req *http.Request) {
	// Verifikasi Token Pengguna
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

	// Ambil ID dari query string
	categoryID := req.URL.Query().Get("id")
	if categoryID == "" {
		var respn model.Response
		respn.Status = "Error: ID Category tidak ditemukan"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	// Convert ID ke ObjectID
	objectID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: ID Category tidak valid"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	// Ambil Data Category dari Database
	filter := bson.M{"_id": objectID}
	_, err = atdb.GetOneDoc[model.Category](config.Mongoconn, "category", filter)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Category tidak ditemukan"
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	// Dekode Request Body untuk Update
	var requestBody struct {
		Name string `json:"name"`
	}
	err = json.NewDecoder(req.Body).Decode(&requestBody)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Gagal membaca data JSON"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	// Update Data Category
	updateData := bson.M{
		"$set": bson.M{"name": requestBody.Name},
	}
	_, err = atdb.UpdateOneDoc(config.Mongoconn, "category", filter, updateData)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Gagal mengupdate category"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusNotModified, respn)
		return
	}

	// Response dengan Data Category yang Diperbarui
	response := map[string]interface{}{
		"status":  "success",
		"message": "Category berhasil diupdate",
		"data":    requestBody.Name,
		"name":    payload.Alias,
	}
	at.WriteJSON(respw, http.StatusOK, response)
}

func DeleteCategory(respw http.ResponseWriter, req *http.Request) {
	// Verifikasi Token Pengguna
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

	// Ambil ID dari query string
	categoryID := req.URL.Query().Get("id")
	if categoryID == "" {
		var respn model.Response
		respn.Status = "Error: ID Category tidak ditemukan"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	// Convert ID ke ObjectID
	objectID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: ID Category tidak valid"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	// Hapus Data Category dari Database
	filter := bson.M{"_id": objectID}
	deleteResult, err := atdb.DeleteOneDoc(config.Mongoconn, "category", filter)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Gagal menghapus category"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusNotModified, respn)
		return
	}

	if deleteResult.DeletedCount == 0 {
		var respn model.Response
		respn.Status = "Error: Category tidak ditemukan"
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	// Response sukses
	response := map[string]interface{}{
		"status":  "success",
		"message": "Category berhasil dihapus",
		"name":    payload.Alias,
	}
	at.WriteJSON(respw, http.StatusOK, response)
}
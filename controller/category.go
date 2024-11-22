package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gocroot/config"
	"github.com/gocroot/model"
	"github.com/gocroot/helper/at"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Fungsi untuk menambahkan kategori baru
func CreateCategory(w http.ResponseWriter, r *http.Request) {
	var newCategory model.Category
	if err := json.NewDecoder(r.Body).Decode(&newCategory); err != nil {
		var response model.Response
		response.Status = "Error: Bad Request"
		response.Response = err.Error()
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	// Set waktu pembuatan kategori
	newCategory.CreatedAt = time.Now()

	// Insert kategori ke MongoDB
	_, err := config.CategoryCollection.InsertOne(context.Background(), newCategory)
	if err != nil {
		var response model.Response
		response.Status = "Error: Gagal Insert Database"
		response.Response = err.Error()
		at.WriteJSON(w, http.StatusInternalServerError, response)
		return
	}

	// Kirim respon sukses
	response := map[string]interface{}{
		"status":  "success",
		"message": "Kategori berhasil dibuat",
		"data":    newCategory,
	}
	at.WriteJSON(w, http.StatusCreated, response)
}

// Fungsi untuk mendapatkan daftar kategori
func GetCategories(w http.ResponseWriter, r *http.Request) {
	var categories []model.Category

	// Ambil data dari MongoDB
	cursor, err := config.CategoryCollection.Find(context.Background(), bson.M{})
	if err != nil {
		var response model.Response
		response.Status = "Error: Gagal mengambil data kategori"
		response.Response = err.Error()
		at.WriteJSON(w, http.StatusInternalServerError, response)
		return
	}
	defer cursor.Close(context.Background())

	// Decode hasil pencarian kategori
	for cursor.Next(context.Background()) {
		var category model.Category
		if err := cursor.Decode(&category); err != nil {
			var response model.Response
			response.Status = "Error: Gagal mendekode kategori"
			at.WriteJSON(w, http.StatusInternalServerError, response)
			return
		}
		categories = append(categories, category)
	}

	// Kirim data kategori sebagai respon
	at.WriteJSON(w, http.StatusOK, categories)
}

// Fungsi untuk mendapatkan detail kategori berdasarkan ID
func GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		var response model.Response
		response.Status = "Error: ID Kategori tidak ditemukan"
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		var response model.Response
		response.Status = "Error: ID Kategori tidak valid"
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	var category model.Category
	err = config.CategoryCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&category)
	if err != nil {
		var response model.Response
		response.Status = "Error: Kategori tidak ditemukan"
		at.WriteJSON(w, http.StatusNotFound, response)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Kategori ditemukan",
		"data":    category,
	}
	at.WriteJSON(w, http.StatusOK, response)
}

// Fungsi untuk mengupdate kategori berdasarkan ID
func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		var response model.Response
		response.Status = "Error: ID Kategori tidak ditemukan"
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		var response model.Response
		response.Status = "Error: ID Kategori tidak valid"
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	var updatedCategory model.Category
	if err := json.NewDecoder(r.Body).Decode(&updatedCategory); err != nil {
		var response model.Response
		response.Status = "Error: Gagal membaca data JSON"
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	updateData := bson.M{
		"name":        updatedCategory.Name,
		"updatedAt":   time.Now(),
	}

	_, err = config.CategoryCollection.UpdateOne(context.Background(), bson.M{"_id": objectID}, bson.M{"$set": updateData})
	if err != nil {
		var response model.Response
		response.Status = "Error: Gagal mengupdate kategori"
		response.Response = err.Error()
		at.WriteJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Kategori berhasil diupdate",
		"data":    updateData,
	}
	at.WriteJSON(w, http.StatusOK, response)
}


// Fungsi untuk menghapus kategori berdasarkan ID
func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		var response model.Response
		response.Status = "Error: ID Kategori tidak ditemukan"
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		var response model.Response
		response.Status = "Error: ID Kategori tidak valid"
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	deleteResult, err := config.CategoryCollection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		var response model.Response
		response.Status = "Error: Gagal menghapus kategori"
		response.Response = err.Error()
		at.WriteJSON(w, http.StatusInternalServerError, response)
		return
	}

	if deleteResult.DeletedCount == 0 {
		var response model.Response
		response.Status = "Error: Kategori tidak ditemukan"
		at.WriteJSON(w, http.StatusNotFound, response)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Kategori berhasil dihapus",
	}
	at.WriteJSON(w, http.StatusOK, response)
}
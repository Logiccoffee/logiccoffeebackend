package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gocroot/config"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/model"
	"go.mongodb.org/mongo-driver/bson"
)

// Fungsi untuk menambahkan kategori baru
func CreateCategory(w http.ResponseWriter, r *http.Request) {
	var newCategory model.Category
	if err := json.NewDecoder(r.Body).Decode(&newCategory); err != nil {
		at.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "Bad request: " + err.Error(),
		})
		return
	}

	// Insert kategori ke MongoDB
	_, err := config.CategoryCollection.InsertOne(context.Background(), bson.M{
		"id":        newCategory.ID,
		"name":      newCategory.Name,
		"createdAt": time.Now(),
	})
	if err != nil {
		at.WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"message": "Gagal insert database: " + err.Error(),
		})
		return
	}

	at.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"status":  "success",
		"message": "Kategori berhasil dibuat",
		"data":    newCategory,
	})
}

// Fungsi untuk mendapatkan daftar kategori
func GetCategories(w http.ResponseWriter, r *http.Request) {
	var categories []model.Category

	// Ambil data dari MongoDB
	cursor, err := config.CategoryCollection.Find(context.Background(), bson.M{})
	if err != nil {
		at.WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"message": "Gagal mengambil data kategori: " + err.Error(),
		})
		return
	}
	defer cursor.Close(context.Background())

	// Decode hasil pencarian kategori
	for cursor.Next(context.Background()) {
		var category bson.M
		if err := cursor.Decode(&category); err != nil {
			at.WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
				"status":  "error",
				"message": "Gagal mendekode kategori: " + err.Error(),
			})
			return
		}

		// Pastikan tipe data sesuai dengan ekspektasi
		id, _ := category["id"].(int32) // Gunakan tipe sesuai database
		name, _ := category["name"].(string)

		categories = append(categories, model.Category{
			ID:   int(id),
			Name: name,
		})
	}

	if err := cursor.Err(); err != nil {
		at.WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"message": "Kesalahan pada cursor: " + err.Error(),
		})
		return
	}

	at.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "success",
		"data":    categories,
	})
}

// Fungsi untuk mendapatkan detail kategori berdasarkan ID
func GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		at.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "ID Kategori tidak ditemukan",
		})
		return
	}

	var category model.Category
	err := config.CategoryCollection.FindOne(context.Background(), bson.M{"id": id}).Decode(&category)
	if err != nil {
		at.WriteJSON(w, http.StatusNotFound, map[string]interface{}{
			"status":  "error",
			"message": "Kategori tidak ditemukan",
		})
		return
	}

	at.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Kategori ditemukan",
		"data":    category,
	})
}

// Fungsi untuk mengupdate kategori berdasarkan ID
func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		at.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "ID Kategori tidak ditemukan",
		})
		return
	}

	var updatedCategory model.Category
	if err := json.NewDecoder(r.Body).Decode(&updatedCategory); err != nil {
		at.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "Gagal membaca data JSON: " + err.Error(),
		})
		return
	}

	updateData := bson.M{
		"name":      updatedCategory.Name,
		"updatedAt": time.Now(),
	}

	_, err := config.CategoryCollection.UpdateOne(context.Background(), bson.M{"id": id}, bson.M{"$set": updateData})
	if err != nil {
		at.WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"message": "Gagal mengupdate kategori: " + err.Error(),
		})
		return
	}

	at.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Kategori berhasil diupdate",
	})
}

// Fungsi untuk menghapus kategori berdasarkan ID
func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		at.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "ID Kategori tidak ditemukan",
		})
		return
	}

	deleteResult, err := config.CategoryCollection.DeleteOne(context.Background(), bson.M{"id": id})
	if err != nil {
		at.WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"message": "Gagal menghapus kategori: " + err.Error(),
		})
		return
	}

	if deleteResult.DeletedCount == 0 {
		at.WriteJSON(w, http.StatusNotFound, map[string]interface{}{
			"status":  "error",
			"message": "Kategori tidak ditemukan",
		})
		return
	}

	at.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Kategori berhasil dihapus",
	})
}
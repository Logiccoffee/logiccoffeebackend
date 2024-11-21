package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/gocroot/model"
)

var (
	categories []model.Category
	mutex      sync.Mutex // Untuk mengamankan akses goroutine
	nextID     = 1        // ID auto-increment
)

// Get all categories
func GetCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// Get category by ID
func GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for _, category := range categories {
		if category.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(category)
			return
		}
	}

	http.Error(w, "Category not found", http.StatusNotFound)
}

// Add a new category
func PostCategory(w http.ResponseWriter, r *http.Request) {
	var newCategory model.Category

	if err := json.NewDecoder(r.Body).Decode(&newCategory); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if newCategory.Name == "" {
		http.Error(w, "Category name is required", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	newCategory.ID = nextID
	nextID++
	categories = append(categories, newCategory)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCategory)
}

// Delete a category
func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	for i, category := range categories {
		if category.ID == id {
			categories = append(categories[:i], categories[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Category not found", http.StatusNotFound)
}

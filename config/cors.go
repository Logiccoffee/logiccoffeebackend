package config

import (
	"net/http"
)

// Daftar origins yang diizinkan
var Origins = []string{
	"https://www.bukupedia.co.id",
	"https://naskah.bukupedia.co.id",
	"https://bukupedia.co.id",
	"https://logiccoffee.id.biz.id", // Menambahkan domain frontend Anda
	// "http://127.0.0.1:5500",  //menambahkan localhost
}

// Fungsi untuk memeriksa apakah origin diizinkan
func isAllowedOrigin(origin string) bool {
	for _, o := range Origins {
		if o == origin {
			return true
		}
	}
	return false
}

// Fungsi untuk mengatur header CORS
func SetAccessControlHeaders(w http.ResponseWriter, r *http.Request) bool {
	origin := r.Header.Get("Origin")

	// Jika origin tidak diizinkan, abaikan CORS dan tidak lanjutkan
	if !isAllowedOrigin(origin) {
		return false
	}

	// Set CORS headers untuk permintaan preflight (OPTIONS)
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Login")
		w.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE,PUT")
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent) // Preflight request berhasil, tanpa body
		return true
	}

	// Set CORS headers untuk permintaan utama (GET, POST, PUT, DELETE, dll)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", origin)
	return false
}

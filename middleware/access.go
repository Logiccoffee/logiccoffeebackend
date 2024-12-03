package main

import (
	"net/http"
)

func checkRoleMiddleware(allowedRoles []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Ambil role dari header atau token
		role := r.Header.Get("Role") // Misalnya, role ada di header (atau bisa dari token yang didecode)

		// Cek apakah role ada dalam daftar yang diizinkan
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				// Jika role valid, lanjutkan ke handler berikutnya
				return
			}
		}

		// Jika role tidak valid, kirim respon error
		http.Error(w, "Forbidden", http.StatusForbidden)
	}
}

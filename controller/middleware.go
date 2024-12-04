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
)

func RoleMiddleware(allowedRoles []string) func(http.HandlerFunc) http.HandlerFunc {
    return func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            payload, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(r))
            if err != nil {
                http.Error(w, "Forbidden: Invalid Token", http.StatusForbidden)
                return
            }

            docuser, err := atdb.GetOneDoc[model.Userdomyikado](config.Mongoconn, "user", bson.M{"phonenumber": payload.Id})
            if err != nil || docuser.Role == "" {
                http.Error(w, "Forbidden: Role not found", http.StatusForbidden)
                return
            }

            for _, allowedRole := range allowedRoles {
                if docuser.Role == allowedRole {
                    next(w, r)
                    return
                }
            }

            http.Error(w, "Forbidden: Access Denied", http.StatusForbidden)
        }
    }
}

// MenuHandler - Handler untuk halaman Menu, hanya untuk role user dan dosen
func MenuHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":  "success",
		"message": "Selamat datang di halaman Menu",
		"data":    "Ini adalah halaman yang dapat diakses oleh user dan dosen.",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// AdminHandler - Handler untuk Dashboard Admin, hanya untuk role admin
func AdminHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":  "success",
		"message": "Selamat datang di Dashboard Admin",
		"data":    "Ini adalah halaman yang hanya dapat diakses oleh admin.",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CashierHandler - Handler untuk Dashboard Kasir, hanya untuk role cashier
func CashierHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":  "success",
		"message": "Selamat datang di Dashboard Kasir",
		"data":    "Ini adalah halaman yang hanya dapat diakses oleh cashier.",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

package model

import "time"

// User struct untuk menyimpan data pengguna
type User struct {
    ID          string    `json:"id" bson:"_id,omitempty"` // ID pengguna (otomatis dibuat oleh MongoDB)
    Name        string    `json:"name" bson:"name"`        // Nama pengguna
    NoWhatsapp  string    `json:"noWhatsapp" bson:"noWhatsapp"` // Nomor WhatsApp pengguna
    Email       string    `json:"email" bson:"email"`      // Email pengguna
    Password    string    `json:"password" bson:"password"` // Password pengguna (harus dihash)
    Role        string    `json:"role" bson:"role"`        // Role pengguna (user, admin, dosen, kasir)
    CreatedAt   time.Time `json:"createdAt" bson:"createdAt"` // Waktu pembuatan akun
    UpdatedAt   time.Time `json:"updatedAt" bson:"updatedAt"` // Waktu terakhir update akun
}

// Default role untuk pengguna baru adalah "user"
func NewUser(name, noWhatsapp, email, password string) User {
    return User{
        Name:        name,
        NoWhatsapp:  noWhatsapp,
        Email:       email,
        Password:    password, // Password harus di-hash sebelum disimpan
        Role:        "user",    // Role default
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }
}

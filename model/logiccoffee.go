package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	ID    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name  string             `json:"name,omitempty" bson:"name,omitempty"`
	Image string             `json:"image,omitempty" bson:"image,omitempty"`
}

type Banner struct {
	ID    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name  string             `json:"name,omitempty" bson:"name,omitempty"`
	Image string             `json:"image,omitempty" bson:"image,omitempty"`
}

type Menu struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CategoryID  primitive.ObjectID `json:"category_id,omitempty" bson:"category_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	Image       string             `json:"image,omitempty" bson:"image,omitempty"`
	Price       float64            `json:"price,omitempty" bson:"price,omitempty"`
}

// Order struct untuk menyimpan informasi pesanan
type Order struct {
    ID                primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`              // ID unik pesanan
    QueueNumber       int                `json:"queue_number,omitempty" bson:"queue_number,omitempty"`
    UniqueOrderNumber string             `json:"unique_order_number,omitempty" bson:"unique_order_number,omitempty"`
    UserID            primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`      // ID pengguna yang memesan (jika dari web)
    UserInfo          UserInfo           `json:"user_info,omitempty" bson:"user_info,omitempty"`  // Informasi user (jika dari web)
    Orders            []OrderItem        `json:"orders,omitempty" bson:"orders,omitempty"`        // Daftar item pesanan
    Total             float64            `json:"total,omitempty" bson:"total,omitempty"`          // Total harga pesanan (harga satuan * kuantitas per item)
    PaymentMethod     string             `json:"payment_method,omitempty" bson:"payment_method,omitempty"` // Metode pembayaran (Transfer/COD/QRIS)
    PaymentInfo       string             `json:"payment_info,omitempty" bson:"payment_info,omitempty"` // Informasi pembayaran (misal rekening transfer atau URL QRIS)
    CreatedBy         string             `json:"created_by,omitempty" bson:"created_by,omitempty"` // Nama siapa yang membuat pesanan (user atau kasir)
    CreatedByRole     string             `json:"created_by_role,omitempty" bson:"created_by_role,omitempty"` // Role siapa yang membuat pesanan
    CreatedAt         int64              `json:"created_at,omitempty" bson:"created_at,omitempty"` // Timestamp pembuatan pesanan
}

// UserInfo struct untuk menyimpan informasi pengguna
type UserInfo struct {
    Name     string `json:"name,omitempty" bson:"name,omitempty"`         // Nama pengguna
    Whatsapp string `json:"whatsapp,omitempty" bson:"whatsapp,omitempty"` // Nomor WhatsApp pengguna
    Note     string `json:"note,omitempty" bson:"note,omitempty"`         // Catatan dari pengguna
}

// OrderItem struct untuk menyimpan detail setiap item dalam pesanan
type OrderItem struct {
    MenuID   primitive.ObjectID `json:"menu_id,omitempty" bson:"menu_id,omitempty"`       // ID dari menu
    Quantity int                `json:"quantity,omitempty" bson:"quantity,omitempty"`    // Kuantitas item
}


package controller

import (
    // "io"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gocroot/config"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/helper/watoken"
	"github.com/gocroot/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
    // "github.com/gocroot/helper/ghupload"
)

func formatRupiah(price float64) string {
	formatter := message.NewPrinter(language.Indonesian)
	return formatter.Sprintf("Rp %.2f", price)
}

// CreateMenu - Tambah Menu Baru
// CreateMenu - Tambah Menu Baru
func CreateMenu(respw http.ResponseWriter, req *http.Request) {
	payload, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
	if err != nil {
		at.WriteJSON(respw, http.StatusForbidden, model.Response{
			Status:   "Error: Token Tidak Valid",
			Location: "Decode Token Error",
			Response: err.Error(),
		})
		return
	}
	var menu model.Menu
	if err := json.NewDecoder(req.Body).Decode(&menu); err != nil {
		at.WriteJSON(respw, http.StatusBadRequest, model.Response{
			Status:   "Error: Bad Request",
			Response: err.Error(),
		})
		return
	}
	if menu.Status != "Tersedia" && menu.Status != "Tidak Tersedia" {
		at.WriteJSON(respw, http.StatusBadRequest, model.Response{
			Status:   "Error: Status Tidak Valid",
			Response: "Status harus 'Tersedia' atau 'Tidak Tersedia'",
		})
		return
	}
	newMenu := model.Menu{
		CategoryID:  menu.CategoryID,
		Name:        menu.Name,
		Description: menu.Description,
		Image:       menu.Image,
		Price:       menu.Price,
		Status:      menu.Status,
	}
	insertResult, err := atdb.InsertOneDoc(config.Mongoconn, "menu", newMenu)
	if err != nil {
		at.WriteJSON(respw, http.StatusNotModified, model.Response{
			Status:   "Error: Gagal Insert Database",
			Response: err.Error(),
		})
		return
	}
	newMenu.ID = insertResult
	response := map[string]interface{}{
		"status":  "success",
		"message": "Menu berhasil ditambahkan",
		"user":    payload.Alias,
		"data": map[string]interface{}{
			"id":          newMenu.ID.Hex(),
			"category_id": newMenu.CategoryID.Hex(),
			"name":        newMenu.Name,
			"description": newMenu.Description,
			"image":       newMenu.Image,
			"price":       formatRupiah(newMenu.Price),
			"status":      newMenu.Status,
		},
	}
	at.WriteJSON(respw, http.StatusOK, response)
}

// GetAllMenu - Ambil Semua Data Menu
func GetAllMenu(respw http.ResponseWriter, req *http.Request) {
	data, err := atdb.GetAllDoc[[]model.Menu](config.Mongoconn, "menu", bson.M{})
	if err != nil {
		at.WriteJSON(respw, http.StatusNotFound, model.Response{
			Status:   "Error: Data menu tidak ditemukan",
			Response: err.Error(),
		})
		return
	}

	if len(data) == 0 {
		at.WriteJSON(respw, http.StatusNotFound, model.Response{
			Status: "Error: Data menu kosong",
		})
		return
	}

	var menus []map[string]interface{}
	for _, menu := range data {
		menus = append(menus, map[string]interface{}{
			"id":          menu.ID.Hex(),
			"category_id": menu.CategoryID.Hex(),
			"name":        menu.Name,
			"description": menu.Description,
			"image":       menu.Image,
			"price":       formatRupiah(menu.Price),
			"status":      menu.Status,
		})
	}

	at.WriteJSON(respw, http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Data menu berhasil diambil",
		"data":    menus,
	})
}

// GetMenuByID - Ambil Menu Berdasarkan ID
func GetMenuByID(respw http.ResponseWriter, req *http.Request) {
	pathParts := strings.Split(req.URL.Path, "/")
	menuID := pathParts[len(pathParts)-1]
	if menuID == "" {
		at.WriteJSON(respw, http.StatusBadRequest, model.Response{
			Status: "Error: ID Menu tidak ditemukan di URL",
		})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(menuID)
	if err != nil {
		at.WriteJSON(respw, http.StatusBadRequest, model.Response{
			Status: "Error: ID Menu tidak valid",
		})
		return
	}

	menu, err := atdb.GetOneDoc[model.Menu](config.Mongoconn, "menu", bson.M{"_id": objectID})
	if err != nil {
		at.WriteJSON(respw, http.StatusNotFound, model.Response{
			Status:   "Error: Menu tidak ditemukan",
			Response: err.Error(),
		})
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Menu ditemukan",
		"data": map[string]interface{}{
			"id":          menu.ID.Hex(),
			"category_id": menu.CategoryID.Hex(),
			"name":        menu.Name,
			"description": menu.Description,
			"image":       menu.Image,
			"price":       formatRupiah(menu.Price),
			"status":      menu.Status,
		},
	}
	at.WriteJSON(respw, http.StatusOK, response)
}

func UpdateMenu(respw http.ResponseWriter, req *http.Request) {
    // Ambil token dari header dan decode menggunakan public key
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
    // Ambil ID menu dari URL
    pathParts := strings.Split(req.URL.Path, "/")
    menuID := pathParts[len(pathParts)-1] // Ambil bagian terakhir dari URL
    if menuID == "" {
        var respn model.Response
        respn.Status = "Error: ID Menu tidak ditemukan di URL"
        at.WriteJSON(respw, http.StatusBadRequest, respn)
        return
    }

    // Konversi ID menu ke ObjectID MongoDB
    objectID, err := primitive.ObjectIDFromHex(menuID)
    if err != nil {
        var respn model.Response
        respn.Status = "Error: ID Menu tidak valid"
        at.WriteJSON(respw, http.StatusBadRequest, respn)
        return
    }

    // Decode body langsung ke map
    var requestBody map[string]interface{}
    err = json.NewDecoder(req.Body).Decode(&requestBody)
    if err != nil {
        var respn model.Response
        respn.Status = "Error: Gagal membaca data JSON"
        respn.Response = err.Error()
        at.WriteJSON(respw, http.StatusBadRequest, respn)
        return
    }

    // Periksa apakah body request kosong
    if len(requestBody) == 0 {
        var respn model.Response
        respn.Status = "Error: Tidak ada data untuk diperbarui"
        at.WriteJSON(respw, http.StatusBadRequest, respn)
        return
    }

    // Menyiapkan data untuk update (langsung menggantikan field yang diberikan)
    updateData := bson.M{}
    if name, exists := requestBody["name"]; exists && name != "" {
        updateData["name"] = name
    }
    if description, exists := requestBody["description"]; exists && description != "" {
        updateData["description"] = description
    }
    if image, exists := requestBody["image"]; exists && image != "" {
        updateData["image"] = image
    }
    if price, exists := requestBody["price"]; exists {
        updateData["price"] = price
    }
    if status, exists := requestBody["status"]; exists && status != "" {
        updateData["status"] = status
    }
    if categoryID, exists := requestBody["category_id"]; exists && categoryID != "" {
        catObjectID, err := primitive.ObjectIDFromHex(categoryID.(string))
        if err != nil {
            var respn model.Response
            respn.Status = "Error: ID Category tidak valid"
            respn.Response = err.Error()
            at.WriteJSON(respw, http.StatusBadRequest, respn)
            return
        }
        updateData["category_id"] = catObjectID
    }

    // Jika tidak ada perubahan data, beri respon error
    if len(updateData) == 0 {
        var respn model.Response
        respn.Status = "Error: Tidak ada perubahan yang dilakukan"
        at.WriteJSON(respw, http.StatusNotModified, respn)
        return
    }

    // Update data menu
    result, err := atdb.UpdateOneDoc(config.Mongoconn, "menu", bson.M{"_id": objectID}, updateData)
    if err != nil {
        var respn model.Response
        respn.Status = "Error: Gagal mengupdate menu"
        respn.Response = err.Error()
        at.WriteJSON(respw, http.StatusInternalServerError, respn)
        return
    }

    // Jika tidak ada dokumen yang dimodifikasi, beri respons error
    if result.ModifiedCount == 0 {
        var respn model.Response
        respn.Status = "Error: Tidak ada perubahan yang dilakukan"
        at.WriteJSON(respw, http.StatusNotModified, respn)
        return
    }

    // Respons sukses
    response := map[string]interface{}{
        "status":  "success",
        "message": "Menu berhasil diupdate",
        "data": map[string]interface{}{
            "id":    objectID.Hex(),
            "updatedFields": updateData,
        },
        "updatedBy": payload.Alias,
    }
    at.WriteJSON(respw, http.StatusOK, response)
}

func DeleteMenu(respw http.ResponseWriter, req *http.Request) {
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

    // Ambil ID menu dari URL
    pathParts := strings.Split(req.URL.Path, "/")
    menuID := pathParts[len(pathParts)-1] // Ambil bagian terakhir dari URL
    if menuID == "" {
        var respn model.Response
        respn.Status = "Error: ID Menu tidak ditemukan di URL"
        at.WriteJSON(respw, http.StatusBadRequest, respn)
        return
    }

    // Konversi ID menu ke ObjectID MongoDB
    objectID, err := primitive.ObjectIDFromHex(menuID)
    if err != nil {
        var respn model.Response
        respn.Status = "Error: ID Menu tidak valid"
        at.WriteJSON(respw, http.StatusBadRequest, respn)
        return
    }

    // Hapus data menu berdasarkan ID
    filter := bson.M{"_id": objectID}
    deleteResult, err := atdb.DeleteOneDoc(config.Mongoconn, "menu", filter)
    if err != nil {
        var respn model.Response
        respn.Status = "Error: Gagal menghapus menu"
        respn.Response = err.Error()
        at.WriteJSON(respw, http.StatusInternalServerError, respn)
        return
    }

    if deleteResult.DeletedCount == 0 {
        var respn model.Response
        respn.Status = "Error: Menu tidak ditemukan"
        at.WriteJSON(respw, http.StatusNotFound, respn)
        return
    }

    // Berhasil menghapus menu
    response := map[string]interface{}{
        "status":  "success",
        "message": "Menu berhasil dihapus",
        "user":    payload.Alias,
        "data":    deleteResult,
    }
    at.WriteJSON(respw, http.StatusOK, response)
}

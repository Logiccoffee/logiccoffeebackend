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

func CreateBanner(respw http.ResponseWriter, req *http.Request) {
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

	var banner model.Banner
	if err := json.NewDecoder(req.Body).Decode(&banner); err != nil {
		var respn model.Response
		respn.Status = "Error: Bad Request"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	// Simpan Banner ke Database
	_, err = atdb.InsertOneDoc(config.Mongoconn, "banner", banner)
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
		"message": "Banner berhasil ditambahkan",
		"name":    payload.Alias,
		"data":    banner,
	}
	at.WriteJSON(respw, http.StatusOK, response)
}

func GetAllBanners(respw http.ResponseWriter, req *http.Request) {
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

	// Ambil Data Banner dari Database
	data, err := atdb.GetAllDoc[[]model.Banner](config.Mongoconn, "banner", primitive.M{})
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Data Banner tidak ditemukan"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	// Cek jika Data Kosong
	if len(data) == 0 {
		var respn model.Response
		respn.Status = "Error: Data Banner kosong"
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	// Transformasi Data Banner
	var banners []map[string]interface{}
	for _, banner := range data {
		banners = append(banners, map[string]interface{}{
			"id":    banner.ID,
			"name":  banner.Name,
			"photo": banner.Photo,
		})
	}

	// Response Berhasil dengan Data Banner
	response := map[string]interface{}{
		"status":  "success",
		"message": "Data Banner berhasil diambil",
		"user":    payload.Alias,
		"data":    banners,
	}
	at.WriteJSON(respw, http.StatusOK, response)
}

func GetBannerByID(respw http.ResponseWriter, req *http.Request) {
	// Ambil ID dari query string
	bannerID := req.URL.Query().Get("id")
	if bannerID == "" {
		var respn model.Response
		respn.Status = "Error: ID Banner tidak ditemukan"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	// Convert ID ke ObjectID
	objectID, err := primitive.ObjectIDFromHex(bannerID)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: ID Banner tidak valid"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	// Ambil Data Banner dari Database
	var banner model.Banner
	filter := bson.M{"_id": objectID}
	_, err = atdb.GetOneDoc[model.Banner](config.Mongoconn, "banner", filter)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Banner tidak ditemukan"
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	// Response dengan Data Banner
	response := map[string]interface{}{
		"status":  "success",
		"message": "Banner ditemukan",
		"data":    banner,
	}
	at.WriteJSON(respw, http.StatusOK, response)
}

func UpdateBanner(respw http.ResponseWriter, req *http.Request) {
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
	bannerID := req.URL.Query().Get("id")
	if bannerID == "" {
		var respn model.Response
		respn.Status = "Error: ID Banner tidak ditemukan"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	// Convert ID ke ObjectID
	objectID, err := primitive.ObjectIDFromHex(bannerID)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: ID Banner tidak valid"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	// Ambil Data Banner dari Database
	filter := bson.M{"_id": objectID}
	_, err = atdb.GetOneDoc[model.Banner](config.Mongoconn, "banner", filter)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Banner tidak ditemukan"
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	// Dekode Request Body untuk Update
	var requestBody struct {
		Name  string `json:"name"`
		Photo string `json:"photo"`
	}
	err = json.NewDecoder(req.Body).Decode(&requestBody)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Gagal membaca data JSON"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	// Update Data Banner
	updateData := bson.M{
		"$set": bson.M{
			"name":  requestBody.Name,
			"photo": requestBody.Photo,
		},
	}
	_, err = atdb.UpdateOneDoc(config.Mongoconn, "banner", filter, updateData)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Gagal mengupdate banner"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusNotModified, respn)
		return
	}

	// Response dengan Data Banner yang Diperbarui
	response := map[string]interface{}{
		"status":  "success",
		"message": "Banner berhasil diupdate",
		"data":    requestBody,
		"name":    payload.Alias,
	}
	at.WriteJSON(respw, http.StatusOK, response)
}

func DeleteBanner(respw http.ResponseWriter, req *http.Request) {
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
	bannerID := req.URL.Query().Get("id")
	if bannerID == "" {
		var respn model.Response
		respn.Status = "Error: ID Banner tidak ditemukan"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	// Convert ID ke ObjectID
	objectID, err := primitive.ObjectIDFromHex(bannerID)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: ID Banner tidak valid"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	// Hapus Data Banner dari Database
	filter := bson.M{"_id": objectID}
	deleteResult, err := atdb.DeleteOneDoc(config.Mongoconn, "banner", filter)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Gagal menghapus banner"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusNotModified, respn)
		return
	}

	// Response Sukses
	if deleteResult.DeletedCount == 0 {
		var respn model.Response
		respn.Status = "Error: Banner tidak ditemukan"
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Banner berhasil dihapus",
		"name":    payload.Alias,
	}
	at.WriteJSON(respw, http.StatusOK, response)
}
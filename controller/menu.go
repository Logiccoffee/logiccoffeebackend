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

func CreateMenu(respw http.ResponseWriter, req *http.Request) {
	payload, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))

	if err != nil {
		payload, err = watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))

		if err != nil {
			var respn model.Response
			respn.Status = "Error: Token Tidak Valid"
			respn.Info = at.GetSecretFromHeader(req)
			respn.Location = "Decode Token Error"
			respn.Response = err.Error()
			at.WriteJSON(respw, http.StatusForbidden, respn)
			return
		}
	}

	var menu model.Menu
	if err := json.NewDecoder(req.Body).Decode(&menu); err != nil {
		var respn model.Response
		respn.Status = "Error: Bad Request"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	_, err = atdb.InsertOneDoc(config.Mongoconn, "menu", menu)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Gagal Insert Database"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusNotModified, respn)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Menu berhasil ditambahkan",
		"name":    payload.Alias,
		"data":    menu,
	}
	at.WriteJSON(respw, http.StatusOK, response)
}

func GetAllMenus(respw http.ResponseWriter, req *http.Request) {
	data, err := atdb.GetAllDoc[[]model.Menu](config.Mongoconn, "menu", primitive.M{})
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Data menu tidak ditemukan"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	if len(data) == 0 {
		var respn model.Response
		respn.Status = "Error: Data menu kosong"
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	var menus []map[string]interface{}
	for _, menu := range data {
		menus = append(menus, map[string]interface{}{
			"id":          menu.ID,
			"name":        menu.Name,
			"description": menu.Description,
			"price":       menu.Price,
			"available":   menu.Available,
			"photo":       menu.Photo,
		})
	}

	at.WriteJSON(respw, http.StatusOK, menus)
}

func GetMenuByID(respw http.ResponseWriter, req *http.Request) {
	menuID := req.URL.Query().Get("id")
	if menuID == "" {
		var respn model.Response
		respn.Status = "Error: ID Menu tidak ditemukan"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(menuID)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: ID Menu tidak valid"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	var menu model.Menu
	filter := bson.M{"_id": objectID}
	_, err = atdb.GetOneDoc[model.Menu](config.Mongoconn, "menu", filter)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Menu tidak ditemukan"
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Menu ditemukan",
		"data":    menu,
	}
	at.WriteJSON(respw, http.StatusOK, response)
}

func UpdateMenu(respw http.ResponseWriter, req *http.Request) {
	payload, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))

	if err != nil {
		payload, err = watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))

		if err != nil {
			var respn model.Response
			respn.Status = "Error: Token Tidak Valid"
			respn.Info = at.GetSecretFromHeader(req)
			respn.Location = "Decode Token Error"
			respn.Response = err.Error()
			at.WriteJSON(respw, http.StatusForbidden, respn)
			return
		}
	}

	menuID := req.URL.Query().Get("id")
	if menuID == "" {
		var respn model.Response
		respn.Status = "Error: ID Menu tidak ditemukan"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(menuID)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: ID Menu tidak valid"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	filter := bson.M{"_id": objectID}
	_, err = atdb.GetOneDoc[model.Menu](config.Mongoconn, "menu", filter)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Menu tidak ditemukan"
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	var requestBody struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Available   bool    `json:"available"`
		Photo       string  `json:"photo"`
	}
	err = json.NewDecoder(req.Body).Decode(&requestBody)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Gagal membaca data JSON"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	updateData := bson.M{}
	if requestBody.Name != "" {
		updateData["name"] = requestBody.Name
	}
	if requestBody.Description != "" {
		updateData["description"] = requestBody.Description
	}
	if requestBody.Price != 0 {
		updateData["price"] = requestBody.Price
	}
	updateData["available"] = requestBody.Available
	if requestBody.Photo != "" {
		updateData["photo"] = requestBody.Photo
	}

	update := bson.M{"$set": updateData}
	_, err = atdb.UpdateOneDoc(config.Mongoconn, "menu", filter, update)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Gagal mengupdate menu"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusNotModified, respn)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Menu berhasil diupdate",
		"data":    updateData,
		"name":    payload.Alias,
	}
	at.WriteJSON(respw, http.StatusOK, response)
}

func DeleteMenu(respw http.ResponseWriter, req *http.Request) {
	payload, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))

	if err != nil {
		payload, err = watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))

		if err != nil {
			var respn model.Response
			respn.Status = "Error: Token Tidak Valid"
			respn.Info = at.GetSecretFromHeader(req)
			respn.Location = "Decode Token Error"
			respn.Response = err.Error()
			at.WriteJSON(respw, http.StatusForbidden, respn)
			return
		}
	}

	menuID := req.URL.Query().Get("id")
	if menuID == "" {
		var respn model.Response
		respn.Status = "Error: ID Menu tidak ditemukan"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(menuID)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: ID Menu tidak valid"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	filter := bson.M{"_id": objectID}
	deleteResult, err := atdb.DeleteOneDoc(config.Mongoconn, "menu", filter)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Gagal menghapus Menu"
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

	response := map[string]interface{}{
		"status":  "success",
		"message": "Menu berhasil dihapus",
		"user":    payload.Alias,
		"data":    deleteResult,
	}
	at.WriteJSON(respw, http.StatusOK, response)
}
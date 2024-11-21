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

func CreateProduct(respw http.ResponseWriter, req *http.Request) {
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

	var product model.Product
	if err := json.NewDecoder(req.Body).Decode(&product); err != nil {
		var respn model.Response
		respn.Status = "Error: Bad Request"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	_, err = atdb.InsertOneDoc(config.Mongoconn, "product", product)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Gagal Insert Database"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusNotModified, respn)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Produk berhasil ditambahkan",
		"name":    payload.Alias,
		"data":    product,
	}
	at.WriteJSON(respw, http.StatusOK, response)
}

func GetAllProducts(respw http.ResponseWriter, req *http.Request) {
	data, err := atdb.GetAllDoc[[]model.Product](config.Mongoconn, "product", primitive.M{})
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Data produk tidak ditemukan"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	if len(data) == 0 {
		var respn model.Response
		respn.Status = "Error: Data produk kosong"
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	var products []map[string]interface{}
	for _, product := range data {
		products = append(products, map[string]interface{}{
			"id":          product.ID,
			"name":        product.Name,
			"description": product.Description,
			"price":       product.Price,
			"available":   product.Available,
			"photo":       product.Photo,
		})
	}

	at.WriteJSON(respw, http.StatusOK, products)
}

func GetProductByID(respw http.ResponseWriter, req *http.Request) {
	productID := req.URL.Query().Get("id")
	if productID == "" {
		var respn model.Response
		respn.Status = "Error: ID Product tidak ditemukan"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: ID Product tidak valid"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	var product model.Product
	filter := bson.M{"_id": objectID}
	_, err = atdb.GetOneDoc[model.Product](config.Mongoconn, "product", filter)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Product tidak ditemukan"
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Product ditemukan",
		"data":    product,
	}
	at.WriteJSON(respw, http.StatusOK, response)
}

func UpdateProduct(respw http.ResponseWriter, req *http.Request) {
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

	productID := req.URL.Query().Get("id")
	if productID == "" {
		var respn model.Response
		respn.Status = "Error: ID Product tidak ditemukan"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: ID Product tidak valid"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	filter := bson.M{"_id": objectID}
	_, err = atdb.GetOneDoc[model.Product](config.Mongoconn, "product", filter)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Product tidak ditemukan"
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
	_, err = atdb.UpdateOneDoc(config.Mongoconn, "product", filter, update)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Gagal mengupdate product"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusNotModified, respn)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Product berhasil diupdate",
		"data":    updateData,
		"name":    payload.Alias,
	}
	at.WriteJSON(respw, http.StatusOK, response)
}

func DeleteProduct(respw http.ResponseWriter, req *http.Request) {
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

	productID := req.URL.Query().Get("id")
	if productID == "" {
		var respn model.Response
		respn.Status = "Error: ID Product tidak ditemukan"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: ID Product tidak valid"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	filter := bson.M{"_id": objectID}
	deleteResult, err := atdb.DeleteOneDoc(config.Mongoconn, "product", filter)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Gagal menghapus Product"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusInternalServerError, respn)
		return
	}

	if deleteResult.DeletedCount == 0 {
		var respn model.Response
		respn.Status = "Error: Product tidak ditemukan"
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Product berhasil dihapus",
		"user":    payload.Alias,
		"data":    deleteResult,
	}
	at.WriteJSON(respw, http.StatusOK, response)
}
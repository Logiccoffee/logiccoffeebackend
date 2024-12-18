package controller

import (
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gocroot/config"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/helper/ghupload"
	"github.com/gocroot/helper/watoken"
	"github.com/gocroot/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func formatRupiah(price float64) string {
	formatter := message.NewPrinter(language.Indonesian)
	return formatter.Sprintf("Rp %.2f", price)
}

// upload menu sekalian sama fotonya.
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

	// Parse multipart form for menu data and image
	err = req.ParseMultipartForm(10 << 20) // Max 10 MB
	if err != nil {
		at.WriteJSON(respw, http.StatusBadRequest, model.Response{
			Status:   "Error: Gagal Memproses Form Data",
			Response: err.Error(),
		})
		return
	}

	// Ambil data menu dari form
	categoryID := req.FormValue("category_id")
	name := req.FormValue("name")
	description := req.FormValue("description")
	price := req.FormValue("price")
	status := req.FormValue("status")

	if status != "Tersedia" && status != "Tidak Tersedia" {
		at.WriteJSON(respw, http.StatusBadRequest, model.Response{
			Status:   "Error: Status Tidak Valid",
			Response: "Status harus 'Tersedia' atau 'Tidak Tersedia'",
		})
		return
	}

	// Parse category ID to ObjectID
	categoryObjectID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		at.WriteJSON(respw, http.StatusBadRequest, model.Response{
			Status:   "Error: Category ID Tidak Valid",
			Response: err.Error(),
		})
		return
	}

	// Handle file upload
	var menuImageURL string
	file, header, err := req.FormFile("menuImage")
	if err == nil {
		defer file.Close()
		fileContent, err := io.ReadAll(file)
		if err != nil {
			at.WriteJSON(respw, http.StatusInternalServerError, model.Response{
				Status:   "Error: Gagal Membaca File",
				Response: err.Error(),
			})
			return
		}

		// Upload image to GitHub
		hashedFileName := ghupload.CalculateHash(fileContent) + header.Filename[strings.LastIndex(header.Filename, "."):]
		pathFile := "menuImages/" + hashedFileName
		content, _, err := ghupload.GithubUpload(
			config.GHAccessToken,
			"Audyardha Nasywa Andini",
			"audyardhaandini@gmail.com",
			fileContent,
			"logiccoffee",
			"img",
			pathFile,
			true,
		)
		if err != nil {
			at.WriteJSON(respw, http.StatusInternalServerError, model.Response{
				Status:   "Error: Gagal Mengupload Gambar Ke GitHub",
				Response: err.Error(),
			})
			return
		}
		menuImageURL = *content.Content.HTMLURL
	}

	// Sanitasi input harga untuk menghapus "Rp", ".", atau ","
	sanitizedPrice := strings.ReplaceAll(price, ".", "")
	sanitizedPrice = strings.ReplaceAll(sanitizedPrice, ",", "")
	sanitizedPrice = strings.ReplaceAll(sanitizedPrice, "Rp", "")
	sanitizedPrice = strings.TrimSpace(sanitizedPrice) // Hapus spasi ekstra

	// Convert price to float64
	menuPrice, err := strconv.ParseFloat(sanitizedPrice, 64)
	if err != nil {
		at.WriteJSON(respw, http.StatusBadRequest, model.Response{
			Status:   "Error: Harga Tidak Valid",
			Response: "Harga harus berupa angka valid",
		})
		return
	}

	// Buat menu baru
	newMenu := model.Menu{
		CategoryID:  categoryObjectID,
		Name:        name,
		Description: description,
		Image:       menuImageURL,
		Price:       menuPrice,
		Status:      status,
	}

	// Simpan ke database
	insertResult, err := atdb.InsertOneDoc(config.Mongoconn, "menu", newMenu)
	if err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, model.Response{
			Status:   "Error: Gagal Menyimpan Menu ke Database",
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

	// Menyiapkan FormData
	err = req.ParseMultipartForm(10 << 20) // 10 MB max size for file uploads
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Gagal memparsing FormData"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	// Ambil data dari FormData
	name := req.FormValue("name")
	description := req.FormValue("description")
	price := req.FormValue("price")
	status := req.FormValue("status")
	categoryID := req.FormValue("category_id")

	// Validasi dan siapkan data untuk update
	updateData := bson.M{}

	// Validasi dan setting nama
	if len(name) > 100 {
		var respn model.Response
		respn.Status = "Error: Nama tidak boleh lebih dari 100 karakter."
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}
	updateData["name"] = name

	// Setting description
	if description != "" {
		updateData["description"] = description
	}

	// Validasi dan setting harga
	if price != "" {
		priceFloat, err := strconv.ParseFloat(price, 64)
		if err != nil {
			var respn model.Response
			respn.Status = "Error: Harga harus berupa angka"
			at.WriteJSON(respw, http.StatusBadRequest, respn)
			return
		}
		updateData["price"] = priceFloat
	}

	// Setting status
	validStatuses := map[string]bool{
		"tersedia":       true,
		"tidak tersedia": true,
		"habis":          true,
	}
	if !validStatuses[status] {
		var respn model.Response
		respn.Status = "Error: Status tidak valid. Hanya 'tersedia', 'tidak tersedia', atau 'habis' yang diperbolehkan."
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}
	updateData["status"] = status

	// Setting category ID
	catObjectID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: ID Category tidak valid"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}
	updateData["category_id"] = catObjectID

	// Handle file upload if "menuImage" is provided in request
	if file, header, err := req.FormFile("menuImage"); err == nil {
		defer file.Close()

		// Baca konten file terlebih dahulu
		fileContent, err := io.ReadAll(file)
		if err != nil {
			var respn model.Response
			respn.Status = "Error: Gagal Membaca File Gambar"
			respn.Response = err.Error()
			at.WriteJSON(respw, http.StatusInternalServerError, respn)
			return
		}

		// Validasi tipe file berdasarkan ekstensi
		if !strings.HasSuffix(strings.ToLower(header.Filename), ".jpg") &&
			!strings.HasSuffix(strings.ToLower(header.Filename), ".png") {
			var respn model.Response
			respn.Status = "Error: File gambar harus berupa JPG atau PNG."
			at.WriteJSON(respw, http.StatusBadRequest, respn)
			return
		}

		// Validasi ukuran file
		if len(fileContent) > 2*1024*1024 { // 2 MB
			var respn model.Response
			respn.Status = "Error: Ukuran file gambar tidak boleh lebih dari 2 MB."
			at.WriteJSON(respw, http.StatusBadRequest, respn)
			return
		}

		// Proses upload gambar ke GitHub
		hashedFileName := ghupload.CalculateHash(fileContent) + header.Filename[strings.LastIndex(header.Filename, "."):]
		pathFile := "menuImages/" + hashedFileName
		content, _, err := ghupload.GithubUpload(
			config.GHAccessToken,
			"Audyardha Nasywa Andini",
			"audyardhaandini@gmail.com",
			fileContent,
			"logiccoffee",
			"img",
			pathFile,
			true,
		)
		if err != nil {
			var respn model.Response
			respn.Status = "Error: Gagal Mengupload Gambar Ke GitHub"
			respn.Response = err.Error()
			at.WriteJSON(respw, http.StatusInternalServerError, respn)
			return
		}

		// Simpan URL gambar ke data update
		updateData["image"] = *content.Content.HTMLURL
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
			"id":            objectID.Hex(),
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

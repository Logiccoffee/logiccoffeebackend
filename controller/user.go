package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gocroot/config"
	"github.com/gocroot/model"
	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gocroot/helper/at"
	"github.com/gocroot/helper/atapi"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/helper/gcallapi"
	"github.com/gocroot/helper/lms"
	"github.com/gocroot/helper/report"
	"github.com/gocroot/helper/watoken"
	"github.com/gocroot/helper/whatsauth"
)

func GetDataUserFromApi(respw http.ResponseWriter, req *http.Request) {
	payload, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
	if err != nil {
		var respn model.Response
		respn.Status = "Error : Token Tidak Valid "
		respn.Info = at.GetSecretFromHeader(req)
		respn.Location = "Decode Token Error: " + at.GetLoginFromHeader(req)
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusForbidden, respn)
		return
	}
	userdt := lms.GetDataFromAPI(payload.Id)
	if userdt.Data.Fullname == "" {
		at.WriteJSON(respw, http.StatusNotFound, userdt)
		return
	}
	at.WriteJSON(respw, http.StatusOK, userdt)
}

func GetDataUser(respw http.ResponseWriter, req *http.Request) {
	payload, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
	if err != nil {
		var respn model.Response
		respn.Status = "Error : Token Tidak Valid "
		respn.Info = config.PublicKeyWhatsAuth
		respn.Location = "Decode Token Error: " + at.GetLoginFromHeader(req)
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusForbidden, respn)
		return
	}
	docuser, err := atdb.GetOneDoc[model.Userdomyikado](config.Mongoconn, "user", primitive.M{"phonenumber": payload.Id})
	if err != nil {
		docuser.PhoneNumber = payload.Id
		docuser.Name = payload.Alias
		at.WriteJSON(respw, http.StatusNotFound, docuser)
		return
	}
	docuser.Name = payload.Alias
	at.WriteJSON(respw, http.StatusOK, docuser)
}

// melakukan pengecekan apakah suda link device klo ada generate token 5tahun
func PutTokenDataUser(respw http.ResponseWriter, req *http.Request) {
	payload, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
	if err != nil {
		var respn model.Response
		respn.Status = "Error : Token Tidak Valid "
		respn.Info = at.GetSecretFromHeader(req)
		respn.Location = "Decode Token Error: " + at.GetLoginFromHeader(req)
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusForbidden, respn)
		return
	}
	docuser, err := atdb.GetOneDoc[model.Userdomyikado](config.Mongoconn, "user", primitive.M{"phonenumber": payload.Id})
	if err != nil {
		docuser.PhoneNumber = payload.Id
		docuser.Name = payload.Alias
		at.WriteJSON(respw, http.StatusNotFound, docuser)
		return
	}
	docuser.Name = payload.Alias
	hcode, qrstat, err := atapi.Get[model.QRStatus](config.WAAPIGetDevice + at.GetLoginFromHeader(req))
	if err != nil {
		at.WriteJSON(respw, http.StatusMisdirectedRequest, docuser)
		return
	}
	if hcode == http.StatusOK && !qrstat.Status {
		docuser.LinkedDevice, err = watoken.EncodeforHours(docuser.PhoneNumber, docuser.Name, config.PrivateKey, 43830)
		if err != nil {
			at.WriteJSON(respw, http.StatusFailedDependency, docuser)
			return
		}
	} else {
		docuser.LinkedDevice = ""
	}
	_, err = atdb.ReplaceOneDoc(config.Mongoconn, "user", primitive.M{"phonenumber": payload.Id}, docuser)
	if err != nil {
		at.WriteJSON(respw, http.StatusExpectationFailed, docuser)
		return
	}
	at.WriteJSON(respw, http.StatusOK, docuser)
}

func PostDataUser(respw http.ResponseWriter, req *http.Request) {
	payload, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
	if err != nil {
		var respn model.Response
		respn.Status = "Error : Token Tidak Valid"
		respn.Info = at.GetSecretFromHeader(req)
		respn.Location = "Decode Token Error"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusForbidden, respn)
		return
	}
	var usr model.Userdomyikado
	err = json.NewDecoder(req.Body).Decode(&usr)
	if err != nil {
		var respn model.Response
		respn.Status = "Error : Body tidak valid"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}
	//pengecekan isian usr
	if usr.NIK == "" || usr.Pekerjaan == "" || usr.AlamatRumah == "" || usr.AlamatKantor == "" {
		var respn model.Response
		respn.Status = "Isian tidak lengkap"
		respn.Response = "Mohon isi lengkap NIK, Pekerjaan, dan kedua alamat"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}
	docuser, err := atdb.GetOneDoc[model.Userdomyikado](config.Mongoconn, "user", primitive.M{"phonenumber": payload.Id})
	if err != nil {
		usr.PhoneNumber = payload.Id
		usr.Name = payload.Alias
		idusr, err := atdb.InsertOneDoc(config.Mongoconn, "user", usr)
		if err != nil {
			var respn model.Response
			respn.Status = "Gagal Insert Database"
			respn.Response = err.Error()
			at.WriteJSON(respw, http.StatusNotModified, respn)
			return
		}
		usr.ID = idusr
		at.WriteJSON(respw, http.StatusOK, usr)
		return
	}
	//jika email belum gsign maka gsign dulu
	if docuser.Email == "" {
		var respn model.Response
		respn.Status = "Email belum terdaftar"
		respn.Response = "Mohon lakukan google sign in dahulu agar email bisa terdaftar"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}
	docuser.NIK = usr.NIK
	docuser.Pekerjaan = usr.Pekerjaan
	docuser.AlamatRumah = usr.AlamatRumah
	docuser.AlamatKantor = usr.AlamatKantor
	_, err = atdb.ReplaceOneDoc(config.Mongoconn, "user", primitive.M{"phonenumber": payload.Id}, docuser)
	if err != nil {
		var respn model.Response
		respn.Status = "Gagal replaceonedoc"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusConflict, respn)
		return
	}
	//melakukan update di seluruh member project
	//ambil project yang member sebagai anggota
	existingprjs, err := atdb.GetAllDoc[[]model.Project](config.Mongoconn, "project", primitive.M{"members._id": docuser.ID})
	if err != nil { //kalo belum jadi anggota project manapun aman langsung ok
		at.WriteJSON(respw, http.StatusOK, docuser)
		return
	}
	if len(existingprjs) == 0 { //kalo belum jadi anggota project manapun aman langsung ok
		at.WriteJSON(respw, http.StatusOK, docuser)
		return
	}
	//loop keanggotaan setiap project dan menggantinya dengan doc yang terupdate
	for _, prj := range existingprjs {
		memberToDelete := model.Userdomyikado{PhoneNumber: docuser.PhoneNumber}
		_, err := atdb.DeleteDocFromArray[model.Userdomyikado](config.Mongoconn, "project", prj.ID, "members", memberToDelete)
		if err != nil {
			var respn model.Response
			respn.Status = "Error : Data project tidak di temukan"
			respn.Response = err.Error()
			at.WriteJSON(respw, http.StatusNotFound, respn)
			return
		}
		_, err = atdb.AddDocToArray[model.Userdomyikado](config.Mongoconn, "project", prj.ID, "members", docuser)
		if err != nil {
			var respn model.Response
			respn.Status = "Error : Gagal menambahkan member ke project"
			respn.Response = err.Error()
			at.WriteJSON(respw, http.StatusExpectationFailed, respn)
			return
		}

	}

	at.WriteJSON(respw, http.StatusOK, docuser)
}

func PostDataBioUser(respw http.ResponseWriter, req *http.Request) {
	payload, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
	if err != nil {
		var respn model.Response
		respn.Status = "Error : Token Tidak Valid"
		respn.Info = at.GetSecretFromHeader(req)
		respn.Location = "Decode Token Error"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusForbidden, respn)
		return
	}
	var usr model.Userdomyikado
	err = json.NewDecoder(req.Body).Decode(&usr)
	if err != nil {
		var respn model.Response
		respn.Status = "Error : Body tidak valid"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}
	docuser, err := atdb.GetOneDoc[model.Userdomyikado](config.Mongoconn, "user", primitive.M{"phonenumber": payload.Id})
	if err != nil {
		usr.PhoneNumber = payload.Id
		usr.Name = payload.Alias
		idusr, err := atdb.InsertOneDoc(config.Mongoconn, "user", usr)
		if err != nil {
			var respn model.Response
			respn.Status = "Gagal Insert Database"
			respn.Response = err.Error()
			at.WriteJSON(respw, http.StatusNotModified, respn)
			return
		}
		usr.ID = idusr
		at.WriteJSON(respw, http.StatusOK, usr)
		return
	}
	//check profpic apakah kosong  atau engga
	if docuser.ProfilePicture == "" {
		var respn model.Response
		respn.Status = "Belum ada Profile Picture"
		respn.Response = "Mohon upload dahulu profile picture anda pada form yang disediakan"
		at.WriteJSON(respw, http.StatusConflict, respn)
		return
	}

	//publish ke blog
	postingan := strings.ReplaceAll(config.ProfPost, "##PROFPIC##", docuser.ProfilePicture)
	postingan = strings.ReplaceAll(postingan, "##BIO##", usr.Bio)
	bpost, err := gcallapi.PostToBlogger(config.Mongoconn, docuser.URLBio, "2587271685863777988", docuser.Name, postingan)
	if err != nil {
		var respn model.Response
		respn.Status = "Gagal post ke blogger"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusConflict, respn)
		return
	}
	//update data content
	docuser.URLBio = bpost.Id
	docuser.PATHBio = bpost.Url
	docuser.Bio = usr.Bio
	//update user data
	_, err = atdb.ReplaceOneDoc(config.Mongoconn, "user", primitive.M{"phonenumber": payload.Id}, docuser)
	if err != nil {
		var respn model.Response
		respn.Status = "Gagal replaceonedoc"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusConflict, respn)
		return
	}
	//melakukan update di seluruh member project
	//ambil project yang member sebagai anggota
	existingprjs, err := atdb.GetAllDoc[[]model.Project](config.Mongoconn, "project", primitive.M{"members._id": docuser.ID})
	if err != nil { //kalo belum jadi anggota project manapun aman langsung ok
		at.WriteJSON(respw, http.StatusOK, docuser)
		return
	}
	if len(existingprjs) == 0 { //kalo belum jadi anggota project manapun aman langsung ok
		at.WriteJSON(respw, http.StatusOK, docuser)
		return
	}
	//loop keanggotaan setiap project dan menggantinya dengan doc yang terupdate
	for _, prj := range existingprjs {
		memberToDelete := model.Userdomyikado{PhoneNumber: docuser.PhoneNumber}
		_, err := atdb.DeleteDocFromArray[model.Userdomyikado](config.Mongoconn, "project", prj.ID, "members", memberToDelete)
		if err != nil {
			var respn model.Response
			respn.Status = "Error : Data project tidak di temukan"
			respn.Response = err.Error()
			at.WriteJSON(respw, http.StatusNotFound, respn)
			return
		}
		_, err = atdb.AddDocToArray[model.Userdomyikado](config.Mongoconn, "project", prj.ID, "members", docuser)
		if err != nil {
			var respn model.Response
			respn.Status = "Error : Gagal menambahkan member ke project"
			respn.Response = err.Error()
			at.WriteJSON(respw, http.StatusExpectationFailed, respn)
			return
		}

	}

	at.WriteJSON(respw, http.StatusOK, docuser)
}

func PostDataUserFromWA(respw http.ResponseWriter, req *http.Request) {
	var resp itmodel.Response
	prof, err := whatsauth.GetAppProfile(at.GetParam(req), config.Mongoconn)
	if err != nil {
		resp.Response = err.Error()
		at.WriteJSON(respw, http.StatusBadRequest, resp)
		return
	}
	if at.GetSecretFromHeader(req) != prof.Secret {
		resp.Response = "Salah secret: " + at.GetSecretFromHeader(req)
		at.WriteJSON(respw, http.StatusUnauthorized, resp)
		return
	}
	var usr model.Userdomyikado
	err = json.NewDecoder(req.Body).Decode(&usr)
	if err != nil {
		resp.Response = "Error : Body tidak valid"
		resp.Info = err.Error()
		at.WriteJSON(respw, http.StatusBadRequest, resp)
		return
	}
	docuser, err := atdb.GetOneDoc[model.Userdomyikado](config.Mongoconn, "user", primitive.M{"phonenumber": usr.PhoneNumber})
	if err != nil {
		idusr, err := atdb.InsertOneDoc(config.Mongoconn, "user", usr)
		if err != nil {
			resp.Response = "Gagal Insert Database"
			resp.Info = err.Error()
			at.WriteJSON(respw, http.StatusNotModified, resp)
			return
		}
		resp.Info = idusr.Hex()
		at.WriteJSON(respw, http.StatusOK, resp)
		return
	}
	docuser.Name = usr.Name
	docuser.Email = usr.Email
	_, err = atdb.ReplaceOneDoc(config.Mongoconn, "user", primitive.M{"phonenumber": usr.PhoneNumber}, docuser)
	if err != nil {
		resp.Response = "Gagal replaceonedoc"
		resp.Info = err.Error()
		at.WriteJSON(respw, http.StatusConflict, resp)
		return
	}
	//melakukan update di seluruh member project
	//ambil project yang member sebagai anggota
	existingprjs, err := atdb.GetAllDoc[[]model.Project](config.Mongoconn, "project", primitive.M{"members._id": docuser.ID})
	if err != nil { //kalo belum jadi anggota project manapun aman langsung ok
		resp.Response = "belum terdaftar di project manapun"
		at.WriteJSON(respw, http.StatusOK, resp)
		return
	}
	if len(existingprjs) == 0 { //kalo belum jadi anggota project manapun aman langsung ok
		resp.Response = "belum terdaftar di project manapun"
		at.WriteJSON(respw, http.StatusOK, resp)
		return
	}
	//loop keanggotaan setiap project dan menggantinya dengan doc yang terupdate
	for _, prj := range existingprjs {
		memberToDelete := model.Userdomyikado{PhoneNumber: docuser.PhoneNumber}
		_, err := atdb.DeleteDocFromArray[model.Userdomyikado](config.Mongoconn, "project", prj.ID, "members", memberToDelete)
		if err != nil {
			resp.Response = "Error : Data project tidak di temukan"
			resp.Info = err.Error()
			at.WriteJSON(respw, http.StatusNotFound, resp)
			return
		}
		_, err = atdb.AddDocToArray[model.Userdomyikado](config.Mongoconn, "project", prj.ID, "members", docuser)
		if err != nil {
			resp.Response = "Error : Gagal menambahkan member ke project"
			resp.Info = err.Error()
			at.WriteJSON(respw, http.StatusExpectationFailed, resp)
			return
		}

	}
	resp.Info = docuser.ID.Hex()
	resp.Info = docuser.Email
	at.WriteJSON(respw, http.StatusOK, resp)
}

func ApproveBimbinganbyPoin(w http.ResponseWriter, r *http.Request) {
	noHp := r.Header.Get("nohp")
	if noHp == "" {
		http.Error(w, "No valid phone number found", http.StatusForbidden)
		return
	}

	var requestData struct {
		NIM   string `json:"nim"`
		Topik string `json:"topik"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil || requestData.NIM == "" || requestData.Topik == "" {
		http.Error(w, "Invalid request body or NIM/Topik not provided", http.StatusBadRequest)
		return
	}

	// Get the API URL from the database
	var conf model.Config
	err = config.Mongoconn.Collection("config").FindOne(context.TODO(), bson.M{"phonenumber": "62895601060000"}).Decode(&conf)
	if err != nil {
		http.Error(w, "Mohon maaf ada kesalahan dalam pengambilan config di database: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare the request body
	requestBody, err := json.Marshal(map[string]string{
		"nim":   requestData.NIM,
		"topik": requestData.Topik,
	})
	if err != nil {
		http.Error(w, "Gagal membuat request body: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create and send the HTTP request
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", conf.ApproveBimbinganURL, bytes.NewBuffer(requestBody))
	if err != nil {
		http.Error(w, "Gagal membuat request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("nohp", noHp)

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Gagal mengirim request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusNotFound:
			http.Error(w, "Token tidak ditemukan! Silahkan Login Kembali", http.StatusNotFound)
		case http.StatusForbidden:
			http.Error(w, "Gagal, Bimbingan telah disetujui!", http.StatusForbidden)
		default:
			http.Error(w, fmt.Sprintf("Gagal approve bimbingan, status code: %d", resp.StatusCode), http.StatusInternalServerError)
		}
		return
	}

	var responseMap map[string]string
	err = json.NewDecoder(resp.Body).Decode(&responseMap)
	if err != nil {
		http.Error(w, "Gagal memproses response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Kurangi poin berdasarkan nomor telepon yang ada di response
	phonenumber := responseMap["no_hp"]
	_, err = report.KurangPoinUserbyPhoneNumber(config.Mongoconn, phonenumber, 13.0)
	if err != nil {
		http.Error(w, "Gagal mengurangi poin: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get updated user data to return the current points
	usr, err := atdb.GetOneDoc[model.Userdomyikado](config.Mongoconn, "user", bson.M{"phonenumber": phonenumber})
	if err != nil {
		http.Error(w, "Gagal mengambil data pengguna: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Add the current points to the response
	responseMap["message"] = "Bimbingan berhasil di approve!"
	responseMap["status"] = "success"
	responseMap["poin_mahasiswa"] = fmt.Sprintf("Poin mahasiswa telah berkurang menjadi: %f", usr.Poin)

	at.WriteJSON(w, http.StatusOK, responseMap)
}

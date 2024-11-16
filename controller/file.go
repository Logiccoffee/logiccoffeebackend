package controller

import (
	"encoding/base64"
	"io"
	"net/http"
	"strings"

	"github.com/gocroot/config"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/helper/dokped"
	"github.com/gocroot/helper/ghupload"
	"github.com/gocroot/helper/watoken"
	"github.com/gocroot/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AksesFileRepoDraft(w http.ResponseWriter, r *http.Request) {
	var respn model.Response
	payload, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(r))
	if err != nil {
		respn.Status = "Error : Token Tidak Valid"
		respn.Info = at.GetSecretFromHeader(r)
		respn.Location = "Decode Token Error"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusForbidden, respn)
		return
	}
	docuser, err := atdb.GetOneDoc[model.Userdomyikado](config.Mongoconn, "user", primitive.M{"phonenumber": payload.Id})
	if err != nil {
		respn.Status = "Error : Data user tidak di temukan"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusNotImplemented, respn)
		return
	}
	pathFileBase64 := at.GetParam(r)
	// Decode string dari Base64
	decoded, err := base64.StdEncoding.DecodeString(pathFileBase64)
	if err != nil {
		respn.Status = "Error : decoding base64"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusNotImplemented, respn)
		return
	}
	pathFile := string(decoded)
	pathslice := strings.Split(pathFile, "/")
	namaprj := pathslice[0]
	//cek apakah user memiliki akses ke project
	prj, err := atdb.GetOneDoc[model.Project](config.Mongoconn, "project", primitive.M{"name": namaprj})
	if err != nil {
		respn.Status = "Error : Data lapak tidak di temukan"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusNotImplemented, respn)
		return
	}
	//check apakah dia owner
	if (prj.Owner.PhoneNumber != docuser.PhoneNumber) && (prj.Editor.PhoneNumber != docuser.PhoneNumber) {
		if !docuser.IsManager { //kalo bukan manager maka ga punya akses dong
			respn.Status = "Error : User bukan owner project tidak berhak"
			respn.Response = "User bukan owner dari project ini"
			at.WriteJSON(w, http.StatusNotImplemented, respn)
			return
		}
	}

	githubOrg := "penerbitbukupedia"
	githubRepo := "draft"
	filecontent, err := ghupload.GithubGetFile(config.GHAccessToken, githubOrg, githubRepo, pathFile)
	if err != nil {
		respn.Status = "Error : Data tidak bisa diambil dari github"
		respn.Info = githubOrg + "/" + githubRepo
		respn.Location = pathFile
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusBadRequest, respn)
		return
	}
	at.WriteFile(w, http.StatusOK, filecontent)
}

func GetFileDraftSPK(w http.ResponseWriter, r *http.Request) {
	var respn model.Response
	payload, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(r))
	if err != nil {
		respn.Status = "Error : Token Tidak Valid"
		respn.Info = at.GetSecretFromHeader(r)
		respn.Location = "Decode Token Error"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusForbidden, respn)
		return
	}
	docuser, err := atdb.GetOneDoc[model.Userdomyikado](config.Mongoconn, "user", primitive.M{"phonenumber": payload.Id})
	if err != nil {
		respn.Status = "Error : Data user tidak di temukan"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusNotImplemented, respn)
		return
	}
	pathFileBase64 := at.GetParam(r)
	// Decode string dari Base64
	decoded, err := base64.StdEncoding.DecodeString(pathFileBase64)
	if err != nil {
		respn.Status = "Error : decoding base64"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusNotImplemented, respn)
		return
	}
	namaprj := string(decoded)
	//cek apakah user memiliki akses ke project
	prj, err := atdb.GetOneDoc[model.Project](config.Mongoconn, "project", primitive.M{"name": namaprj})
	if err != nil {
		respn.Status = "Error : Data lapak tidak di temukan"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusNotImplemented, respn)
		return
	}
	//check apakah dia owner
	if prj.Owner.PhoneNumber != docuser.PhoneNumber {
		respn.Status = "Error : User bukan owner project tidak berhak"
		respn.Response = "User bukan owner dari project ini"
		at.WriteJSON(w, http.StatusNotImplemented, respn)
		return
	}
	filecontent, err := dokped.GenerateSPK(prj, config.AESKey)
	if err != nil {
		respn.Status = "Error : Dokumen gagal di generate"
		respn.Info = prj.Name
		respn.Location = prj.ID.Hex()
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusBadRequest, respn)
		return
	}
	at.WriteFile(w, http.StatusOK, filecontent)
}

func GetFileDraftSPKT(w http.ResponseWriter, r *http.Request) {
	var respn model.Response
	payload, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(r))
	if err != nil {
		respn.Status = "Error : Token Tidak Valid"
		respn.Info = at.GetSecretFromHeader(r)
		respn.Location = "Decode Token Error"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusForbidden, respn)
		return
	}
	docuser, err := atdb.GetOneDoc[model.Userdomyikado](config.Mongoconn, "user", primitive.M{"phonenumber": payload.Id})
	if err != nil {
		respn.Status = "Error : Data user tidak di temukan"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusNotImplemented, respn)
		return
	}
	pathFileBase64 := at.GetParam(r)
	// Decode string dari Base64
	decoded, err := base64.StdEncoding.DecodeString(pathFileBase64)
	if err != nil {
		respn.Status = "Error : decoding base64"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusNotImplemented, respn)
		return
	}
	namaprj := string(decoded)
	//cek apakah user memiliki akses ke project
	prj, err := atdb.GetOneDoc[model.Project](config.Mongoconn, "project", primitive.M{"name": namaprj})
	if err != nil {
		respn.Status = "Error : Data lapak tidak di temukan"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusNotImplemented, respn)
		return
	}
	//check apakah dia owner
	if prj.Owner.PhoneNumber != docuser.PhoneNumber {
		respn.Status = "Error : User bukan owner project tidak berhak"
		respn.Response = "User bukan owner dari project ini"
		at.WriteJSON(w, http.StatusNotImplemented, respn)
		return
	}
	filecontent, err := dokped.GenerateSPKT(prj, config.AESKey)
	if err != nil {
		respn.Status = "Error : Dokumen gagal di generate"
		respn.Info = prj.Name
		respn.Location = prj.ID.Hex()
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusBadRequest, respn)
		return
	}
	at.WriteFile(w, http.StatusOK, filecontent)
}

func GetFileDraftSPI(w http.ResponseWriter, r *http.Request) {
	var respn model.Response
	payload, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(r))
	if err != nil {
		respn.Status = "Error : Token Tidak Valid"
		respn.Info = at.GetSecretFromHeader(r)
		respn.Location = "Decode Token Error"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusForbidden, respn)
		return
	}
	docuser, err := atdb.GetOneDoc[model.Userdomyikado](config.Mongoconn, "user", primitive.M{"phonenumber": payload.Id})
	if err != nil {
		respn.Status = "Error : Data user tidak di temukan"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusNotImplemented, respn)
		return
	}
	pathFileBase64 := at.GetParam(r)
	// Decode string dari Base64
	decoded, err := base64.StdEncoding.DecodeString(pathFileBase64)
	if err != nil {
		respn.Status = "Error : decoding base64"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusNotImplemented, respn)
		return
	}
	pathFile := string(decoded)
	pathslice := strings.Split(pathFile, "/")
	namaprj := pathslice[0]
	//cek apakah user memiliki akses ke project
	prj, err := atdb.GetOneDoc[model.Project](config.Mongoconn, "project", primitive.M{"name": namaprj})
	if err != nil {
		respn.Status = "Error : Data lapak tidak di temukan"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusNotImplemented, respn)
		return
	}
	//check apakah dia owner
	if prj.Owner.PhoneNumber != docuser.PhoneNumber {
		respn.Status = "Error : User bukan owner project tidak berhak"
		respn.Response = "User bukan owner dari project ini"
		at.WriteJSON(w, http.StatusNotImplemented, respn)
		return
	}
	//ambil surat pengantar
	filecontentpengantar, err := dokped.GenerateSPI(prj, config.AESKey)
	if err != nil {
		respn.Status = "Error : Dokumen gagal di generate"
		respn.Info = prj.Name
		respn.Location = prj.ID.Hex()
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusBadRequest, respn)
		return
	}
	//gabungkan dengan pdf sampul
	/* githubOrg := "penerbitbukupedia"
	githubRepo := "draft"
	filecontentsampul, err := ghupload.GithubGetFile(config.GHAccessToken, githubOrg, githubRepo, pathFile)
	if err != nil {
		respn.Status = "Error : Data tidak bisa diambil dari github"
		respn.Info = githubOrg + "/" + githubRepo
		respn.Location = pathFile
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusBadRequest, respn)
		return
	}
	filecontent, err := fpdf.MergePDFBytes(filecontentpengantar, filecontentsampul)
	if err != nil {
		respn.Status = "Error : Dokumen gagal di merge"
		respn.Info = prj.Name
		respn.Location = prj.ID.Hex()
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusBadRequest, respn)
		return
	} */
	at.WriteFile(w, http.StatusOK, filecontentpengantar)
}

func UploadProfilePictureHandler(w http.ResponseWriter, r *http.Request) {
	var respn model.Response
	payload, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(r))
	if err != nil {
		respn.Status = "Error : Token Tidak Valid"
		respn.Info = at.GetSecretFromHeader(r)
		respn.Location = "Decode Token Error"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusForbidden, respn)
		return
	}
	userdoc, err := atdb.GetOneDoc[model.Userdomyikado](config.Mongoconn, "user", primitive.M{"phonenumber": payload.Id})
	if err != nil {
		respn.Status = "Error : Data user tidak di temukan"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusNotImplemented, respn)
		return
	}
	file, header, err := r.FormFile("profpic")
	if err != nil {
		respn.Status = "Error : File tidak ada"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusBadRequest, respn)
		return
	}
	defer file.Close()
	// Read the file content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		respn.Status = "Error : File tidak bisa dibaca"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusInternalServerError, respn)
		return
	}
	// Calculate hash of the file content
	//hashedFileName := ghupload.CalculateHash(fileContent)

	// Get GitHub credentials and other details from the request or environment variables
	GitHubAccessToken := config.GHAccessToken
	GitHubAuthorName := "Rolly Maulana Awangga"
	GitHubAuthorEmail := "awangga@gmail.com"
	githubOrg := "penerbitbukupedia"
	githubRepo := "profile"
	pathFile := "picture/" + userdoc.ID.Hex() + "/" + userdoc.ID.Hex() + header.Filename[strings.LastIndex(header.Filename, "."):] // Append the original file extension
	replace := true

	// Use GithubUpload function to upload the file to GitHub
	content, _, err := ghupload.GithubUpload(GitHubAccessToken, GitHubAuthorName, GitHubAuthorEmail, fileContent, githubOrg, githubRepo, pathFile, replace)
	if err != nil {
		respn.Status = "Error : File tidak bisa diupload ke github"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusInternalServerError, respn)
		return
	}
	//update data profpic
	userdoc.ProfilePicture = "https://raw.githubusercontent.com/" + githubOrg + "/" + githubRepo + "/main/" + *content.Content.Path
	atdb.ReplaceOneDoc(config.Mongoconn, "user", bson.M{"_id": userdoc.ID}, userdoc)

	// Respond with success message
	respn.Info = userdoc.ID.Hex()
	respn.Location = userdoc.ProfilePicture
	respn.Response = *content.Content.URL
	respn.Status = *content.Content.HTMLURL
	at.WriteJSON(w, http.StatusOK, respn)
}

func FileUploadWithParamFileHandler(w http.ResponseWriter, r *http.Request) {
	var respn model.Response
	payload, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(r))
	if err != nil {
		respn.Status = "Error : Token Tidak Valid"
		respn.Info = at.GetSecretFromHeader(r)
		respn.Location = "Decode Token Error"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusForbidden, respn)
		return
	}
	_, err = atdb.GetOneDoc[model.Userdomyikado](config.Mongoconn, "user", primitive.M{"phonenumber": payload.Id})
	if err != nil {
		respn.Status = "Error : Data user tidak di temukan"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusNotImplemented, respn)
		return
	}
	prjid := at.GetParam(r)
	objectId, _ := primitive.ObjectIDFromHex(prjid)
	prj, err := atdb.GetOneDoc[model.Project](config.Mongoconn, "project", primitive.M{"_id": objectId})
	if err != nil {
		respn.Status = "Error : Data lapak tidak di temukan"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusNotImplemented, respn)
		return
	}

	file, header, err := r.FormFile("profpic")
	if err != nil {
		respn.Status = "Error : File tidak ada"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusBadRequest, respn)
		return
	}
	defer file.Close()
	// Read the file content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		respn.Status = "Error : File tidak bisa dibaca"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusInternalServerError, respn)
		return
	}
	// Calculate hash of the file content
	hashedFileName := ghupload.CalculateHash(fileContent)

	// Get GitHub credentials and other details from the request or environment variables
	GitHubAccessToken := config.GHAccessToken
	GitHubAuthorName := "Rolly Maulana Awangga"
	GitHubAuthorEmail := "awangga@gmail.com"
	githubOrg := "penerbitbukupedia"
	githubRepo := "img"
	pathFile := prj.Name + "/menu/" + hashedFileName + header.Filename[strings.LastIndex(header.Filename, "."):] // Append the original file extension
	replace := true

	// Use GithubUpload function to upload the file to GitHub
	content, _, err := ghupload.GithubUpload(GitHubAccessToken, GitHubAuthorName, GitHubAuthorEmail, fileContent, githubOrg, githubRepo, pathFile, replace)
	if err != nil {
		respn.Status = "Error : File tidak bisa diupload ke github"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusInternalServerError, respn)
		return
	}

	// Respond with success message
	respn.Info = hashedFileName
	respn.Location = "/" + githubRepo + "/" + *content.Content.Path
	respn.Response = *content.Content.URL
	respn.Status = *content.Content.HTMLURL
	at.WriteJSON(w, http.StatusOK, respn)
}

func UploadCoverBukuWithParamFileHandler(w http.ResponseWriter, r *http.Request) {
	var respn model.Response
	payload, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(r))
	if err != nil {
		respn.Status = "Error : Token Tidak Valid"
		respn.Info = at.GetSecretFromHeader(r)
		respn.Location = "Decode Token Error"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusForbidden, respn)
		return
	}
	docuser, err := atdb.GetOneDoc[model.Userdomyikado](config.Mongoconn, "user", primitive.M{"phonenumber": payload.Id})
	if err != nil {
		respn.Status = "Error : Data user tidak di temukan"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusNotImplemented, respn)
		return
	}
	prjid := at.GetParam(r)
	objectId, _ := primitive.ObjectIDFromHex(prjid)
	prj, err := atdb.GetOneDoc[model.Project](config.Mongoconn, "project", primitive.M{"_id": objectId})
	if err != nil {
		respn.Status = "Error : Data lapak tidak di temukan"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusNotImplemented, respn)
		return
	}
	//check apakah dia owner
	if prj.Owner.PhoneNumber != docuser.PhoneNumber {
		respn.Status = "Error : User bukan owner project tidak berhak"
		respn.Response = "User bukan owner dari project ini"
		at.WriteJSON(w, http.StatusNotImplemented, respn)
		return
	}

	file, header, err := r.FormFile("coverbuku")
	if err != nil {
		respn.Status = "Error : File tidak ada"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusBadRequest, respn)
		return
	}
	defer file.Close()
	// Read the file content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		respn.Status = "Error : File tidak bisa dibaca"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusInternalServerError, respn)
		return
	}
	// Calculate hash of the file content
	//hashedFileName := ghupload.CalculateHash(fileContent)

	// Get GitHub credentials and other details from the request or environment variables
	GitHubAccessToken := config.GHAccessToken
	GitHubAuthorName := "Rolly Maulana Awangga"
	GitHubAuthorEmail := "awangga@gmail.com"
	githubOrg := "penerbitbukupedia"
	githubRepo := "katalog"
	pathFile := prj.Name + "/cover/" + prj.ID.Hex() + header.Filename[strings.LastIndex(header.Filename, "."):] // Append the original file extension
	replace := true

	// Use GithubUpload function to upload the file to GitHub
	content, _, err := ghupload.GithubUpload(GitHubAccessToken, GitHubAuthorName, GitHubAuthorEmail, fileContent, githubOrg, githubRepo, pathFile, replace)
	if err != nil {
		respn.Status = "Error : File tidak bisa diupload ke github"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusInternalServerError, respn)
		return
	}
	//update data profpic
	prj.CoverBuku = "https://raw.githubusercontent.com/" + githubOrg + "/" + githubRepo + "/main/" + *content.Content.Path
	atdb.ReplaceOneDoc(config.Mongoconn, "project", bson.M{"_id": prj.ID}, prj)

	// Respond with success message
	respn.Info = prj.ID.Hex()
	respn.Location = prj.CoverBuku
	respn.Response = *content.Content.URL
	respn.Status = *content.Content.HTMLURL
	at.WriteJSON(w, http.StatusOK, respn)
}

func UploadDraftBukuWithParamFileHandler(w http.ResponseWriter, r *http.Request) {
	var respn model.Response
	payload, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(r))
	if err != nil {
		respn.Status = "Error : Token Tidak Valid"
		respn.Info = at.GetSecretFromHeader(r)
		respn.Location = "Decode Token Error"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusForbidden, respn)
		return
	}
	docuser, err := atdb.GetOneDoc[model.Userdomyikado](config.Mongoconn, "user", primitive.M{"phonenumber": payload.Id})
	if err != nil {
		respn.Status = "Error : Data user tidak di temukan"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusNotImplemented, respn)
		return
	}
	prjid := at.GetParam(r)
	objectId, _ := primitive.ObjectIDFromHex(prjid)
	prj, err := atdb.GetOneDoc[model.Project](config.Mongoconn, "project", primitive.M{"_id": objectId})
	if err != nil {
		respn.Status = "Error : Data lapak tidak di temukan"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusNotImplemented, respn)
		return
	}
	//check apakah dia owner
	if prj.Owner.PhoneNumber != docuser.PhoneNumber {
		respn.Status = "Error : User bukan owner project tidak berhak"
		respn.Response = "User bukan owner dari project ini"
		at.WriteJSON(w, http.StatusNotImplemented, respn)
		return
	}

	file, header, err := r.FormFile("draftbuku")
	if err != nil {
		respn.Status = "Error : File tidak ada"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusBadRequest, respn)
		return
	}
	defer file.Close()
	// Read the file content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		respn.Status = "Error : File tidak bisa dibaca"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusInternalServerError, respn)
		return
	}
	// Calculate hash of the file content
	//hashedFileName := ghupload.CalculateHash(fileContent)

	// Get GitHub credentials and other details from the request or environment variables
	GitHubAccessToken := config.GHAccessToken
	GitHubAuthorName := "Rolly Maulana Awangga"
	GitHubAuthorEmail := "awangga@gmail.com"
	githubOrg := "penerbitbukupedia"
	githubRepo := "draft"
	pathFile := prj.Name + "/draft/" + prj.ID.Hex() + header.Filename[strings.LastIndex(header.Filename, "."):] // Append the original file extension
	replace := true

	// Use GithubUpload function to upload the file to GitHub
	content, _, err := ghupload.GithubUpload(GitHubAccessToken, GitHubAuthorName, GitHubAuthorEmail, fileContent, githubOrg, githubRepo, pathFile, replace)
	if err != nil {
		respn.Status = "Error : File tidak bisa diupload ke github"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusInternalServerError, respn)
		return
	}
	//update data profpic
	//prj.DraftBuku = "https://raw.githubusercontent.com/" + githubOrg + "/" + githubRepo + "/main/" + *content.Content.Path
	prj.DraftBuku = *content.Content.Path //karena repo private
	atdb.ReplaceOneDoc(config.Mongoconn, "project", bson.M{"_id": prj.ID}, prj)

	// Respond with success message
	respn.Info = prj.ID.Hex()
	respn.Location = prj.DraftBuku
	respn.Response = *content.Content.URL
	respn.Status = *content.Content.HTMLURL
	at.WriteJSON(w, http.StatusOK, respn)
}

func UploadDraftBukuPDFWithParamFileHandler(w http.ResponseWriter, r *http.Request) {
	var respn model.Response
	payload, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(r))
	if err != nil {
		respn.Status = "Error : Token Tidak Valid"
		respn.Info = at.GetSecretFromHeader(r)
		respn.Location = "Decode Token Error"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusForbidden, respn)
		return
	}
	docuser, err := atdb.GetOneDoc[model.Userdomyikado](config.Mongoconn, "user", primitive.M{"phonenumber": payload.Id})
	if err != nil {
		respn.Status = "Error : Data user tidak di temukan"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusNotImplemented, respn)
		return
	}
	prjid := at.GetParam(r)
	objectId, _ := primitive.ObjectIDFromHex(prjid)
	prj, err := atdb.GetOneDoc[model.Project](config.Mongoconn, "project", primitive.M{"_id": objectId})
	if err != nil {
		respn.Status = "Error : Data lapak tidak di temukan"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusNotImplemented, respn)
		return
	}
	//check apakah dia owner
	if prj.Owner.PhoneNumber != docuser.PhoneNumber {
		respn.Status = "Error : User bukan owner project tidak berhak"
		respn.Response = "User bukan owner dari project ini"
		at.WriteJSON(w, http.StatusNotImplemented, respn)
		return
	}

	file, header, err := r.FormFile("draftpdfbuku")
	if err != nil {
		respn.Status = "Error : File tidak ada"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusBadRequest, respn)
		return
	}
	defer file.Close()
	// Read the file content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		respn.Status = "Error : File tidak bisa dibaca"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusInternalServerError, respn)
		return
	}
	// Calculate hash of the file content
	//hashedFileName := ghupload.CalculateHash(fileContent)

	// Get GitHub credentials and other details from the request or environment variables
	GitHubAccessToken := config.GHAccessToken
	GitHubAuthorName := "Rolly Maulana Awangga"
	GitHubAuthorEmail := "awangga@gmail.com"
	githubOrg := "penerbitbukupedia"
	githubRepo := "draft"
	pathFile := prj.Name + "/pdf/" + prj.ID.Hex() + header.Filename[strings.LastIndex(header.Filename, "."):] // Append the original file extension
	replace := true

	// Use GithubUpload function to upload the file to GitHub
	content, _, err := ghupload.GithubUpload(GitHubAccessToken, GitHubAuthorName, GitHubAuthorEmail, fileContent, githubOrg, githubRepo, pathFile, replace)
	if err != nil {
		respn.Status = "Error : File tidak bisa diupload ke github"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusInternalServerError, respn)
		return
	}
	//update data profpic
	prj.DraftPDFBuku = *content.Content.Path
	atdb.ReplaceOneDoc(config.Mongoconn, "project", bson.M{"_id": prj.ID}, prj)

	// Respond with success message
	respn.Info = prj.ID.Hex()
	respn.Location = prj.DraftPDFBuku
	respn.Response = *content.Content.URL
	respn.Status = *content.Content.HTMLURL
	at.WriteJSON(w, http.StatusOK, respn)
}

func UploadSPKPDFWithParamFileHandler(w http.ResponseWriter, r *http.Request) {
	var respn model.Response
	payload, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(r))
	if err != nil {
		respn.Status = "Error : Token Tidak Valid"
		respn.Info = at.GetSecretFromHeader(r)
		respn.Location = "Decode Token Error"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusForbidden, respn)
		return
	}
	docuser, err := atdb.GetOneDoc[model.Userdomyikado](config.Mongoconn, "user", primitive.M{"phonenumber": payload.Id})
	if err != nil {
		respn.Status = "Error : Data user tidak di temukan"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusNotImplemented, respn)
		return
	}
	prjid := at.GetParam(r)
	objectId, _ := primitive.ObjectIDFromHex(prjid)
	prj, err := atdb.GetOneDoc[model.Project](config.Mongoconn, "project", primitive.M{"_id": objectId})
	if err != nil {
		respn.Status = "Error : Data lapak tidak di temukan"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusNotImplemented, respn)
		return
	}
	//check apakah dia owner
	if prj.Owner.PhoneNumber != docuser.PhoneNumber {
		respn.Status = "Error : User bukan owner project tidak berhak"
		respn.Response = "User bukan owner dari project ini"
		at.WriteJSON(w, http.StatusNotImplemented, respn)
		return
	}

	file, header, err := r.FormFile("spk")
	if err != nil {
		respn.Status = "Error : File tidak ada"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusBadRequest, respn)
		return
	}
	defer file.Close()
	// Read the file content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		respn.Status = "Error : File tidak bisa dibaca"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusInternalServerError, respn)
		return
	}
	// Calculate hash of the file content
	//hashedFileName := ghupload.CalculateHash(fileContent)

	// Get GitHub credentials and other details from the request or environment variables
	GitHubAccessToken := config.GHAccessToken
	GitHubAuthorName := "Rolly Maulana Awangga"
	GitHubAuthorEmail := "awangga@gmail.com"
	githubOrg := "penerbitbukupedia"
	githubRepo := "draft"
	pathFile := prj.Name + "/spk/" + prj.ID.Hex() + header.Filename[strings.LastIndex(header.Filename, "."):] // Append the original file extension
	replace := true

	// Use GithubUpload function to upload the file to GitHub
	content, _, err := ghupload.GithubUpload(GitHubAccessToken, GitHubAuthorName, GitHubAuthorEmail, fileContent, githubOrg, githubRepo, pathFile, replace)
	if err != nil {
		respn.Status = "Error : File tidak bisa diupload ke github"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusInternalServerError, respn)
		return
	}
	//update data profpic
	prj.SPK = *content.Content.Path
	atdb.ReplaceOneDoc(config.Mongoconn, "project", bson.M{"_id": prj.ID}, prj)

	// Respond with success message
	respn.Info = prj.ID.Hex()
	respn.Location = prj.SPK
	respn.Response = *content.Content.URL
	respn.Status = *content.Content.HTMLURL
	at.WriteJSON(w, http.StatusOK, respn)
}

func UploadSPIPDFWithParamFileHandler(w http.ResponseWriter, r *http.Request) {
	var respn model.Response
	payload, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(r))
	if err != nil {
		respn.Status = "Error : Token Tidak Valid"
		respn.Info = at.GetSecretFromHeader(r)
		respn.Location = "Decode Token Error"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusForbidden, respn)
		return
	}
	docuser, err := atdb.GetOneDoc[model.Userdomyikado](config.Mongoconn, "user", primitive.M{"phonenumber": payload.Id})
	if err != nil {
		respn.Status = "Error : Data user tidak di temukan"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusNotImplemented, respn)
		return
	}
	prjid := at.GetParam(r)
	objectId, _ := primitive.ObjectIDFromHex(prjid)
	prj, err := atdb.GetOneDoc[model.Project](config.Mongoconn, "project", primitive.M{"_id": objectId})
	if err != nil {
		respn.Status = "Error : Data lapak tidak di temukan"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusNotImplemented, respn)
		return
	}
	//check apakah dia owner
	if prj.Owner.PhoneNumber != docuser.PhoneNumber {
		respn.Status = "Error : User bukan owner project tidak berhak"
		respn.Response = "User bukan owner dari project ini"
		at.WriteJSON(w, http.StatusNotImplemented, respn)
		return
	}

	file, header, err := r.FormFile("spi")
	if err != nil {
		respn.Status = "Error : File tidak ada"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusBadRequest, respn)
		return
	}
	defer file.Close()
	// Read the file content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		respn.Status = "Error : File tidak bisa dibaca"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusInternalServerError, respn)
		return
	}
	// Calculate hash of the file content
	//hashedFileName := ghupload.CalculateHash(fileContent)

	// Get GitHub credentials and other details from the request or environment variables
	GitHubAccessToken := config.GHAccessToken
	GitHubAuthorName := "Rolly Maulana Awangga"
	GitHubAuthorEmail := "awangga@gmail.com"
	githubOrg := "penerbitbukupedia"
	githubRepo := "draft"
	pathFile := prj.Name + "/spi/" + prj.ID.Hex() + header.Filename[strings.LastIndex(header.Filename, "."):] // Append the original file extension
	replace := true

	// Use GithubUpload function to upload the file to GitHub
	content, _, err := ghupload.GithubUpload(GitHubAccessToken, GitHubAuthorName, GitHubAuthorEmail, fileContent, githubOrg, githubRepo, pathFile, replace)
	if err != nil {
		respn.Status = "Error : File tidak bisa diupload ke github"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusInternalServerError, respn)
		return
	}
	//update data profpic
	prj.SPI = *content.Content.Path
	atdb.ReplaceOneDoc(config.Mongoconn, "project", bson.M{"_id": prj.ID}, prj)

	// Respond with success message
	respn.Info = prj.ID.Hex()
	respn.Location = prj.SPI
	respn.Response = *content.Content.URL
	respn.Status = *content.Content.HTMLURL
	at.WriteJSON(w, http.StatusOK, respn)
}

func UploadSampulBukuPDFWithParamFileHandler(w http.ResponseWriter, r *http.Request) {
	var respn model.Response
	payload, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(r))
	if err != nil {
		respn.Status = "Error : Token Tidak Valid"
		respn.Info = at.GetSecretFromHeader(r)
		respn.Location = "Decode Token Error"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusForbidden, respn)
		return
	}
	docuser, err := atdb.GetOneDoc[model.Userdomyikado](config.Mongoconn, "user", primitive.M{"phonenumber": payload.Id})
	if err != nil {
		respn.Status = "Error : Data user tidak di temukan"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusNotImplemented, respn)
		return
	}
	prjid := at.GetParam(r)
	objectId, _ := primitive.ObjectIDFromHex(prjid)
	prj, err := atdb.GetOneDoc[model.Project](config.Mongoconn, "project", primitive.M{"_id": objectId})
	if err != nil {
		respn.Status = "Error : Data lapak tidak di temukan"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusNotImplemented, respn)
		return
	}
	//check apakah dia owner
	if prj.Owner.PhoneNumber != docuser.PhoneNumber {
		respn.Status = "Error : User bukan owner project tidak berhak"
		respn.Response = "User bukan owner dari project ini"
		at.WriteJSON(w, http.StatusNotImplemented, respn)
		return
	}

	file, header, err := r.FormFile("sampulpdfbuku")
	if err != nil {
		respn.Status = "Error : File tidak ada"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusBadRequest, respn)
		return
	}
	defer file.Close()
	// Read the file content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		respn.Status = "Error : File tidak bisa dibaca"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusInternalServerError, respn)
		return
	}
	// Calculate hash of the file content
	//hashedFileName := ghupload.CalculateHash(fileContent)

	// Get GitHub credentials and other details from the request or environment variables
	GitHubAccessToken := config.GHAccessToken
	GitHubAuthorName := "Rolly Maulana Awangga"
	GitHubAuthorEmail := "awangga@gmail.com"
	githubOrg := "penerbitbukupedia"
	githubRepo := "draft"
	pathFile := prj.Name + "/sampul/" + prj.ID.Hex() + header.Filename[strings.LastIndex(header.Filename, "."):] // Append the original file extension
	replace := true

	// Use GithubUpload function to upload the file to GitHub
	content, _, err := ghupload.GithubUpload(GitHubAccessToken, GitHubAuthorName, GitHubAuthorEmail, fileContent, githubOrg, githubRepo, pathFile, replace)
	if err != nil {
		respn.Status = "Error : File tidak bisa diupload ke github"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusInternalServerError, respn)
		return
	}
	//update data profpic
	prj.SampulPDFBuku = *content.Content.Path
	atdb.ReplaceOneDoc(config.Mongoconn, "project", bson.M{"_id": prj.ID}, prj)

	// Respond with success message
	respn.Info = prj.ID.Hex()
	respn.Location = prj.SampulPDFBuku
	respn.Response = *content.Content.URL
	respn.Status = *content.Content.HTMLURL
	at.WriteJSON(w, http.StatusOK, respn)
}

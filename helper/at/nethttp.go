package at

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func GetSecretFromHeader(r *http.Request) (secret string) {
	if r.Header.Get("secret") != "" {
		secret = r.Header.Get("secret")
	} else if r.Header.Get("Secret") != "" {
		secret = r.Header.Get("Secret")
	}
	return
}

func GetLoginFromHeader(r *http.Request) (secret string) {
	if r.Header.Get("login") != "" {
		secret = r.Header.Get("login")
	} else if r.Header.Get("Login") != "" {
		secret = r.Header.Get("Login")
	}
	return
}

func Jsonstr(strc interface{}) string {
	jsonData, err := json.Marshal(strc)
	if err != nil {
		log.Fatal(err)
	}
	return string(jsonData)
}

func WriteJSON(respw http.ResponseWriter, statusCode int, content interface{}) {
	respw.Header().Set("Content-Type", "application/json")
	respw.WriteHeader(statusCode)
	respw.Write([]byte(Jsonstr(content)))
}

func WriteFile(w http.ResponseWriter, statusCode int, fileContent []byte) {
	w.Header().Set("Content-Disposition", "attachment; filename=\"file.ext\"")
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprint(len(fileContent)))
	w.WriteHeader(statusCode)
	w.Write(fileContent)
}

func WriteString(respw http.ResponseWriter, statusCode int, content string) {
	respw.WriteHeader(statusCode)
	respw.Write([]byte(content))
}

package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"
	"time"
	"strings"
	"unicode"

	"github.com/gocroot/config"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/helper/auth"
	"github.com/gocroot/helper/watoken"
	"github.com/gocroot/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func RegisterGmailAuth(w http.ResponseWriter, r *http.Request) {
	logintoken, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(r))
	if err != nil {
		var respn model.Response
		respn.Status = "Error : Token Tidak Valid "
		respn.Info = at.GetSecretFromHeader(r)
		respn.Location = "Decode Token Error: " + at.GetLoginFromHeader(r)
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusForbidden, respn)
		return
	}
	var request struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid request"})
		return
	}

	// Ambil kredensial dari database
	creds, err := atdb.GetOneDoc[auth.GoogleCredential](config.Mongoconn, "credentials", bson.M{})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(map[string]string{"message": "Database Connection Problem: Unable to fetch credentials"})
		return
	}

	// Verifikasi ID token menggunakan client_id
	payload, err := auth.VerifyIDToken(request.Token, creds.ClientID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid token: Token verification failed"})
		return
	}

	userInfo := model.Userdomyikado{
		Name:                 payload.Claims["name"].(string),
		PhoneNumber:          logintoken.Id,
		Email:                payload.Claims["email"].(string),
		GoogleProfilePicture: payload.Claims["picture"].(string),
		Role:                 "user",
	}

	// Simpan atau perbarui informasi pengguna di database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := config.Mongoconn.Collection("user")
	filter := bson.M{"phonenumber": logintoken.Id}

	var existingUser model.Userdomyikado
	err = collection.FindOne(ctx, filter).Decode(&existingUser)
	if err != nil || existingUser.PhoneNumber == "" {
		// User does not exist or exists but has no phone number, insert into db
		id, err := atdb.InsertOneDoc(config.Mongoconn, "user", userInfo)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadGateway)
			json.NewEncoder(w).Encode(map[string]string{"message": "Database Connection Problem: Unable to fetch credentials"})
			return
		}
		response := map[string]interface{}{
			"message": "User Berhasil Terdaftar",
			"user":    userInfo,
			"id":      id.Hex(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	} else if existingUser.PhoneNumber != "" {
		existingUser.Email = userInfo.Email
		existingUser.GoogleProfilePicture = userInfo.GoogleProfilePicture
		_, err := atdb.ReplaceOneDoc(config.Mongoconn, "user", bson.M{"_id": existingUser.ID}, existingUser)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadGateway)
			json.NewEncoder(w).Encode(map[string]string{"message": "Database Connection Problem: Unable to update user"})
			return
		}
		response := map[string]interface{}{
			"message": "Authenticated successfully",
			"user":    existingUser,
			"id":      existingUser.ID,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	update := bson.M{
		"$set": userInfo,
	}
	opts := options.Update().SetUpsert(true)
	_, err = collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to save user info: Database update failed"})
		return
	}

	response := map[string]interface{}{
		"user": userInfo,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// register manual
// RegisterUser handles user registration with only essential fields
// RegisterUser handles user registration with only essential fields
func RegisterUser(respw http.ResponseWriter, req *http.Request) {
    // Decode the incoming request body into the Userdomyikado struct
    var usr model.Userdomyikado
    err := json.NewDecoder(req.Body).Decode(&usr)
    if err != nil {
        var respn model.Response
        respn.Status = "Error: Body tidak valid"
        respn.Response = err.Error()
        at.WriteJSON(respw, http.StatusBadRequest, respn)
        return
    }

    // Validate required fields
    if usr.Name == "" || usr.PhoneNumber == "" || usr.Email == "" || usr.Password == "" {
        var respn model.Response
        respn.Status = "Isian tidak lengkap"
        respn.Response = "Mohon isi lengkap Nama, WhatsApp, Email, dan Password"
        at.WriteJSON(respw, http.StatusBadRequest, respn)
        return
    }

    // Normalize phone number (convert +62, 62, or 08 to 62)
    usr.PhoneNumber = normalizePhoneNumber(usr.PhoneNumber)

    // **Check for unique email and phone number**
    collection := config.Mongoconn.Collection("user")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Cek apakah email atau nomor telepon sudah digunakan
    filter := bson.M{
        "$or": []bson.M{
            {"email": usr.Email},
            {"phonenumber": usr.PhoneNumber},
        },
    }
    var existingUser model.Userdomyikado
    err = collection.FindOne(ctx, filter).Decode(&existingUser)
    if err == nil {
        // Jika ditemukan, berarti email atau nomor telepon sudah digunakan
        var respn model.Response
        respn.Status = "Email atau Nomor WhatsApp sudah terdaftar"
        respn.Response = "Silakan gunakan email atau nomor WhatsApp yang berbeda"
        at.WriteJSON(respw, http.StatusConflict, respn)
        return
    }

    // Hash the password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(usr.Password), bcrypt.DefaultCost)
    if err != nil {
        var respn model.Response
        respn.Status = "Error: Gagal meng-hash password"
        respn.Response = err.Error()
        at.WriteJSON(respw, http.StatusInternalServerError, respn)
        return
    }
    usr.Password = string(hashedPassword)

    // Set default role to "user"
    usr.Role = "user"

    // Generate a new ID for the user
    usr.ID = primitive.NewObjectID()

    // Insert the user into the database
    _, err = atdb.InsertOneDoc(config.Mongoconn, "user", usr)
    if err != nil {
        var respn model.Response
        respn.Status = "Error: Gagal insert database"
        respn.Response = err.Error()
        at.WriteJSON(respw, http.StatusNotModified, respn)
        return
    }

    // Clear the password before sending the response
    usr.Password = "" // Jangan kirim password

    // Kirim data user langsung tanpa wrapping "user"
    at.WriteJSON(respw, http.StatusOK, usr)
}

// Function to normalize phone number to start with 62
func normalizePhoneNumber(phone string) string {
	// Remove all non-numeric characters
	phone = removeNonNumeric(phone)

	// Remove leading '+' if exists
	phone = strings.TrimPrefix(phone, "+")

	phone = strings.TrimPrefix(phone, "+")
	if strings.HasPrefix(phone, "08") {
		return "62" + phone[2:]
	} else if phone[0] == '0' {
		return "62" + phone[1:]
	}
	// Jika sudah '62', tidak perlu perubahan
	return phone

}

// Helper function to remove all non-numeric characters
func removeNonNumeric(s string) string {
	var result []rune
	for _, r := range s {
		if unicode.IsDigit(r) {
			result = append(result, r)
		}
	}
	return string(result)
}

func Auth(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid request"})
		return
	}

	// Ambil kredensial dari database
	creds, err := atdb.GetOneDoc[auth.GoogleCredential](config.Mongoconn, "credentials", bson.M{})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(map[string]string{"message": "Database Connection Problem: Unable to fetch credentials"})
		return
	}

	// Verifikasi ID token menggunakan client_id
	payload, err := auth.VerifyIDToken(request.Token, creds.ClientID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid token: Token verification failed"})
		return
	}

	userInfo := model.Userdomyikado{
		Name:                 payload.Claims["name"].(string),
		Email:                payload.Claims["email"].(string),
		GoogleProfilePicture: payload.Claims["picture"].(string),
	}

	// Simpan atau perbarui informasi pengguna di database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := config.Mongoconn.Collection("user")
	filter := bson.M{"email": userInfo.Email}

	var existingUser model.Userdomyikado
	err = collection.FindOne(ctx, filter).Decode(&existingUser)
	if err != nil || existingUser.PhoneNumber == "" {
		// User does not exist or exists but has no phone number, request QR scan
		response := map[string]interface{}{
			"message": "Please scan the QR code to provide your phone number",
			"user":    userInfo,
			"token":   "",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	} else if existingUser.PhoneNumber != "" {
		token, err := watoken.EncodeforHours(existingUser.PhoneNumber, existingUser.Name, config.PrivateKey, 18) // Generating a token for 18 hours
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Token generation failed"})
			return
		}
		response := map[string]interface{}{
			"message": "Authenticated successfully",
			"user":    userInfo,
			"token":   token,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	update := bson.M{
		"$set": userInfo,
	}
	opts := options.Update().SetUpsert(true)
	_, err = collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to save user info: Database update failed"})
		return
	}

	response := map[string]interface{}{
		"user": userInfo,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parse the incoming request body
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid request"})
		return
	}

	// Ensure both email and password are provided
	if request.Email == "" || request.Password == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Email and password are required"})
		return
	}

	// Connect to the MongoDB collection
	collection := config.Mongoconn.Collection("user")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find the user by email
	var user model.Userdomyikado
	err := collection.FindOne(ctx, bson.M{"email": request.Email}).Decode(&user)
	if err != nil {
		// Handle case where user is not found
		if err == mongo.ErrNoDocuments {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"message": "Invalid email or user not registered"})
		} else {
			// For other errors (e.g., database issues), respond with 500
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Database query failed", "error": err.Error()})
		}
		return
	}

	// Validate the password using bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid password"})
		return
	}

	// Generate the authentication token
	token, err := watoken.EncodeforHours(user.PhoneNumber, user.Name, config.PrivateKey, 18) // Token valid for 18 hours
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Token generation failed", "error": err.Error()})
		return
	}

	// Return response with user info and token
	response := map[string]interface{}{
		"message": "Login successful",
		"user": map[string]interface{}{
			"name":  user.Name,
			"email": user.Email,
		},
		"token": token,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}


func GeneratePasswordHandler(respw http.ResponseWriter, r *http.Request) {
	var request struct {
		PhoneNumber string `json:"phonenumber"`
		Captcha     string `json:"captcha"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		var respn model.Response
		respn.Status = "Invalid Request"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}
	// Validate CAPTCHA
	captchaResponse, err := http.PostForm("https://challenges.cloudflare.com/turnstile/v0/siteverify", url.Values{
		"secret":   {"0x4AAAAAAAfj2NjfaHRBhkd2VjcfmRe5gvI"},
		"response": {request.Captcha},
	})
	if err != nil {
		var respn model.Response
		respn.Status = "Failed to verify captcha"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusServiceUnavailable, respn)
		return
	}
	defer captchaResponse.Body.Close()

	var captchaResult struct {
		Success bool `json:"success"`
	}
	if err := json.NewDecoder(captchaResponse.Body).Decode(&captchaResult); err != nil {
		var respn model.Response
		respn.Status = "Failed to decode captcha response"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusInternalServerError, respn)
		return
	}
	if !captchaResult.Success {
		var respn model.Response
		respn.Status = "Unauthorized"
		respn.Response = "Invalid captcha"
		at.WriteJSON(respw, http.StatusUnauthorized, respn)
		return
	}

	// Validate phone number
	re := regexp.MustCompile(`^62\d{9,15}$`)
	if !re.MatchString(request.PhoneNumber) {
		var respn model.Response
		respn.Status = "Bad Request"
		respn.Response = "Invalid phone number format"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	// Check if phone number exists in the 'user' collection
	userFilter := bson.M{"phonenumber": request.PhoneNumber}
	_, err = atdb.GetOneDoc[model.Userdomyikado](config.Mongoconn, "user", userFilter)
	if err != nil {
		var respn model.Response
		respn.Status = "Unauthorized"
		respn.Response = "Phone number not registered"
		at.WriteJSON(respw, http.StatusUnauthorized, respn)
		return
	}

	// Generate random password
	randomPassword, err := auth.GenerateRandomPassword(12)
	if err != nil {
		var respn model.Response
		respn.Status = "Failed to generate password"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusInternalServerError, respn)
		return
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(randomPassword)
	if err != nil {
		var respn model.Response
		respn.Status = "Failed to hash password"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusInternalServerError, respn)
		return
	}

	// Update or insert the user in the database
	stpFilter := bson.M{"phonenumber": request.PhoneNumber}
	_, err = atdb.GetOneDoc[model.Stp](config.Mongoconn, "stp", stpFilter)
	var responseMessage string

	if err == mongo.ErrNoDocuments {
		// Document not found, insert new one
		newUser := model.Stp{
			PhoneNumber:  request.PhoneNumber,
			PasswordHash: hashedPassword,
			CreatedAt:    time.Now(),
			Role:         "user", // Set the role for new user
		}
		_, err = atdb.InsertOneDoc(config.Mongoconn, "stp", newUser)
		if err != nil {
			var respn model.Response
			respn.Status = "Failed to insert new user"
			respn.Response = err.Error()
			at.WriteJSON(respw, http.StatusNotModified, respn)
			return
		}
		responseMessage = "New user created and password generated successfully"
	} else {
		// Document found, update the existing one
		stpUpdate := bson.M{
			"phonenumber": request.PhoneNumber,
			"password":    hashedPassword,
			"createdAt":   time.Now(),
			"role":        "user", // Ensure the role is set for existing user
		}
		_, err = atdb.UpdateOneDoc(config.Mongoconn, "stp", stpFilter, stpUpdate)
		if err != nil {
			var respn model.Response
			respn.Status = "Failed to update user"
			respn.Response = err.Error()
			at.WriteJSON(respw, http.StatusInternalServerError, respn)
			return
		}
		responseMessage = "User info updated and password generated successfully"
	}

	// Respond with success and the generated password
	response := map[string]interface{}{
		"message":     responseMessage,
		"phonenumber": request.PhoneNumber,
	}
	at.WriteJSON(respw, http.StatusOK, response)

	// Send the random password via WhatsApp
	auth.SendWhatsAppPassword(respw, request.PhoneNumber, randomPassword)
}

var (
	rl = auth.NewRateLimiter(1, 5) // 1 request per second, burst of 5
)

func VerifyPasswordHandler(respw http.ResponseWriter, r *http.Request) {
	var request struct {
		PhoneNumber string `json:"phonenumber"`
		Password    string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		var respn model.Response
		respn.Status = "Invalid Request"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	// Implementasi rate limiting
	limiter := rl.GetLimiter(request.PhoneNumber)
	if !limiter.Allow() {
		var respn model.Response
		respn.Status = "Too Many Requests"
		respn.Response = "Please try again later."
		at.WriteJSON(respw, http.StatusTooManyRequests, respn)
		return
	}

	// Find user in the database
	userFilter := bson.M{"phonenumber": request.PhoneNumber}
	user, err := atdb.GetOneDoc[model.Stp](config.Mongoconn, "stp", userFilter)
	if err != nil {
		var respn model.Response
		respn.Status = "Failed to verify password"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusUnauthorized, respn)
		return
	}

	// Verify password and expiry
	if time.Now().After(user.CreatedAt.Add(4 * time.Minute)) {
		var respn model.Response
		respn.Status = "Unauthorized"
		respn.Response = "Password Expired"
		at.WriteJSON(respw, http.StatusUnauthorized, respn)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password))
	if err != nil {
		var respn model.Response
		respn.Status = "Failed to verify password"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusUnauthorized, respn)
		return
	}

	// Find user in the 'user' collection
	myiUserFilter := bson.M{"phonenumber": request.PhoneNumber}
	existingUser, err := atdb.GetOneDoc[model.Userdomyikado](config.Mongoconn, "user", myiUserFilter)
	if err != nil {
		var respn model.Response
		respn.Status = "Unauthorized"
		respn.Response = "Phone number not registered"
		at.WriteJSON(respw, http.StatusUnauthorized, respn)
		return
	}

	token, err := watoken.EncodeforHours(existingUser.PhoneNumber, existingUser.Name, config.PrivateKey, 18)
	if err != nil {
		var respn model.Response
		respn.Status = "Failed to give the token"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusInternalServerError, respn)
		return
	}

	response := map[string]interface{}{
		"message": "Authenticated successfully",
		"token":   token,
		"name":    existingUser.Name,
	}

	// Respond with success
	at.WriteJSON(respw, http.StatusOK, response)
}

func ResendPasswordHandler(respw http.ResponseWriter, r *http.Request) {
	var request struct {
		PhoneNumber string `json:"phonenumber"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		var respn model.Response
		respn.Status = "Invalid Request"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	// Generate random password
	randomPassword, err := auth.GenerateRandomPassword(12)
	if err != nil {
		var respn model.Response
		respn.Status = "Failed to generate password"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusInternalServerError, respn)
		return
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(randomPassword)
	if err != nil {
		var respn model.Response
		respn.Status = "Failed to hash password"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusInternalServerError, respn)
		return
	}

	// Check if phone number exists in the 'stp' collection
	stpFilter := bson.M{"phonenumber": request.PhoneNumber}
	_, stpErr := atdb.GetOneDoc[model.Stp](config.Mongoconn, "stp", stpFilter)

	if stpErr == mongo.ErrNoDocuments {
		// Document not found, insert new one
		newUser := model.Stp{
			PhoneNumber:  request.PhoneNumber,
			PasswordHash: hashedPassword,
			CreatedAt:    time.Now(),
		}
		_, err = atdb.InsertOneDoc(config.Mongoconn, "stp", newUser)
		if err != nil {
			var respn model.Response
			respn.Status = "Failed to insert new user"
			respn.Response = err.Error()
			at.WriteJSON(respw, http.StatusInternalServerError, respn)
			return
		}
		responseMessage := "New user created and password generated successfully"

		// Respond with success and the generated password
		response := map[string]interface{}{
			"message":     responseMessage,
			"phonenumber": request.PhoneNumber,
		}
		at.WriteJSON(respw, http.StatusOK, response)

		// Send the random password via WhatsApp
		auth.SendWhatsAppPassword(respw, request.PhoneNumber, randomPassword)
		return
	} else if stpErr != nil {
		var respn model.Response
		respn.Status = "Failed to fetch user info"
		respn.Response = stpErr.Error()
		at.WriteJSON(respw, http.StatusInternalServerError, respn)
		return
	}

	// Document found, update the existing one
	stpUpdate := bson.M{
		"phonenumber": request.PhoneNumber,
		"password":    hashedPassword,
		"createdAt":   time.Now(),
	}
	_, err = atdb.UpdateOneDoc(config.Mongoconn, "stp", stpFilter, stpUpdate)
	if err != nil {
		var respn model.Response
		respn.Status = "Failed to update user"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusInternalServerError, respn)
		return
	}
	responseMessage := "User info updated and password generated successfully"

	// Respond with success and the generated password
	response := map[string]interface{}{
		"message":     responseMessage,
		"phonenumber": request.PhoneNumber,
	}
	at.WriteJSON(respw, http.StatusOK, response)

	// Send the random password via WhatsApp
	auth.SendWhatsAppPassword(respw, request.PhoneNumber, randomPassword)
}

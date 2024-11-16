package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"sync"

	"github.com/gocroot/config"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/helper/atapi"
	"github.com/gocroot/helper/whatsauth"
	"github.com/gocroot/model"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/time/rate"
	"google.golang.org/api/idtoken"
)

type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.Mutex
	r        rate.Limit
	burst    int
}

func NewRateLimiter(r rate.Limit, burst int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		r:        r,
		burst:    burst,
	}
}

func (rl *RateLimiter) GetLimiter(key string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[key]
	if !exists {
		limiter = rate.NewLimiter(rl.r, rl.burst)
		rl.limiters[key] = limiter
	}

	return limiter
}

func VerifyIDToken(idToken string, audience string) (*idtoken.Payload, error) {
	payload, err := idtoken.Validate(context.Background(), idToken, audience)
	if err != nil {
		return nil, fmt.Errorf("id token validation failed: %v", err)
	}
	return payload, nil
}

func GenerateRandomPassword(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func SendWhatsAppPassword(respw http.ResponseWriter, phoneNumber string, password string) {
    // Prepare WhatsApp message
    dt := &whatsauth.TextMessage{
        To:      phoneNumber,
        IsGroup: false,
        Messages: "Hi! Your login password is: *" + password + "*.\n\n" +
        "Enter this password on the STP page within 4 minutes. The password will expire after that. " +
        "To copy the password, press and hold the password.",
    }

    // Send WhatsApp message
    _, resp, err := atapi.PostStructWithToken[model.Response]("Token", config.WAAPIToken, dt, config.WAAPIMessage)
    if err != nil {
		resp.Info = "message: unauthorized"
		resp.Response = err.Error()
		at.WriteJSON(respw, http.StatusUnauthorized, resp)
		return
    }

}





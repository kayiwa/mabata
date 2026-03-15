package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/kayiwa/mabata/internal/models"
)

type sessionPayload struct {
	User      models.User `json:"user"`
	ExpiresAt int64       `json:"expires_at"`
}

type SessionManager struct {
	secret []byte
}

func NewSessionManager(secret string) *SessionManager {
	return &SessionManager{secret: []byte(secret)}
}

func (s *SessionManager) Set(w http.ResponseWriter, user models.User) error {
	payload := sessionPayload{
		User:      user,
		ExpiresAt: time.Now().Add(8 * time.Hour).Unix(),
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	data := base64.RawURLEncoding.EncodeToString(b)
	sig := s.sign(data)
	cookie := &http.Cookie{
		Name:     "mabata_session",
		Value:    data + "." + sig,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
	return nil
}

func (s *SessionManager) Clear(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "mabata_session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
}

func (s *SessionManager) Get(r *http.Request) (*models.User, error) {
	c, err := r.Cookie("mabata_session")
	if err != nil {
		return nil, err
	}
	parts := strings.Split(c.Value, ".")
	if len(parts) != 2 {
		return nil, errors.New("invalid session format")
	}
	if !hmac.Equal([]byte(parts[1]), []byte(s.sign(parts[0]))) {
		return nil, errors.New("invalid session signature")
	}
	b, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, err
	}
	var payload sessionPayload
	if err := json.Unmarshal(b, &payload); err != nil {
		return nil, err
	}
	if time.Now().Unix() > payload.ExpiresAt {
		return nil, errors.New("session expired")
	}
	return &payload.User, nil
}

func (s *SessionManager) sign(data string) string {
	h := hmac.New(sha256.New, s.secret)
	_, _ = h.Write([]byte(data))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

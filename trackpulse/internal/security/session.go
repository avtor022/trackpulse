package security

import (
	"crypto/rand"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Session struct {
	UserID    string
	Token     string
	CreatedAt time.Time
}

func NewSession(userID string) *Session {
	token := generateToken()
	return &Session{
		UserID:    userID,
		Token:     token,
		CreatedAt: time.Now(),
	}
}

func (s *Session) IsValid() bool {
	// Сессия действительна в течение 24 часов
	expirationTime := s.CreatedAt.Add(24 * time.Hour)
	return time.Now().Before(expirationTime)
}

func (s *Session) Logout() {
	s.Token = ""
}

func generateToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func ComparePassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
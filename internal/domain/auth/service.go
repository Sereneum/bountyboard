package auth

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"sync"
	"time"
)

type service struct {
	secret []byte
	users  map[string]string
	mu     sync.RWMutex
}

func New(secret string) Service {
	return &service{
		secret: []byte(secret),
		users:  make(map[string]string),
	}
}

func (s *service) Register(ctx context.Context, username, password string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[username]; exists {
		return errors.New("user already exists")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	s.users[username] = string(hashed)
	return nil
}
func (s *service) Login(ctx context.Context, username, password string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	hashed, ok := s.users[username]
	if !ok {
		return "", errors.New("user not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	claims := jwt.MapClaims{
		"user": username,
		"exp":  time.Now().Add(24 * 7 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *service) ValidateToken(ctx context.Context, tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid claims")
	}
	userID, ok := claims["user"].(string)
	if !ok {
		return "", errors.New("user not in token")
	}
	return userID, nil
}

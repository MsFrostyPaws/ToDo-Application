package service

import (
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/MsFrostyPaws/todo-list"
	"github.com/MsFrostyPaws/todo-list/pkg/repository"
	"github.com/dgrijalva/jwt-go"
)

const (
	salt       = "dcgye7gfcbOWU8yfdugtsujcahbia"
	tokenTTL   = 12 * time.Hour
	signingKey = "fsvcGNTw65dv743tdHDN734r833uh"
)

type tokenClains struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user todo.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)

	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClains{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})

	return token.SignedString([]byte(signingKey))
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

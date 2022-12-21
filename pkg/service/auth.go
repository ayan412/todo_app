package service

import (
	"crypto/sha1"
	"fmt"
	"github.com/ayan412/zhashkevych_rest_api/todo-app"
	"github.com/ayan412/zhashkevych_rest_api/todo-app/pkg/repository"
	"github.com/dgrijalva/jwt-go"
)

const (
	
	salt = "img223sds74sdi"
	signingKey = "dfad$sdf874giHdaDEtd43IRas"
	tokenTTL = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandartClaims
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
	// get user from DB
	user, err := s.repo.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigninMethodHS256, &tokenClaims{
		jwt.StandartClaims{
		ExpiresAt: time.Now().Add(tokenTTL).Unix(), 
		IssuedAt: time.Now().Unix(),
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


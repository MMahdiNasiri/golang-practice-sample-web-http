package authenticate

import (
	"context"
	"os"
	"sample-web-http/internal/user"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	secretKey []byte
}

func NewService() *TokenService {
	return &TokenService{
		secretKey: []byte(os.Getenv("SECRET_KEY")),
	}
}

func (s *TokenService) GenerateToken(ctx context.Context, user *user.User) (string, error) {
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)

}

func (s *TokenService) ValidateToken(ctx context.Context, tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return s.secretKey, nil
	})
	if err != nil || !token.Valid {
		return 0, err
	}
	claims := token.Claims.(jwt.MapClaims)
	return claims["sub"].(int), nil
}

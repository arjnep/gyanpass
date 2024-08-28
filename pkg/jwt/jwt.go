package jwt

import (
	"errors"
	"log"
	"time"

	"github.com/arjnep/gyanpass/config"
	"github.com/arjnep/gyanpass/internal/entity"
	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	User *entity.User `json:"user"`
	jwt.RegisteredClaims
}

type Service interface {
	GenerateToken(u *entity.User) (string, error)
	ValidateToken(token string) (*TokenClaims, error)
}

type jwtService struct {
	secretKey string
	issuer    string
	cfg       *config.Configuration
}

func NewJWTService(cfg *config.Configuration) Service {
	return &jwtService{
		secretKey: cfg.Server.JWTSecret,
		issuer:    "gyanpass",
		cfg:       config.GetConfig(),
	}
}

func (s *jwtService) GenerateToken(u *entity.User) (string, error) {

	currentTime := time.Now()
	tokenExp := currentTime.Add(time.Duration(s.cfg.Server.JWTExpiry) * time.Second)

	claims := TokenClaims{
		User: u,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(tokenExp),
			Issuer:    s.issuer,
			IssuedAt:  jwt.NewNumericDate(currentTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		log.Println("Failed to sign id token string")
		return "", err
	}

	return signedToken, nil
}

func (s *jwtService) ValidateToken(tokenString string) (*TokenClaims, error) {
	claims := &TokenClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*TokenClaims)

	if !ok {
		return nil, errors.New("valid token but couldn't parse claims")
	}

	return claims, nil

}

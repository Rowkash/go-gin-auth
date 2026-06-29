package auth

import (
	"errors"
	"time"

	"github.com/Rowkash/go-gin-auth/internal/users"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

//var jwtKey = []byte("billy_herrington")
//var jwtExpiresIn = time.Hour * 24

//	type UserSession struct {
//		UserId   uint   `json:"id"`
//		UserRole string `json:"role"`
//	}

type Claims struct {
	UserId   uint   `json:"userId"`
	UserRole string `json:"userRole"`
	jwt.RegisteredClaims
}

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type TokenService struct {
	secretKey []byte
	expiresIn time.Duration
}

func NewTokenService(secret string, expiresIn time.Duration) *TokenService {
	return &TokenService{
		secretKey: []byte(secret),
		expiresIn: expiresIn,
	}
}

func (s *TokenService) generate(user users.User) (*Tokens, error) {
	//expirationTime := time.Now().Add(s.jwtExpiresIn)
	expirationTime := time.Now().Add(s.expiresIn)

	claims := &Claims{
		UserId:   user.ID,
		UserRole: string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//accessToken, err := token.SignedString(jwtKey)
	accessToken, err := token.SignedString(s.secretKey)
	if err != nil {
		return nil, err
	}
	refreshToken := uuid.New()

	return &Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.String(),
	}, nil
}

func (s *TokenService) validate(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		//return jwtKey, nil
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

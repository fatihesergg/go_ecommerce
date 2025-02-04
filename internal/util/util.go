package util

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var JWTSECRET string

type JwtTokenClaims struct {
	jwt.RegisteredClaims
	Role string `json:"role"`
}

type ApiResponse struct {
	Status  int         `json:"status_code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func WriteJson(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func CreateJWT(userID string, role string) (string, error) {
	claims := JwtTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "go_ecommerce",
			Audience:  []string{"go_ecommerce_api"},
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Role: role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JWTSECRET))
}

func ParseJWT(jwtToken string) (*JwtTokenClaims, error) {
	token, err := jwt.ParseWithClaims(jwtToken, &JwtTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTSECRET), nil
	})

	if claims, ok := token.Claims.(*JwtTokenClaims); ok && token.Valid {
		if claims.Issuer != "go_ecommerce" {
			return nil, fmt.Errorf("invalid issuer")
		}
		return claims, nil
	} else {
		return nil, err
	}
}

func EncryptPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

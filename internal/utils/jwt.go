package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var JwtKey = []byte("randomtest")

type Claims struct {
	UserID string `json:"user_id"`
	RoomID string `json:"room_id"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID string, roomID string) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID: userID,
		RoomID: roomID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "WhenRU3",
			ExpiresAt: &jwt.NumericDate{Time: now.Add(time.Hour * 24)},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtKey)
}

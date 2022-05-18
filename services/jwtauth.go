package services

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateTokenWithMethod(signMethod *jwt.SigningMethodHMAC, claims jwt.MapClaims, secret *string) (*string, error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(signMethod, claims)
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(*secret))
	if err != nil {
		return nil, err
	}
	return &tokenString, nil
}

func GenerateTokenWithClaims(email string) string {
	claims := jwt.MapClaims{
		"exp":        time.Now().Add(time.Hour * 3).Unix(),
		"iat":        time.Now().Unix(),
		"user_email": email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return t

}

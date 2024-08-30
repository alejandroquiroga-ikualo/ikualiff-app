package internal

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	Email string
	jwt.StandardClaims
}

func GenerateJwt(email string) string {
	claims := Claims{
		email,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	apiJwtKey := GetEnv()[API_JWT_KEY]

	tokenString, err := token.SignedString([]byte(apiJwtKey))
	if err != nil {
		log.Fatalf("Error signing token: %v", err)
	}

	return tokenString
}

func ParseJwtToken(token string) *Claims {
	parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		apiJwtKey := GetEnv()[API_JWT_KEY]
		return []byte(apiJwtKey), nil
	})

	if err != nil {
		log.Fatalf("Error parsing JWT token: %v", err)
	} else if claims, ok := parsedToken.Claims.(*Claims); ok {
		return claims
	}

	log.Fatalf("Unkwnown claims type")
	return nil
}

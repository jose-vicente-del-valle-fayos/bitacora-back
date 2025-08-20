package utilidades

import (
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

// secreto stores de JWT secret
var secreto = []byte(os.Getenv("SECRET_JWT"))

// GenerarJWT generates a JSON Web Token from user's ID (string)
func GenerarJWT(issuer string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    issuer,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 12)), // Expira en medio día
	})
	return claims.SignedString(secreto)
}

// ParsearJWT gets a JWT and returns a user
func ParsearJWT(cookie string) (string, error) {
	token, err := jwt.ParseWithClaims(cookie, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secreto), nil
	})
	if err != nil || !token.Valid {
		return "", err
	}
	claims := token.Claims.(*jwt.RegisteredClaims)
	return claims.GetIssuer()
}

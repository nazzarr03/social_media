package middleware

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/nazzarr03/social-media/models"
)

var secret = []byte(os.Getenv("JWT_SECRET"))

type JwtCustomClaims struct {
	Username string `json:"username"`
	ID       uint   `json:"id"`
	jwt.StandardClaims
}

func GenerateToken(user *models.User, expiry int) (string, error) {
	exp := time.Now().Add(time.Hour * time.Duration(expiry)).Unix()
	claims := &JwtCustomClaims{
		Username: user.Username,
		ID:       uint(user.UserID),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessToken, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func ValidateToken(tokenString string) (*JwtCustomClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JwtCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JwtCustomClaims)
	if !ok {
		return nil, errors.New("invalid token")
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("token has expired")
	}

	return claims, nil
}

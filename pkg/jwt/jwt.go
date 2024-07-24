package jwt_service

import (
	"errors"
	"time"

	"github.com/Frozelo/music-rate-service/internal/domain/entity"
	"github.com/golang-jwt/jwt"
)

var jwtKey = []byte("some secret key")

type Claims struct {
	UserId int    `json:"userId"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

func GenerateJWT(authUser *entity.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserId: authUser.ID,
		Email:  authUser.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New("invalid token")
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return nil, errors.New("token is expired or not active yet")
			} else {
				return nil, errors.New("couldn't handle this token")
			}
		}
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

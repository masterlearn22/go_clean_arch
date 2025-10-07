package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go_clean/app/models"
	"go_clean/config"
)


func GenerateToken(u models.User) (string, error) {
	jwtCfg := config.LoadJWT()
	claims := models.JWTClaims{
		UserID:   u.ID,
		Username: u.Username,
		Role:     u.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(jwtCfg.TTLHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString(jwtCfg.Secret)
}

func ValidateToken(tokenStr string) (*models.JWTClaims, error) {
	jwtCfg := config.LoadJWT()
	tok, err := jwt.ParseWithClaims(tokenStr, &models.JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtCfg.Secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := tok.Claims.(*models.JWTClaims); ok && tok.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}

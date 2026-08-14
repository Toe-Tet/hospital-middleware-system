package helper

import (
	"errors"
	"fmt"
	"time"

	"hospital-middleware-system/src/config"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	StaffID    int    `json:"staff_id"`
	HospitalID int    `json:"hospital_id"`
	Role       string `json:"role"`
	jwt.RegisteredClaims
}

type TokenResponse struct {
	Token     string    `json:"token"`
	ExpiresAt int       `json:"expires_at"`
}

func GenerateToken(staffID, hospitalID int) (TokenResponse, error) {
	secret := []byte(config.AppConfig.JWTSecret)
	expiresAt := time.Now().Add(time.Duration(config.AppConfig.JWTExpiresHrs) * time.Hour)

	claims := JWTClaims{
		StaffID:    staffID,
		HospitalID: hospitalID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("staff:%d", staffID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(secret)
	if err != nil {
		return TokenResponse{}, err
	}
	return TokenResponse{
		Token:     tokenStr,
		ExpiresAt: int(expiresAt.Unix()),
	}, nil
}

func ValidateToken(tokenStr string) (*JWTClaims, error) {
	secret := []byte(config.AppConfig.JWTSecret)

	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token claims")
}

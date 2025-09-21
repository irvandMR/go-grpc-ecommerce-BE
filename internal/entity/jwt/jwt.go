package jwt

import (
	"context"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/irvandMR/go-grpc-ecommerce-BE/internal/utils"
)

type JWTClaim struct {
	jwt.RegisteredClaims
	Fullname string `json:"fullname"`
	Role     string `json:"role"`
	Email    string `json:"email"`
}

type JwtEntityContexKey string

var JwtEntityContexKeyValue JwtEntityContexKey = "JwtEntity"

func GetClaimsFromToken(jwtToken string) (*JWTClaim, error) {
	tClaim, err := jwt.ParseWithClaims(jwtToken, &JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		secretKey := os.Getenv("JWT_SECRET")
		return []byte(secretKey), nil
	})
	
	if !tClaim.Valid {
		return nil, utils.UnauthenticatedResponse()
	}
	if err != nil {
		return nil, utils.UnauthenticatedResponse()
	}
	
	claims, ok := tClaim.Claims.(*JWTClaim)
	if !ok {
		return nil, utils.UnauthenticatedResponse()
	}
	return claims, nil
}

func (jc *JWTClaim) SteToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, JwtEntityContexKeyValue, jc)
}
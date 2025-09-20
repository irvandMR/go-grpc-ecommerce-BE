package entity

import "github.com/golang-jwt/jwt/v5"

type JWTClaim struct {
	jwt.RegisteredClaims
	Fullname string `json:"fullname"`
	Role     string `json:"role"`
	Email    string `json:"email"`
}
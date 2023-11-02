package model

import "github.com/golang-jwt/jwt"

type Claims struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	jwt.StandardClaims
}

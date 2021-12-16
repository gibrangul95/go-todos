package util

import (
	"github.com/gibrangul95/go-todos/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gibrangul95/go-todos/database"
	"github.com/dgrijalva/jwt-go"
	"time"
	"github.com/gibrangul95/go-todos/internal/model"
)

var jwtKey = []byte(config.Config("PRIV_KEY"))

func GenerateTokens(uuid string) (string, string) {
	claim, accessToken := GenerateAccessClaims(uuid)
	refreshToken := GenerateRefreshClaims(claim)

	return accessToken, refreshToken
}

func GenerateAccessClaims(uuid string) (*model.Claims, string) {
	t := time.Now()

	claim := &model.Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer: uuid,
			ExpiresAt: t.Add(60 * time.Minute).Unix(),
			Subject: "access_token",
			IssuedAt: t.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		panic(err)
	}

	return claim, tokenString
}

func GenerateRefreshClaims(cl *model.Claims) string {
	result := database.DB.Where(&model.Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer: cl.Issuer,
		},
	}).Find(&model.Claims{})

	if result.RowsAffected > 3 {
		database.DB.Where(&model.Claims{
			StandardClaims: jwt.StandardClaims{
				Issuer: cl.Issuer,
			},
		}).Delete(&model.Claims{})
	}

	t := time.Now()
	refreshClaim := &model.Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer: cl.Issuer, 
			ExpiresAt: t.Add(30 * 24 * time.Hour).Unix(),
			Subject: "refresh_token",
			IssuedAt: t.Unix(),
		},
	}

	database.DB.Create(&refreshClaim)

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaim)
	refreshTokenString, err := refreshToken.SignedString(jwtKey)
	if err != nil {
		panic(err)
	}

	return refreshTokenString
}

func GetAuthCookies(accessToken, refreshToken string) (*fiber.Cookie, *fiber.Cookie) {
	accessCookie := &fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
	}

	refreshCookie := &fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(10 * 24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
	}

	return accessCookie, refreshCookie
}
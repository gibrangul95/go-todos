package middleware

import (
	"github.com/gibrangul95/go-todos/config"
	"time"
	"github.com/dgrijalva/jwt-go"
	"github.com/gibrangul95/go-todos/internal/model"
	"github.com/gofiber/fiber/v2"
)

var jwtKey = []byte(config.Config("PRIV_KEY"))

func SecureAuth() func(*fiber.Ctx) error {
	return func (c *fiber.Ctx) error {
		accessToken := c.Get("access_token")
		claims := new(model.Claims)

		token, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) { 
			return jwtKey, nil
		})

		if token.Valid {
			if claims.ExpiresAt < time.Now().Unix() {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": true,
					"general": "Token Expired",
				})
			}
		} else if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				c.ClearCookie("access_token", "refresh_token")
				return c.SendStatus(fiber.StatusForbidden)
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return c.SendStatus(fiber.StatusUnauthorized)
			} else {
				c.ClearCookie("access_token", "refresh_token")
				return c.SendStatus(fiber.StatusForbidden)
			}
		}

		c.Locals("id", claims.Issuer)
		return c.Next()
	}
}
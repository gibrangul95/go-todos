package userHandler

import (
	"github.com/google/uuid"
	"time"
	"github.com/gibrangul95/go-todos/config"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"github.com/gibrangul95/go-todos/database"
	"github.com/gibrangul95/go-todos/util"
	"github.com/gibrangul95/go-todos/internal/model"
	"github.com/gofiber/fiber/v2"
)

var jwtKey = []byte(config.Config("PRIV_KEY"))

func CreateUser(c *fiber.Ctx) error {
	u := new(model.User)

	if err := c.BodyParser(u); err != nil {
		return c.JSON(fiber.Map{
			"error": true,
			"input": "Please review your input",
		})
	}

	errors := util.ValidateRegister(u)
	if errors.Err {
		return c.JSON(errors)
	}

	if count := database.DB.Where(&model.User{Email: u.Email}).First(new(model.User)).RowsAffected; count > 0 {
		errors.Err, errors.Email = true, "Email is already registered"
	}

	if count := database.DB.Where(&model.User{Username: u.Username}).First(new(model.User)).RowsAffected; count > 0 {
		errors.Err, errors.Username = true, "Username is already registered"
	}

	if errors.Err {
		return c.JSON(errors)
	}

	password := []byte(u.Password)

	hashedPassword, err := bcrypt.GenerateFromPassword(
		password,
		rand.Intn(bcrypt.MaxCost-bcrypt.MinCost)+bcrypt.MinCost,
	)
	if err != nil {
		panic(err)
	}

	u.ID = uuid.New()
	u.Password = string(hashedPassword)
	
	if err := database.DB.Create(&u).Error; err != nil {
		return c.JSON(fiber.Map{
				"error":   true,
				"general": "Something went wrong, please try again later. 😕",
		})
	}

	accessToken, refreshToken := util.GenerateTokens(u.ID.String())
	accessCookie, refreshCookie := util.GetAuthCookies(accessToken, refreshToken)
	c.Cookie(accessCookie)
	c.Cookie(refreshCookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func LoginUser(c *fiber.Ctx) error {
	type LoginInput struct {
		Identity string `json:"identity"`
		Password string `json:"password"`
	}

	input := new(LoginInput)

	if err := c.BodyParser(input); err != nil {
		return c.JSON(fiber.Map{"error": true, "input": "Please review your input"})
	}

	// check if a user exists
	u := new(model.User)
	if res := database.DB.Where(
		&model.User{Email: input.Identity}).Or(
		&model.User{Username: input.Identity},
	).First(&u); res.RowsAffected <= 0 {
		return c.JSON(fiber.Map{"error": true, "general": "Invalid Credentials."})
	}

	// Comparing the password with the hash
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input.Password)); err != nil {
		return c.JSON(fiber.Map{"error": true, "general": "Invalid Credentials."})
	}

	// setting up the authorization cookies
	accessToken, refreshToken := util.GenerateTokens(u.ID.String())
	accessCookie, refreshCookie := util.GetAuthCookies(accessToken, refreshToken)
	c.Cookie(accessCookie)
	c.Cookie(refreshCookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func GetUserData(c *fiber.Ctx) error {
	id := c.Locals("id")

	u := new(model.User)
	if res := database.DB.Where("id = ?", id).First(&u); res.RowsAffected <= 0 {
			return c.JSON(fiber.Map{"error": true, "general": "Cannot find the User"})
	}

	return c.JSON(u)
}

func GetAccessToken(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")

	refreshClaims := new(model.Claims)
	token, _ := jwt.ParseWithClaims(refreshToken, refreshClaims,
		func(token *jwt.Token) (interface{}, error) {
				return jwtKey, nil
		})

	if res := database.DB.Where(
		"expires_at = ? AND issued_at = ? AND issuer = ?",
		refreshClaims.ExpiresAt, refreshClaims.IssuedAt, refreshClaims.Issuer,
	).First(&model.Claims{}); res.RowsAffected <= 0 {
		c.ClearCookie("access_token", "refresh_token")
		return c.SendStatus(fiber.StatusForbidden)
	}

	if token.Valid {
		if refreshClaims.ExpiresAt < time.Now().Unix() {
			c.ClearCookie("access_token", "refresh_token")
			return c.SendStatus(fiber.StatusForbidden)
		}
	} else {
		c.ClearCookie("access_token", "refresh_token")
		return c.SendStatus(fiber.StatusForbidden)
	}

	_, accessToken := util.GenerateAccessClaims(refreshClaims.Issuer)

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
	})

	return c.JSON(fiber.Map{"access_token": accessToken})
}

func Logout(c *fiber.Ctx) error {
	id := c.Locals("id")

	err := database.DB.Delete(model.Claims{}, "issuer LIKE ?", id).Error
	if err != nil {
		panic(err)
	}
	c.ClearCookie("access_token", "refresh_token")
	return c.SendStatus(fiber.StatusOK)
}
package model

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID uuid.UUID `gorm:"type:uuid`
	Username string `json:"username" gorm:"unique"`
	Email string `json:"email" gorm:"unique"`
	Password string `json:"password"`
}

type UserErrors struct {
	Err bool `json:"error"`
	Email string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	jwt.StandardClaims
	ID uint `gorm:"primaryKey"`
}

type Todo struct {
	gorm.Model
	ID uuid.UUID `gorm:"type:uuid`
	Owner uuid.UUID `gorm:"type:uuid`
	Title string
	AssignedTo string
	Completed bool
	DueDate int
}
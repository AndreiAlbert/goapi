package models

import (
	"gorm.io/gorm"
)

type UserEntity struct {
	gorm.Model
	Username string `gorm:"column: username;unique;not null" json:"username"`
	Email    string `gorm:"column:email; unique; not null" json:"email"`
    Password string `gorm:"column: password; not null" json:"-"`
}

type LoginRequest struct {
    Email string `json:"email"`
    Password string `json:"password"`
}

type LoginResponse struct {
    Email string `json:"email"`
    Jwt string `json:"jwt"`
}

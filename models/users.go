package models

import (
	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	Username string `json:"username" gorm:"unique;not null"`
	Password string `json:"password" gorm:"not null"`
	Email    string `json:"email" gorm:"unique;not null"`
	Token    string `json:"Token"`
}

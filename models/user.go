package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	gorm.Model
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Username       string    `gorm:"unique" json:"username"`
	Password       string    `json:"password"`
	AccountCreated time.Time `json:"account_created"`
	AccountUpdated time.Time `json:"account_updated"`
}

type UpdateUserInput struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

package auth

import (
	"time"
)

type User struct {
	Id           int64
	Email        string
	PasswordHash string
	PasswordSalt string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

package users

import "time"

type User struct {
	ID           string
	Email        string
	PasswordHash string
	Provider     string
	ProviderID   string
	IsVerified   bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

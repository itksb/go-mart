package auth

import "time"

// User - user provided by identity provider system. It is not part of the domain model
type User struct {
	ID           string    `json:"id"`
	Login        string    `json:"login"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

package auth

type User struct {
	ID           string `json:"id"`
	Login        string `json:"login"`
	PasswordHash string `json:"-"`
}
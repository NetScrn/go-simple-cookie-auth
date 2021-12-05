package users

type User struct {
	Id             int    `json:"id"`
	Email          string `json:"email"`
	PasswordDigest string `json:"password_digest"`
	IsConfirmed    bool   `json:"is_confirmed"`
}
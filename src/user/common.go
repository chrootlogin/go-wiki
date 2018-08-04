package user

type User struct {
	Username     string   `json:"-"`
	PasswordHash string   `json:"password-hash"`
	Email        string   `json:"email"`
	IsEnabled    bool     `json:"is-enabled"`
	Permissions  []string `json:"permissions"`
}
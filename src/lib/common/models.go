package common

type User struct {
	Username     string   `json:"-"`
	PasswordHash string   `json:"password-hash"`
	Email        string   `json:"email"`
	IsEnabled    bool     `json:"is-enabled"`
	Permissions  []string `json:"permissions"`
}

type ApiResponse struct {
	Message string `json:"message"`
}

type File struct {
	ContentType string            `json:"content-type"`
	Content     string 			  `json:"content"`
	Metadata	map[string]string `json:"metadata"`
}

type Page struct {
	Content  string			   `json:"content"`
	Metadata map[string]string `json:"metadata"`
}

type Configuration struct {
	Title        string `json:"title"`
	Registration bool   `json:"registration"`
}
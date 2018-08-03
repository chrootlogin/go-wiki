package common

const (
	WrongAPIUsageError = "Invalid api call - parameters did not match to method definition"
)

type ApiResponse struct {
	Message string `json:"message"`
}

type File struct {
	ContentType string            `json:"content-type"`
	Content 	string 			  `json:"content"`
	MetaData	map[string]string `json:"metadata"`
}
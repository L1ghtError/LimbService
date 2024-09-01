package model

// GoogleResponse is the response sent by google
type GoogleResponse struct {
	ID       string `json:"id"`
	UserName string `json:"given_name"`
	Email    string `json:"email"`
	Verified bool   `json:"verified_email"`
	Picture  string `json:"picture"`
	Fullname string `json:"name"`
}

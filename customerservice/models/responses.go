package models

// Response represents a generic HTTP response
type Response struct {
	Message string `json:"message" example:"Success message or error details"`
}

// TokenResponse represents the response for successful login
type TokenResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

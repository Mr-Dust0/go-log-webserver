package types

// Create AuthInput structure to store the json data from login.html into
type AuthInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

package models

// Create User Structre and give metadata on the corresponding json keys and make ID primary and Email and Username are unqiue
type User struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	Email    string `json:"email" gorm:"unique"`
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`
}

package controllers

import (
	"net/http"
	"os"
	"time"
	"webserver/initializers"
	"webserver/middleware"
	"webserver/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func Login(ctx *gin.Context) {

	var authInput models.AuthInput

	// Take Json from POST request from login.html and put the data into authInput
	if err := ctx.ShouldBindJSON(&authInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userFound models.User
	initializers.DB.Where("username=?", authInput.Username).Find(&userFound)

	// Check if the Username matches with an user
	if userFound.ID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	// Hash the Password given in the POST request and compare it to the hash in the database to see if they are the same
	if err := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(authInput.Password)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid password"})
		return
	}

	// Create an jwt token with the userId in the claims and expires in 24 hours so they user doesnt have to keep logining back in.
	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userFound.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign the token so end users cant change the claims and become an differnt user or some kind of attack
	token, err := generateToken.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to generate token"})
	}

	// Set Authorization cookie that can last for an day and allow to be sent over http because securecookie is set to false.
	ctx.SetCookie("Authorization", token, 24*3600, "", "", false, true)
}
func ChangePassword(ctx *gin.Context) {

	userName := ctx.PostForm("username")
	oldPassword := ctx.PostForm("oldpassword")
	newPassword := ctx.PostForm("newpassword")
	// Check to see if the password has actaully changed
	if newPassword == oldPassword {
		ctx.HTML(http.StatusConflict, "resetpassword.html", gin.H{"error_message": "Old password cannot equal new password"})
		return
	}
	var user models.User
	initializers.DB.Where("userName = ?", userName).Find(&user)
	// Check if the user exists
	if user.ID == 0 {
		// Dont want to say the username is wrong because that can give attackers an idea of what usernames to brute force
		ctx.HTML(http.StatusConflict, "resetpassword.html", gin.H{"error_message": "Invalid user details"})
		return
	}
	// Hash the old password and see if matches the one stored in the database
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		ctx.HTML(http.StatusConflict, "resetpassword.html", gin.H{"error_message": "Invalid user details"})
		return
	}
	// Hash new password before storing it in database
	newPasswordHash, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	// Update password to new password that is hashed
	initializers.DB.Where("ID = ?", user.ID).Update("Password", newPasswordHash)
	ctx.Redirect(http.StatusFound, "/") // Equivalent to HTTP 302 redirect, which forces a GET request.
	return

}
func GetResetPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "resetpassword.html", gin.H{"userName": middleware.LoggedInUser})

}
func GetLoginPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login.html", gin.H{})
}

package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"webserver/initializers"
	"webserver/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func Login(ctx *gin.Context) {

	var authInput models.AuthInput

	if err := ctx.ShouldBindJSON(&authInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userFound models.User
	initializers.DB.Where("username=?", authInput.Username).Find(&userFound)

	if userFound.ID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(authInput.Password)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid password"})
		return
	}

	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userFound.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := generateToken.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to generate token"})
	}

	fmt.Println(userFound.Email)
	ctx.SetCookie("Authorization", token, 24*3600, "", "", false, true)
	ctx.JSON(200, gin.H{
		"token": token,
	})
}
func ChangePassword(ctx *gin.Context) {

	userName := ctx.PostForm("username")
	oldPassword := ctx.PostForm("oldpassword")
	newPassword := ctx.PostForm("newpassword")
	if newPassword == oldPassword {
		ctx.HTML(http.StatusConflict, "resetpassword.html", gin.H{"error_message": "Old password cannot equal new password"})
		return
	}
	var user models.User
	initializers.DB.Where("userName = ?", userName).Find(&user)
	if user.ID == 0 {
		// Dont want to say the username is wrong because that can give attackers an idea of what usernames to brute force
		ctx.HTML(http.StatusConflict, "resetpassword.html", gin.H{"error_message": "Invalid user details"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		ctx.HTML(http.StatusConflict, "resetpassword.html", gin.H{"error_message": "Invalid user details"})
		return
	}
	initializers.DB.Where("ID = ?", user.ID).Update("Password", newPassword)
	ctx.HTML(http.StatusOK, "index.html", gin.H{"email": user.Email})
}

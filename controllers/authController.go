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

	email := ctx.PostForm("email")
	if email == "" {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{})
		return
	}
	var user models.User
	initializers.DB.Where("Email = ?", email).Find(&user)
	if user.ID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user not found with that email please try again"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"email": user.Email})
}

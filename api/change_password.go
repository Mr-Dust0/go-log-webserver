package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (api *API) GetChangePasswordPage(ctx *gin.Context) {
	username, _ := ctx.Get("userName")
	ctx.HTML(http.StatusOK, "resetpassword.html", gin.H{"userName": username.(string)})
}

func (api *API) RequestChangePassword(ctx *gin.Context) {
	username := ctx.PostForm("username")
	oldPassword := ctx.PostForm("oldpassword")
	newPassword := ctx.PostForm("newpassword")

	// Check to see if the password has actaully changed
	if newPassword == oldPassword {
		ctx.HTML(http.StatusConflict, "resetpassword.html", gin.H{"error_message": "Old password cannot equal new password"})
		return
	}

	user, err := api.Database.GetUserByUsername(username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	// Hash the Password given in the POST request and compare it to the hash in the database to see if they are the same
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		ctx.HTML(http.StatusConflict, "resetpassword.html", gin.H{"error_message": "Invalid user details"})
		return
	}

	// Hash new password before storing it in database
	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to hash password"})
		return
	}

	// Update password to new password that is hashed
	err = api.Database.UpdatePassword(user, newPasswordHash)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to update password"})
		return
	}

	ctx.Redirect(http.StatusFound, "/") // Equivalent to HTTP 302 redirect, which forces a GET request.
	return
}

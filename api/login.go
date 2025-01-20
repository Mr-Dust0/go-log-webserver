package api

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func (api *API) GetLoginPage(ctx *gin.Context) {
	username := ""
	usernameInterface, ok := ctx.Get("userName")
	if ok {
		username = usernameInterface.(string)
	}

	ctx.HTML(http.StatusOK, "login.html", gin.H{"userName": username})
}

func (api *API) RequestLogin(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")

	user, err := api.Database.GetUserByUsername(username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	// Hash the Password given in the POST request and compare it to the hash in the database to see if they are the same
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid password"})
		return
	}

	// Create an jwt token with the userId in the claims and expires in 24 hours so they user doesnt have to keep logining back in.
	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign the token so end users cant change the claims and become an differnt user or some kind of attack
	token, err := generateToken.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to generate token"})
	}

	// Set Authorization cookie that can last for an day and allow to be sent over http because securecookie is set to false.
	ctx.SetCookie("Authorization", token, 24*3600, "", "", false, true)
	// Statusfound doesnt allow posts i think
	ctx.Redirect(http.StatusFound, "/")
}

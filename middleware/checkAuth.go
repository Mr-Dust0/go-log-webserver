package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"webserver/initializers"
	"webserver/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func DecryptJwt(ctx *gin.Context) jwt.MapClaims {
	// Get data stored in Authorization cookie
	emptyClaims := jwt.MapClaims{}
	tokenString, err := ctx.Cookie("Authorization")

	// Check if there is data in the cookie
	if tokenString == "" {
		// Error message will be displayed at the top of the login page in red
		return emptyClaims
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check sigining method is legit
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret that was used to sign the token
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil || !token.Valid {
		return emptyClaims
	}

	// Get claims stored in the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return emptyClaims
	}
	// See if current time is gretater than expired time which means the token is expired and no longer valid
	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		return emptyClaims
	}
	return claims

}

func CheckAuth(ctx *gin.Context) {
	claims := DecryptJwt(ctx)
	var user models.User
	initializers.DB.Where("ID=?", claims["id"]).Find(&user)

	// Check to see if ID in token maps to an valid user
	if user.ID == 0 {
		ctx.Redirect(302, "/login")
	}

	ctx.Next()

}
func GetUsedLoggedIn(ctx *gin.Context) {
	claims := DecryptJwt(ctx)
	var user models.User
	initializers.DB.Where("ID=?", claims["id"]).Find(&user)

	if user.ID == 0 {
		ctx.Set("userName", "")
		ctx.Next()
	}

	ctx.Set("userName", "Welcome "+user.Username)
	ctx.Next()

}

func GetUser(ctx *gin.Context) {
	claims := DecryptJwt(ctx)
	var user models.User
	var message string
	initializers.DB.Where("ID=?", claims["id"]).Find(&user)

	if user.ID == 0 {
		message = ""

	} else {
		message = "Welcome: " + user.Username
	}
	ctx.HTML(http.StatusOK, "username.html", gin.H{"userName": message})

}

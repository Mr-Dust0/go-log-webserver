package middlewares

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

func CheckAuth(ctx *gin.Context) {

	// Get data stored in Authorization cookie
	tokenString, err := ctx.Cookie("Authorization")

	// Check if there is data in the cookie
	if tokenString == "" {
		// Error message will be displayed at the top of the login page in red
		ctx.HTML(http.StatusOK, "login.html", gin.H{
			"error_message": "Need to be authorized to access that page",
		})
		// Stop processesing futher handlers so the desired page ins't loaded
		ctx.Abort()
		return
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
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Get claims stored in the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		ctx.Abort()
		return
	}
	// See if current time is gretater than expired time which means the token is expired and no longer valid
	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var user models.User
	initializers.DB.Where("ID=?", claims["id"]).Find(&user)

	// Check to see if ID in token maps to an valid user
	if user.ID == 0 {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Alow Next handler to run
	ctx.Next()

}

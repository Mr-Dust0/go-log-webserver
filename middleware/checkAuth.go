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

	authHeader, err := ctx.Cookie("Authorization")

	if authHeader == "" {

		ctx.HTML(http.StatusOK, "login.html", gin.H{
			"error_message": "Need to be authorized to access that page",
		})
		ctx.Abort()
		return
	}
	fmt.Println(authHeader)
	tokenString := authHeader
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil || !token.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		ctx.Abort()
		return
	}

	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var user models.User
	initializers.DB.Where("ID=?", claims["id"]).Find(&user)

	if user.ID == 0 {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx.Next()

}

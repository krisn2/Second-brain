package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

//steps
// function that return a func
// write the Auth middleware in return func where para is context
// get the header use context
// if there is nothing in header return missing token and abort the process
// trim the auth header and store it in tokenStr
// parse the token it valid or not
// if not valid return invalid token
// if valid get the claims
// set the userId in context
// next

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Token"})
			ctx.Abort()
			return
		}
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			log.Fatal("JWT_SECRET environment variable not set")
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ") // don't forget to add the space
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			ctx.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			ctx.Set("userId", claims["userId"])
		}
		ctx.Next()
	}
}

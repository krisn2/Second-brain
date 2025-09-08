package handlers

import (
	"net/http"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/krisn2/second-brain/db"
	"github.com/krisn2/second-brain/models"
	"golang.org/x/crypto/bcrypt"
)

func Isvalidpassword(password string) bool {
	if len(password) < 8 || len(password) > 20 {
		return false
	}
	var hasLower, hasUpper, hasNumber, hasSpecial bool

	for _, ch := range password {
		switch {
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasNumber = true
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSpecial = true
		}
	}
	return hasLower && hasUpper && hasNumber && hasSpecial
}

func Register(c *gin.Context) {

	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid Input"})
		return
	}

	if len(body.Username) < 3 || len(body.Username) > 10 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"Error": "Username must be 3-10 chars"})
		return
	}

	if !Isvalidpassword(body.Password) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Password must be 8-20 chars, with upper, lower, number, special"})
		return
	}

	var existing models.User
	if err := db.DB.Where("username = ?", body.Username).First(&existing).Error; err == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "User already exists"})
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(body.Password), 12)
	user := models.User{Username: body.Username, Password: string(hashed)}

	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sign up successful",
	})
}

func Login(c *gin.Context) {

	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid input"})
		return
	}

	var user models.User
	if err := db.DB.Where("username = ?", body.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Wrong username or password"})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)) != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Wrong username or password"})
		return
	}

	claims := jwt.MapClaims{
		"userId": user.ID.String(),
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, _ := token.SignedString([]byte("MY_SECRET"))

	c.JSON(http.StatusOK, gin.H{"token": signed})
}

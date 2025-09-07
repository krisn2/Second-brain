package handlers

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	models "github.com/krisn2/second-brain"
	"github.com/krisn2/second-brain/db"
	"golang.org/x/crypto/bcrypt"
)

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

	passRegex := `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[\W_]).{8,20}$`

	if !regexp.MustCompile(passRegex).MatchString(body.Password) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Weak Password"})
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
		Username string `gorm:"username"`
		Password string `gorm:"password"`
	}

	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid input"})
		return
	}

	var user models.User
	if err := db.DB.Where("username = ?", body.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Wrong username or password"})
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Username), []byte(body.Password)) != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Wrong username or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "login endpoint",
	})
}

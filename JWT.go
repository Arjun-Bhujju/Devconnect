package controllers

import (
	"devconnect/config"
	"devconnect/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// JWT secret key
var jwtKey = []byte("your-secret-key")

// Hash password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Verify password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Generate JWT token
func GenerateToken(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// Register User
func RegisterUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	hashedPassword, _ := HashPassword(user.Password)
	user.Password = hashedPassword

	_, err := config.DB.NamedExec("INSERT INTO users (username, email, password) VALUES (:username, :email, :password)", user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}

	token, _ := GenerateToken(user)
	c.JSON(http.StatusOK, gin.H{"message": "User registered", "token": token})
}

// Login User
func LoginUser(c *gin.Context) {
	var input models.User
	var dbUser models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := config.DB.Get(&dbUser, "SELECT * FROM users WHERE email=$1", input.Email)
	if err != nil || !CheckPasswordHash(input.Password, dbUser.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, _ := GenerateToken(dbUser)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

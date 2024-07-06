package controllers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"webapp/models"
	"webapp/utils"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input models.User

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if a user with the given username (or email) already exists
		var existingUser models.User
		result := db.Where("username = ?", input.Username).First(&existingUser)
		if result.Error == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "A user account with the given username already exists."})
			return
		} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// An unexpected error occurred while checking for existing user
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing user."})
			return
		}

		hashedPassword, err := utils.HashPassword(input.Password)
		fmt.Print(len(hashedPassword))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		user := models.User{
			ID:             uuid.New(),
			FirstName:      input.FirstName,
			LastName:       input.LastName,
			Username:       input.Username,
			Password:       hashedPassword,
			AccountCreated: time.Now(),
			AccountUpdated: time.Now(),
		}

		db.Create(&user)

		c.JSON(http.StatusCreated, gin.H{
			"id":              user.ID,
			"first_name":      user.FirstName,
			"last_name":       user.LastName,
			"username":        user.Username,
			"account_created": user.AccountCreated.Format(time.RFC3339),
			"account_updated": user.AccountUpdated.Format(time.RFC3339),
		})
	}
}

func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Assuming LoginInput is defined in your models to accept login credentials
		var credentials models.LoginInput
		if err := c.ShouldBindJSON(&credentials); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Assuming the User model has a Username field
		var user models.User
		result := db.Where("username = ?", credentials.Username).First(&user)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect username"})
			return
		} else if result.Error != nil {
			// Log this error for debugging purposes
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		fmt.Print(len(user.Password))

		if !utils.CheckPasswordHash(credentials.Password, user.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password"})
			return
		}

		// Assuming GenerateToken is a function that creates a new JWT token
		token, err := utils.GenerateToken(user.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}

func UpdateCurrentUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract Basic Auth credentials
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Basic" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}

		decoded, err := base64.StdEncoding.DecodeString(headerParts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Base64 credentials"})
			return
		}

		creds := strings.SplitN(string(decoded), ":", 2)
		if len(creds) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials format"})
			return
		}
		username, password := creds[0], creds[1]

		// Authenticate the user
		var user models.User
		if err := db.Where("username = ?", username).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			return
		}

		// Verify the password
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		// Proceed with updating the user's information
		var input map[string]interface{}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Validate and hash password if present in input
		if pass, ok := input["password"].(string); ok && pass != "" {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
				return
			}
			input["password"] = string(hashedPassword)
		} else if ok {
			delete(input, "password")
		}

		// Perform the update
		result := db.Model(&user).Updates(input)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user information"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
	}
}

func GetCurrentUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		// Extracting the encoded credentials from the Authorization header
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Basic" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}

		// Decoding the Base64 encoded credentials
		decoded, err := base64.StdEncoding.DecodeString(headerParts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Base64 credentials"})
			return
		}

		// Splitting the decoded credentials into username and password
		creds := strings.SplitN(string(decoded), ":", 2)
		if len(creds) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials format"})
			return
		}
		username, password := creds[0], creds[1]

		// Retrieve the user from the database using the username
		var user models.User
		result := db.Where("username = ?", username).First(&user)
		if result.Error != nil || result.RowsAffected == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		// Comparing the provided password with the hashed password in the database
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		// Responding with user details if authentication is successful
		c.JSON(http.StatusOK, gin.H{
			"id":              user.ID,
			"first_name":      user.FirstName,
			"last_name":       user.LastName,
			"username":        user.Username,
			"account_created": user.AccountCreated.Format(time.RFC3339),
			"account_updated": user.AccountUpdated.Format(time.RFC3339),
		})
	}
}

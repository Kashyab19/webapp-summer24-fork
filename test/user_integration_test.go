package test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"webapp/controllers" // Ensure this is the correct import path for your controllers package
	"webapp/models"      // Ensure this is the correct import path for your models package

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func setupDatabase() *gorm.DB {
	dsn := "host=localhost user=kashyabmurali dbname=postgres password=postgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		os.Exit(1)
	}
	db.AutoMigrate(&models.User{})
	return db
}

func setupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.POST("/v1/user", controllers.CreateUser(db))
	r.GET("/v1/user/self", controllers.GetCurrentUser(db))
	r.PUT("/v1/user/self", controllers.UpdateCurrentUser(db))
	return r
}

func TestMain(m *testing.M) {
	db = setupDatabase()
	code := m.Run()
	db.Exec("DELETE FROM users") // Cleanup database
	os.Exit(code)
}

func TestCreateUser(t *testing.T) {
	router := setupRouter(db)

	user := models.User{
		FirstName: "Test",
		LastName:  "User",
		Username:  "test2@example.com",
		Password:  "password123",
	}

	userBytes, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/v1/user", bytes.NewBuffer(userBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetUserSelf(t *testing.T) {
	router := setupRouter(db)

	credentials := "test2@example.com:password123"
	encodedCredentials := base64.StdEncoding.EncodeToString([]byte(credentials))

	req, _ := http.NewRequest("GET", "/v1/user/self", nil)
	req.Header.Set("Authorization", "Basic "+encodedCredentials)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	fmt.Println("Response status:", w.Code)
	fmt.Println("Response body:", w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateUser(t *testing.T) {
	router := setupRouter(db)

	credentials := "test2@example.com:password123"
	encodedCredentials := base64.StdEncoding.EncodeToString([]byte(credentials))

	updateInfo := map[string]string{
		"first_name": "Updated",
		"last_name":  "User",
		"password":   "newPassword123",
	}

	updateBytes, _ := json.Marshal(updateInfo)
	req, _ := http.NewRequest("PUT", "/v1/user/self", bytes.NewBuffer(updateBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+encodedCredentials)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

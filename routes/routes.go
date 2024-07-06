package routes

import (
	"webapp/controllers"
	"webapp/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	// Middleware for handling authentication
	authMiddleware := utils.TokenAuthMiddleware()

	// Public routes
	r.POST("/v1/user", controllers.CreateUser(db))

	r.PUT("/v1/user/bearer-token", authMiddleware, controllers.UpdateCurrentUser(db))
	r.POST("/login", controllers.Login(db))

	r.PUT("/v1/user/self", controllers.UpdateCurrentUser(db))
	r.GET("/v1/user/self", controllers.GetCurrentUser(db))

	return r
}

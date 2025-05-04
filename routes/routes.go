package routes

import (
	"ecommerce-go/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	api.POST("/register", controllers.Register)
	api.POST("/login", controllers.Login)
}

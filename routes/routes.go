package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/huelton/ecommerce-go/controllers"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	api.POST("/register", controllers.Register)
}

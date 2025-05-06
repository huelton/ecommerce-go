package routes

import (
	"ecommerce/controllers"
	"ecommerce/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	api.POST("/register", controllers.Register)
	api.POST("/login", controllers.Login)

	api.GET("/me", middleware.Autenticated(), func(c *gin.Context) {
		userID := c.GetInt("user_id")
		isAdmin := c.GetBool("is_admin")
		c.JSON(200, gin.H{"user_id": userID, "is_admin": isAdmin})
	})

	//Products
	api.GET("/products", controllers.ListProducts)
	api.POST("/products", middleware.Autenticated(), controllers.CreateProduct)

}

package routes

import (
	"ecommerce/controllers"
	"ecommerce/middleware"
	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")

	api.GET("/exemplo", func(c *gin.Context) {
		txn := nrgin.Transaction(c)
		txn.AddAttribute("rota", "/exemplo")
		c.JSON(200, gin.H{"mensagem": "rota monitorada"})
	})

	api.GET("/produtos", func(c *gin.Context) {
		txn := nrgin.Transaction(c)
		txn.AddAttribute("rota", "/produtos")
		txn.AddAttribute("usuario_id", "123")
		c.JSON(200, gin.H{"mensagem": "Produtos listados"})
	})

	api.POST("/register", controllers.Register)
	api.POST("/login", controllers.Login)

	api.GET("/me", middleware.Autenticated(), func(c *gin.Context) {
		userID := c.GetInt("user_id")
		isAdmin := c.GetBool("is_admin")
		c.JSON(200, gin.H{"user_id": userID, "is_admin": isAdmin})
	})

	//Products
	api.POST("/products", middleware.Autenticated(), controllers.CreateProduct)
	api.GET("/products", controllers.ListProducts)

	//Orders
	api.POST("/orders", middleware.Autenticated(), controllers.CreateOrder)
	api.GET("/orders", middleware.Autenticated(), controllers.ListOrdersUser)

	//Payment Order
	api.PUT("/orders/:id/payment", middleware.Autenticated(), controllers.OrderPayment)

	//Cancel Order
	api.PUT("/orders/:id/cancel", middleware.Autenticated(), controllers.CancelOrder)

	//Admin routes
	api.GET("/admin/orders", middleware.Autenticated(), controllers.ListAllOrdersAdmin)
	api.GET("/admin/dashboard/orders", controllers.CountAllOrdersAdmin)

}

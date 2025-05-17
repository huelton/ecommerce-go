package main

import (
	"ecommerce/config"
	"ecommerce/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	config.ConectDataBase()

	r := gin.Default()
	r.Use(CORSMiddleware())
	routes.RegisterRoutes(r)

	r.Run(":8080")
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Para requisições OPTIONS, responda diretamente
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

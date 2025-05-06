package main

import (
	"ecommerce/config"
	"ecommerce/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	config.ConectDataBase()

	r := gin.Default()
	routes.RegisterRoutes(r)

	r.Run(":8080")
}

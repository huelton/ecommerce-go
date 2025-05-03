package main

import (
	"github.com/gin-gonic/gin"
	"github.com/huelton/ecommerce-go/config"
	"github.com/huelton/ecommerce-go/routes"
)

func main() {
	config.ConectDataBase()

	r := gin.Default()
	routes.RegisterRoutes(r)

	r.Run(":8080")
}

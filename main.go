package main

import (
	"ecommerce/config"
	"ecommerce/docs"
	"ecommerce/routes"
	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
)

func main() {
	// Inicializa o New Relic
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("Ecommerce-Golang"),
		newrelic.ConfigLicense("804eea2d000fca630fc87e9fbccf62b7FFFFNRAL"), // use variável de ambiente
		newrelic.ConfigDistributedTracerEnabled(true),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Swagger
	docs.SwaggerInfo.Title = "Sistema Ecommerce GOLANG"
	docs.SwaggerInfo.Description = "This is a sample server Ecommerce server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// Inicializa banco de dados
	config.ConectDataBase()

	// Inicializa servidor
	r := gin.Default()
	r.Use(CORSMiddleware())

	// Adiciona middleware do New Relic
	r.Use(nrgin.Middleware(app))

	// Rotas
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Rotas da aplicação
	routes.RegisterRoutes(r)

	// Inicia servidor
	log.Println("Servidor iniciado em http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, x-api-key")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

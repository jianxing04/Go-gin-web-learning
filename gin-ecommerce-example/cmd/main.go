package main

import (
	"gin-ecommerce-example/internal/handlers"
	"gin-ecommerce-example/internal/models"
	"gin-ecommerce-example/pkg/config"
	"gin-ecommerce-example/pkg/database"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Config load error:", err)
	}
	if err := database.ConnectMySQL(); err != nil {
		log.Fatal("MySQL connect error:", err)
	}
	database.ConnectRedis()
	if err := database.ConnectES(); err != nil {
		log.Fatal("ES connect error:", err)
	}

	// 自动迁移模型
	database.DB.AutoMigrate(&models.Product{})

	r := gin.Default()
	r.Static("/frontend", "./frontend")
	r.GET("/", func(c *gin.Context) {
		c.File("./frontend/index.html")
	})
	api := r.Group("/api")
	{
		api.POST("/products", handlers.CreateProductHandler)
		api.GET("/products/:id", handlers.GetProductHandler)
		api.GET("/products", handlers.GetAllProductsHandler)
		api.GET("/search", handlers.SearchProductsHandler) // e.g., /search?q=keyword
		api.PUT("/products/:id", handlers.UpdateProductHandler)
		api.DELETE("/products/:id", handlers.DeleteProductHandler)
	}

	r.Run(config.AppConfig.ServerPort)
}

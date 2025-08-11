package main

import (
	"employee-crud/config"
	"employee-crud/database"
	"employee-crud/handlers"
	"employee-crud/models"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()
	database.InitDB()
	if err := database.DB.AutoMigrate(&models.Employee{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})
	r.Static("/static", "./static")
	api := r.Group("/api/employees")
	{
		api.POST("/", handlers.CreateEmployee)
		api.GET("/", handlers.GetEmployees)
		api.GET("/:id", handlers.GetEmployeeByID)
		api.PUT("/:id", handlers.UpdateEmployee)
		api.DELETE("/:id", handlers.DeleteEmployee)
	}
	addr := fmt.Sprintf(":%d", config.Cfg.Server.Port)
	log.Printf("Starting server on %s", addr)
	r.Run(addr)
}

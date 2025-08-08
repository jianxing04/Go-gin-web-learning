package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserRequest struct {
	UserInput string `json:"user_input"`
}

func main() {
	router := gin.Default()
	router.Static("/", "./static")
	api := router.Group("/api")
	{
		api.POST("/process", func(c *gin.Context) {
			var requestData UserRequest
			if err := c.ShouldBindJSON(&requestData); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Invalid input",
				})
				return
			}
			processedData := "Processed: " + requestData.UserInput
			c.JSON(http.StatusOK, gin.H{
				"message":        "Data processed successfully",
				"processed_data": processedData,
			})
		})
	}
	router.Run(":8080")
}

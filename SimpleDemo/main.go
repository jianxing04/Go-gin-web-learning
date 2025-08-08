package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func main(){
	router:=gin.Default()
	router.GET("/",func(c *gin.Context){
		c.JSON(http.StatusOK,gin.H{
			"message":"Welcome to the Gin web server!",
		})
	})
	router.GET("/ping",func(c *gin.Context){
		c.JSON(http.StatusOK,gin.H{
			"message":"pong",
		})
	})
	router.GET("/hello/:name",func(c *gin.Context){
		name:=c.Param("name")
		c.JSON(http.StatusOK,gin.H{
			"message":"Hello,"+name+"!",
		})
	})
	router.Run(":8080")
}
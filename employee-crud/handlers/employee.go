package handlers

import (
	"net/http"

	"employee-crud/database"
	"employee-crud/models"

	"github.com/gin-gonic/gin"
)

func CreateEmployee(c *gin.Context) {
	var input struct {
		Name  string `json:"name" binding:"required"`
		Age   int    `json:"age"`
		Email string `json:"email"`
		Dept  string `json:"dept"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	emp := models.Employee{
		Name:  input.Name,
		Age:   input.Age,
		Email: input.Email,
		Dept:  input.Dept,
	}
	if err := database.DB.Create(&emp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create employee",
		})
		return
	}
	c.JSON(http.StatusOK, emp)
}

func GetEmployees(c *gin.Context) {
	var employees []models.Employee
	database.DB.Find(&employees)
	c.JSON(http.StatusOK, employees)
}

func GetEmployeeByID(c *gin.Context) {
	id := c.Param("id")
	var emp models.Employee
	if err := database.DB.First(&emp, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Employee not found",
		})
		return
	}
	c.JSON(http.StatusOK, emp)
}

func UpdateEmployee(c *gin.Context) {
	id := c.Param("id")
	var emp models.Employee
	if err := database.DB.First(&emp, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Employee not found",
		})
		return
	}
	var input struct {
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Email string `json:"email"`
		Dept  string `json:"dept"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	updates := models.Employee{
		Name:  input.Name,
		Age:   input.Age,
		Email: input.Email,
		Dept:  input.Dept,
	}
	if err := database.DB.Model(&emp).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update employee",
		})
		return
	}
	c.JSON(http.StatusOK, emp)
}

func DeleteEmployee(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.Employee{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Employee not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Employee deleted successfully",
	})
}

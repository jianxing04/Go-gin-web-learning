package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func main() {
	InitConfig()

	r := gin.Default()
	r.Static("/static", "./static")

	r.POST("/login", login)

	authorized := r.Group("/api", AuthMiddleware)
	authorized.POST("/employees", createEmployee)
	authorized.GET("/employees", getEmployees)
	authorized.PUT("/employees/:id", updateEmployee)
	authorized.DELETE("/employees/:id", deleteEmployee)

	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	r.Run(":8080")
}

func login(c *gin.Context) {
	var loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名和密码不能为空"})
		return
	}

	var user User
	if err := DB.Where("username = ?", loginRequest.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
		return
	}

	if !user.CheckPassword(loginRequest.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "密码错误"})
		return
	}

	token, err := generateToken(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"token":   token,
	})
}

func generateToken(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "employee-app",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func AuthMiddleware(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供授权头"})
		c.Abort()
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的授权头"})
		c.Abort()
		return
	}

	tokenString := parts[1]
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的token"})
		c.Abort()
		return
	}

	c.Set("username", claims.Username)
	c.Next()
}

func createEmployee(c *gin.Context) {
	var emp Employee
	if err := c.ShouldBindJSON(&emp); err != nil {
		c.JSON(http.StatusBadRequest, Response{Code: 400, Msg: err.Error()})
		return
	}
	DB.Create(&emp)
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "添加成功", Data: emp})
}

func getEmployees(c *gin.Context) {
	var employees []Employee
	var total int64

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "5"))
	name := c.Query("name")

	query := DB.Model(&Employee{})
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	query.Count(&total)

	query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&employees)

	c.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  "查询成功",
		Data: gin.H{
			"list":  employees,
			"total": total,
		},
	})
}

func updateEmployee(c *gin.Context) {
	id := c.Param("id")
	var emp Employee
	if err := DB.First(&emp, id).Error; err != nil {
		c.JSON(http.StatusNotFound, Response{Code: 404, Msg: "员工不存在"})
		return
	}
	var input Employee
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, Response{Code: 400, Msg: err.Error()})
		return
	}
	emp.Name = input.Name
	emp.Position = input.Position
	emp.Salary = input.Salary
	emp.HireDate = input.HireDate
	DB.Save(&emp)
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "更新成功", Data: emp})
}

func deleteEmployee(c *gin.Context) {
	id := c.Param("id")
	if err := DB.Delete(&Employee{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, Response{Code: 500, Msg: "删除失败"})
		return
	}
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "删除成功"})
}

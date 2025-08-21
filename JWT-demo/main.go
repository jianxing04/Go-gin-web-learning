package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// 1. 定义用户模型
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"` // 序列化时忽略密码（安全起见）
}

// 2. 模拟数据库（实际项目用真实数据库）
var userStore = make(map[string]User) // key: username, value: User
var lastUserID = 0

// 3. JWT配置（实际项目中密钥应放在环境变量，不要硬编码）
const (
	JWTSecretKey = "your-secret-key-keep-safe" // 密钥（生产环境需更换为复杂密钥）
	TokenExpiry  = 24 * time.Hour              // Token有效期24小时
)

// 自定义JWT声明（包含用户信息）
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// 4. 生成JWT Token
func generateToken(user User) (string, error) {
	// 设置过期时间
	expirationTime := time.Now().Add(TokenExpiry)

	// 创建自定义声明
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),     // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),     // 生效时间（立即生效）
			Issuer:    "gin-jwt-example",                  // 签发者
		},
	}

	// 创建Token（使用HS256算法）
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名Token
	tokenString, err := token.SignedString([]byte(JWTSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// 5. 验证JWT Token并解析用户信息
func validateToken(tokenString string) (*Claims, error) {
	// 解析Token
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			// 验证签名算法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(JWTSecretKey), nil
		},
	)

	if err != nil {
		return nil, err
	}

	// 验证claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// 6. JWT验证中间件（保护需要登录的接口）
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// 格式要求：Bearer <token>（注意空格）
		var tokenString string
		fmt.Sscanf(authHeader, "Bearer %s", &tokenString)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization format must be Bearer {token}"})
			c.Abort()
			return
		}

		// 验证Token
		claims, err := validateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		// 将用户信息存入上下文，供后续接口使用
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)

		c.Next() // 继续执行后续处理
	}
}

// 7. 注册接口处理函数
func registerHandler(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"` // 必须提供用户名
		Password string `json:"password" binding:"required"` // 必须提供密码
	}

	// 绑定并验证请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户名是否已存在
	if _, exists := userStore[req.Username]; exists {
		c.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
		return
	}

	// 密码加密（使用bcrypt，自动加盐）
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost, // 加密强度（默认10）
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	// 创建用户
	lastUserID++
	user := User{
		ID:       fmt.Sprintf("%d", lastUserID),
		Username: req.Username,
		Password: string(hashedPassword), // 存储加密后的密码
	}
	userStore[user.Username] = user

	// 返回成功信息（不包含密码）
	c.JSON(http.StatusCreated, gin.H{
		"message": "user registered successfully",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}

// 8. 登录接口处理函数
func loginHandler(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// 绑定请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查询用户
	user, exists := userStore[req.Username]
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	// 验证密码（对比明文与加密后的密码）
	err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	// 生成JWT Token
	token, err := generateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	// 返回Token和用户信息
	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"token":   token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}

// 9. 受保护的示例接口（需要JWT验证）
func profileHandler(c *gin.Context) {
	// 从上下文获取用户信息（JWT中间件已存入）
	userID, _ := c.Get("userID")
	username, _ := c.Get("username")

	c.JSON(http.StatusOK, gin.H{
		"message": "this is a protected profile endpoint",
		"user": gin.H{
			"id":       userID,
			"username": username,
		},
	})
}

func main() {
	r := gin.Default()

	// 公开接口（无需登录）
	public := r.Group("/api")
	{
		public.POST("/register", registerHandler) // 注册
		public.POST("/login", loginHandler)       // 登录
	}

	// 受保护接口（需要JWT验证）
	protected := r.Group("/api")
	protected.Use(JWTAuthMiddleware()) // 应用JWT中间件
	{
		protected.GET("/profile", profileHandler) // 个人信息接口
	}

	// 启动服务
	r.Run(":8080")
}


/*
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"username":"alice", "password":"123456"}'

curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"alice", "password":"123456"}'

curl http://localhost:8080/api/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."  # 替换为实际Token
*/
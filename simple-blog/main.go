package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Blog struct {
	ID        int       `json:"id"`
	Text      string    `json:"text"`
	ImageURL  string    `json:"image_url"`
	VideoURL  string    `json:"video_url"`
	CreatedAt time.Time `json:"created_at"`
}

var (
	blogs []Blog
	mu    sync.Mutex
)

const pageSize = 5 // 每页显示条数

// 保存上传文件
func saveFile(c *gin.Context, formKey, uploadDir string) (string, error) {
	file, err := c.FormFile(formKey)
	if err != nil {
		return "", nil
	}
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", err
	}
	filePath := filepath.Join(uploadDir, file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		return "", err
	}
	return "/" + uploadDir + "/" + file.Filename, nil
}

func main() {
	r := gin.Default()

	// 提供静态资源：上传的文件
	r.Static("/uploads/images", "./uploads/images")
	r.Static("/uploads/videos", "./uploads/videos")

	// 提供前端静态页面
	r.Static("/app", "./frontend")

	// 获取分页博客
	r.GET("/api/blogs", func(c *gin.Context) {
		pageStr := c.Query("page")
		page, _ := strconv.Atoi(pageStr)
		if page < 1 {
			page = 1
		}

		mu.Lock()
		defer mu.Unlock()

		start := (page - 1) * pageSize
		end := start + pageSize
		if start >= len(blogs) {
			c.JSON(http.StatusOK, []Blog{})
			return
		}
		if end > len(blogs) {
			end = len(blogs)
		}
		c.JSON(http.StatusOK, blogs[start:end])
	})

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/app/index.html")
		// 或者使用 http.StatusFound (302) 进行临时重定向
		// c.Redirect(http.StatusFound, "/app/index.html")
	})

	// 发布博客
	r.POST("/api/blog", func(c *gin.Context) {
		text := c.PostForm("text")
		imageURL, _ := saveFile(c, "image", "uploads/images")
		videoURL, _ := saveFile(c, "video", "uploads/videos")

		mu.Lock()
		newBlog := Blog{
			ID:        len(blogs) + 1,
			Text:      text,
			ImageURL:  imageURL,
			VideoURL:  videoURL,
			CreatedAt: time.Now(),
		}
		blogs = append([]Blog{newBlog}, blogs...) // 最新在前
		mu.Unlock()

		c.JSON(http.StatusOK, gin.H{"message": "博客发布成功"})
	})

	r.Run(":8080")
}

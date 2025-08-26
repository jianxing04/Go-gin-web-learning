package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建存放图片的目录
	const uploadDir = "./uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		panic("创建上传目录失败: " + err.Error())
	}

	r := gin.Default()

	// 上传图片接口
	r.POST("/upload", func(c *gin.Context) {
		// 取到上传的文件
		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "获取图片失败"})
			return
		}

		// 拼接保存路径
		filePath := filepath.Join(uploadDir, file.Filename)

		// 保存文件
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
			return
		}

		// 返回访问路径
		c.JSON(http.StatusOK, gin.H{
			"message": "上传成功",
			"url":     fmt.Sprintf("/images/%s", file.Filename),
		})
	})

	// 提供静态文件服务，让前端可以访问已上传的图片
	r.Static("/images", uploadDir)

	// 启动服务
	r.Run(":8080")
}

/*
curl -X POST -F "image=@\"/mnt/c/Users/weijianxing/Pictures/Saved Pictures/test.jpg\"" http://localhost:8080/upload

*/

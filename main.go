package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
	"os"
	"path/filepath"
)

var redisClient *redis.Client

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func main() {
	router := gin.Default()

	router.Use(CORSMiddleware())
	router.POST("/enqueue", EnqueueHandler)
	router.GET("/get_enqueue_items", GetEnqueueItemsHandler)
	router.GET("/list_reports", ListReportsHandler)

	router.Run(":8080")
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func EnqueueHandler(c *gin.Context) {
	var request struct {
		ImageName string `json:"image_name" binding:"required"`
	}

	// Đọc dữ liệu từ request body
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	items, err := redisClient.LRange(context.Background(), "image_queue", 0, -1).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get enqueue items"})
		return
	}

	for _, item := range items {
		if item == request.ImageName {
			c.JSON(http.StatusOK, gin.H{"message": "Item already exists"})
			return
		}
	}

	err = redisClient.RPush(context.Background(), "image_queue", request.ImageName).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enqueue item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item enqueued successfully"})
}

func GetEnqueueItemsHandler(c *gin.Context) {
	items, err := redisClient.LRange(context.Background(), "image_queue", 0, -1).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get enqueue items"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

func ListReportsHandler(c *gin.Context) {
	dirPath := "./reports"

	files, err := ListItemsInDirectory(dirPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error listing files"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"files": files})
}

func ListItemsInDirectory(dirPath string) ([]string, error) {
	var items []string

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		items = append(items, path)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return items, nil
}

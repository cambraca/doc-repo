package service

import (
	"api/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

func RunHttpServer(version, buildTime string) {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		//AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Length", "Content-Type"},
		//AllowCredentials: false,
		//MaxAge:           12 * time.Hour,
	}))

	r.GET("/ping", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	r.GET("/status", func(c *gin.Context) {
		c.JSON(200, struct {
			Version    string `json:"version"`
			Status     string `json:"status"`
			BuildTime  string `json:"build_time"`
			BucketName string `json:"bucket_name"`
		}{
			Version:    version,
			Status:     "ok",
			BuildTime:  buildTime,
			BucketName: os.Getenv("DOCUMENTS_BUCKET_NAME"),
		})
	})

	r.GET("/api/v1/documents", controllers.DocumentsIndex)
	r.POST("/api/v1/documents", controllers.DocumentsCreate)

	err := r.Run()
	if err != nil {
		log.Fatalf("Server error: %s", err)
	}
}

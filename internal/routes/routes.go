package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/krishnadwypayan/shorturl/internal/snowflake"
)

func RegisterSnowflakeRoutes(r *gin.Engine, generator *snowflake.Generator) {
	// Route: Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Route: Generate a new ID
	r.GET("/generate", func(c *gin.Context) {
		id := generator.Next()
		c.JSON(http.StatusOK, gin.H{"id": id})
	})
}

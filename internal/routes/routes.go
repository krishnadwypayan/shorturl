package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/krishnadwypayan/shorturl/internal/model"
	"github.com/krishnadwypayan/shorturl/internal/shortify"
	"github.com/krishnadwypayan/shorturl/internal/snowflake"
)

func RegisterSnowflakeRoutes(r *gin.Engine, generator *snowflake.Generator) {
	// Route: Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Route: Generate a new ID
	r.GET("/generate", func(c *gin.Context) {
		res := model.SnowflakeResponse{
			ID: generator.NextString(),
		}
		c.JSON(http.StatusOK, res)
	})
}

func RegisterShortURLRoutes(r *gin.Engine) {
	// Route: Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Route: Generate a new short URL ID
	r.POST("/shortify", func(c *gin.Context) {
		var req model.ShortURLRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		res, err := shortify.Shortify(req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, res)
		}
	})
}

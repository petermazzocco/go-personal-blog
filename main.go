package main

import (
	"net/http"
	"personal-blog/initializers"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.InitENV()
	initializers.InitDB()
}

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Landing Page",
		})
	})

	r.POST("/new", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Created message",
		})
	})

	r.POST("/edit/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Updated message",
		})
	})

	r.GET("/admin", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Admin portal",
		})
	})
	r.Run()
}

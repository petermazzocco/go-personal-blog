package main

import (
	"net/http"
	"personal-blog/auth"
	"personal-blog/initializers"
	"personal-blog/renderer"
	"personal-blog/views"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.InitENV()
	initializers.InitDB()
}

func main() {
	r := gin.Default()
	r.LoadHTMLFiles("views/*.html")
	htmlRenderer := r.HTMLRender
	r.HTMLRender = &renderer.HTMLTemplRenderer{FallbackHtmlRenderer: htmlRenderer}

	// Disable trusted proxies warning.
	r.SetTrustedProxies(nil)

	// Serve static files
	r.Static("/static", "./static")

	// Serve your templates
	r.GET("/", func(c *gin.Context) {
		views.Index().Render(c.Request.Context(), c.Writer)
	})

	// Auth pages
	r.GET("/signup", func(c *gin.Context) {
		views.SignUp().Render(c.Request.Context(), c.Writer)
	})
	r.GET("/signin", func(c *gin.Context) {
		views.SignIn().Render(c.Request.Context(), c.Writer)
	})

	// User Authentication
	r.POST("/signin", auth.SignInWithCredentials)
	r.POST("/signup", auth.SignUpWithCredentials)

	// Protected API Routes
	api := r.Group("/api")
	api.Use(auth.Middleware())
	{
		// Create a new post
		api.GET("/new", func(c *gin.Context) {
			views.NewPost().Render(c.Request.Context(), c.Writer)
		})
		api.POST("/new", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Created message",
			})
		})

		// View or edit a post
		api.GET("/post/:id", func(c *gin.Context) {
			views.ViewPost().Render(c.Request.Context(), c.Writer)
		})
		api.POST("/post/:id", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Updated message",
			})
		})

		// Admin portal
		api.GET("/admin", func(c *gin.Context) {
			views.AdminPortal().Render(c.Request.Context(), c.Writer)
		})
	}

	r.Run(":8080")
}

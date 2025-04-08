package main

import (
	"net/http"
	"personal-blog/auth"
	"personal-blog/initializers"
	"personal-blog/models"
	"personal-blog/renderer"
	"personal-blog/views"
	"strconv"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.InitENV()
	initializers.InitDB()
}

func main() {
	// Instantiate gin router, loading html and renderer
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

	// Auth Pages
	r.GET("/signup", func(c *gin.Context) {
		views.SignUp().Render(c.Request.Context(), c.Writer)
	})
	r.GET("/signin", func(c *gin.Context) {
		views.SignIn().Render(c.Request.Context(), c.Writer)
	})

	// User Authentication Routes
	r.POST("/signin", auth.SignInWithCredentials)
	r.POST("/signup", auth.SignUpWithCredentials)

	// Protected API Routes
	api := r.Group("/authenticated")
	api.Use(auth.Middleware())
	{
		// New post page
		api.GET("/new", func(c *gin.Context) {
			views.NewPost().Render(c.Request.Context(), c.Writer)
		})
		// New post route
		api.POST("/new", func(c *gin.Context) {
			//Get form values
			title := c.PostForm("title")
			content := c.PostForm("content")

			// Get user id from cookie
			userIDValue, exists := c.Get("userId")
			if !exists {
				views.NewPostError("Error occured while creating the post. Try again").Render(c.Request.Context(), c.Writer)
				return
			}
			userIDStr := userIDValue.(string)
			userID, _ := strconv.ParseUint(userIDStr, 10, 64)

			// Query the user in the db to verify the user is actually in the database
			var user models.User
			if err := initializers.DB.First(&user, userID).Error; err != nil {
				views.NewPostError("Error occured while creating the post. Try again").Render(c.Request.Context(), c.Writer)
				return
			}

			// Implement logic to create a new post
			post := models.Post{Title: title, Content: content}
			if err := initializers.DB.Create(&post).Error; err != nil {
				views.NewPostError("Error occured while creating the post. Try again").Render(c.Request.Context(), c.Writer)
				return
			}

			// Swap inner html
			views.SuccessNewPost().Render(c.Request.Context(), c.Writer)
		})

		// View a post page
		api.GET("/post/:id", func(c *gin.Context) {
			views.ViewPost().Render(c.Request.Context(), c.Writer)
		})
		// Edit a post route
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

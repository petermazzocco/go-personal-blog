package main

import (
	"fmt"
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
			userIDValue, exists := c.Get("userID")
			if !exists {
				fmt.Println("UserID is not found")
				views.NewPostError("Error occured while creating the post. Try again").Render(c.Request.Context(), c.Writer)
				return
			}
			userIDStr := userIDValue.(string)
			userID, _ := strconv.ParseUint(userIDStr, 10, 64)

			// Query the user in the db to verify the user is actually in the database
			var user models.User
			if err := initializers.DB.First(&user, userID).Error; err != nil {
				fmt.Println("Error occured querying db for user", err.Error())
				views.NewPostError("Error occured while creating the post. Try again").Render(c.Request.Context(), c.Writer)
				return
			}

			// Implement logic to create a new post
			post := models.Post{Title: title, Content: content, UserID: uint(userID)}
			if err := initializers.DB.Create(&post).Error; err != nil {
				fmt.Println("Error occured creating post", err.Error())
				views.NewPostError("Error occured while creating the post. Try again").Render(c.Request.Context(), c.Writer)
				return
			}

			// Swap inner html
			views.SuccessNewPost().Render(c.Request.Context(), c.Writer)
		})

		// View all posts
		api.GET("/posts", func(c *gin.Context) {
			var posts []models.Post
			if err := initializers.DB.Find(&posts).Error; err != nil {
				fmt.Println("Error occured querying db for posts", err.Error())
				views.NewPostError("Error occured while fetching the posts. Try again").Render(c.Request.Context(), c.Writer)
				return
			}
			views.ViewPosts(posts).Render(c.Request.Context(), c.Writer)
		})

		// View a post page
		api.GET("/posts/:id", func(c *gin.Context) {
			var post models.Post
			id := c.Param("id")
			if err := initializers.DB.First(&post, id).Error; err != nil {
				fmt.Println("Error occured querying db for post", err.Error())
				views.NewPostError("Error occured while fetching the post. Try again").Render(c.Request.Context(), c.Writer)
				return
			}
			views.ViewPost(post).Render(c.Request.Context(), c.Writer)
		})
		// Edit a post route
		api.POST("/posts/:id", func(c *gin.Context) {
			title := c.PostForm("title")
			content := c.PostForm("content")

			var post models.Post
			id := c.Param("id")
			if err := initializers.DB.First(&post, id).Error; err != nil {
				fmt.Println("Error occured querying db for post", err.Error())
				views.NewPostError("Error occured while fetching the post. Try again").Render(c.Request.Context(), c.Writer)
				return
			}

			post.Title = title
			post.Content = content

			if err := initializers.DB.Save(&post).Error; err != nil {
				fmt.Println("Error occured while updating the post", err.Error())
				views.NewPostError("Error occured while updating the post. Try again").Render(c.Request.Context(), c.Writer)
				return
			}

			views.SuccessNewPost().Render(c.Request.Context(), c.Writer)
		})

	}

	r.Run(":8080")
}

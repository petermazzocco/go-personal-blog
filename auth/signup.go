package auth

import (
	"net/http"
	"personal-blog/helpers"
	"personal-blog/initializers"
	"personal-blog/models"
	"personal-blog/views"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func SignUpWithCredentials(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")
	var user models.User
	// Generate a salt for the user
	salt, err := helpers.GenerateSalt()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating salt"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}
	password = string(hashedPassword)

	// Store salt in database
	user.Salt = helpers.EncodeToBase64(salt)
	user.Email = email
	user.Password = password
	// Save user to DB
	if err := initializers.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	// Return a success message using HTMX
	views.SignupSuccess().Render(c.Request.Context(), c.Writer)

}

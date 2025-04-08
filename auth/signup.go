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
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Generate a salt for the user
	salt, err := helpers.GenerateSalt()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating salt"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}
	user.Password = string(hashedPassword)

	// Store salt in database
	user.Salt = helpers.EncodeToBase64(salt)

	// Save user to DB
	if err := initializers.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	// Return a success message using HTMX
	views.SignupSuccess().Render(c.Request.Context(), c.Writer)

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

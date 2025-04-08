package auth

import (
	"encoding/base64"
	"net/http"
	"personal-blog/helpers"
	"personal-blog/initializers"
	"personal-blog/models"
	"personal-blog/views"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func SignInWithCredentials(c *gin.Context) {
	// Get request data
	email := c.PostForm("email")
	password := c.PostForm("password")

	// Find user in DB
	var user models.User
	if err := initializers.DB.Where("email = ?", email).First(&user).Error; err != nil {
		views.AuthError("Invalid username or password").Render(c.Request.Context(), c.Writer)
		return
	}

	// Compare stored hashed password with the given password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		views.AuthError("Invalid username or password").Render(c.Request.Context(), c.Writer)
		return
	}

	// Decode user's salt
	salt, err := base64.StdEncoding.DecodeString(user.Salt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding salt"})
		return
	}

	// Derive the encryption key
	encryptionKey := helpers.DeriveKey(password, salt)
	encodedKey := helpers.EncodeToBase64(encryptionKey)
	// Generate a session token (JWT)
	token, err := GenerateJWT(strconv.FormatUint(uint64(user.ID), 10), encodedKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	// Set the JWT as an HTTPOnly cookie
	c.SetCookie(
		"token",
		token,
		3600*24*30, // one month
		"/",
		"",
		false, // HTTPS only
		true,  // not accessible via JavaScript
	)

	// Return success response using HTMX
	views.SigninSuccess(user.Email).Render(c.Request.Context(), c.Writer)
}

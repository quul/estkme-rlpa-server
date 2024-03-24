package apiServer

import (
	"github.com/damonto/estkme-rlpa-server/internal/config"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

func login(c *gin.Context) {
	session := sessions.Default(c)
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Validate form input
	if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Parameters can't be empty"})
		return
	}

	switch config.C.APIServerMode {
	case "singleUser":
		if username != config.C.DefaultUserName || password != config.C.DefaultPassword {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Authentication failed"})
			return
		}
		break
	case "multiUser":
		var user User
		DB.First(&user, "username = ?", username)
		if user.ID == 0 || checkPasswordHash(user.PasswordHash, user.PasswordHash) {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Authentication failed"})
			return
		}
	default:
		panic("Unreachable")
	}

	// Save the username in the session
	session.Set("user", username)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully authenticated user"})
}

func logout(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid session token"})
		return
	}
	session.Delete("user")
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func registerEnabledCheck(c *gin.Context) {
	if config.C.EnableRegister {
		c.JSON(http.StatusOK, gin.H{"message": "Register is enabled"})
	} else {
		c.JSON(http.StatusForbidden, gin.H{"message": "Register is disabled"})
	}
}

func register(c *gin.Context) {
	if !config.C.EnableRegister {
		c.JSON(http.StatusForbidden, gin.H{"message": "Register is disabled"})
		return
	}
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Validate form input
	if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Parameters can't be empty"})
		return
	}

	passwordHash, _ := hashPassword(password)

	user := &User{
		Username:     username,
		PasswordHash: passwordHash,
	}
	DB.Create(user)
	c.JSON(http.StatusOK, gin.H{"message": "Successfully registered user"})
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

package middleware

import (
	"RMS_deploy/initializers"
	Log "RMS_deploy/log"
	"RMS_deploy/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func VerifyAdmin(c *gin.Context) {

	log := Log.GetLogger(c)
	log.Info("verifying Admin role ")

	// Retrieve user from context
	userDetail, exists := c.Get(initializers.UserString)
	if !exists {
		log.Info("user not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}

	// Correct type assertion for models.Users
	user, ok := userDetail.(models.Users)
	if !ok {
		log.Info("invalid user type in context")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid user type in context"})
		return
	}

	// Verify user role
	if user.Role != models.Admin {
		log.Info("invalid user role in context")
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden: admin access required"})
		return
	}

	log.Info("user verified as admin")
	c.Next()
}
func VerifySubAdmin(c *gin.Context) {

	log := Log.GetLogger(c)
	log.Info("verifying sub admin role ")

	// Retrieve user from context
	userDetail, exists := c.Get(initializers.UserString)
	if !exists {
		log.Info("user not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}

	// Correct type assertion for models.Users
	user, ok := userDetail.(models.Users)
	if !ok {
		log.Info("invalid user type in context")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid user type in context"})
		return
	}

	// Verify user role
	if user.Role != models.SubAdmin {
		log.Info("invalid user role in context")
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden: sub admin access required"})
		return
	}

	log.Info("user verified as sub admin")
	c.Next()
}

func VerifyUser(c *gin.Context) {

	log := Log.GetLogger(c)
	log.Info("verifying user role ")
	// Retrieve user from context
	userDetail, exists := c.Get(initializers.UserString)
	if !exists {
		log.Info("user not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}

	// Correct type assertion for models.Users
	user, ok := userDetail.(models.Users)
	if !ok {
		log.Info("invalid user type in context")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid user type in context"})
		return
	}

	// Verify user role
	if user.Role != models.User {
		log.Info("invalid user role in context")
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden: user access required"})
		return
	}

	log.Info("user verified as user")
	c.Next()
}

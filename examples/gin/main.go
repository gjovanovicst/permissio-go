// Package main demonstrates the Permissio.io Go SDK with Gin middleware.
package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/permissio/permissio-go/pkg/config"
	"github.com/permissio/permissio-go/pkg/enforcement"
	"github.com/permissio/permissio-go/pkg/models"
	"github.com/permissio/permissio-go/pkg/permissio"
	"go.uber.org/zap"
)

var permisClient *permissio.Client

// UserIn represents the input for user registration.
type UserIn struct {
	Email     string `json:"email" binding:"required,email"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

func main() {
	// Load environment variables
	apiKey := os.Getenv("PERMIS_API_KEY")
	if apiKey == "" {
		panic("PERMIS_API_KEY environment variable is required")
	}

	// Create logger
	logger, _ := zap.NewDevelopment()

	// Initialize Permis client
	cfg := config.NewConfigBuilder(apiKey).
		WithApiUrl("http://localhost:3001").
		WithDebug(true).
		WithLogger(logger).
		Build()

	permisClient = permissio.New(cfg)

	// Create Gin router
	router := gin.Default()

	// Health check endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, Gin with Permissio.io!"})
	})

	// Register new user
	router.POST("/register", registerUserHandler)

	// Protected endpoint: only authors can create posts
	router.POST("/posts", AuthMiddleware("create", "Post"), func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{
			"message": "Post created successfully",
		})
	})

	// Protected endpoint: only users with read permission can view posts
	router.GET("/posts", AuthMiddleware("read", "Post"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"posts": []gin.H{
				{"id": 1, "title": "First Post"},
				{"id": 2, "title": "Second Post"},
			},
		})
	})

	// Protected endpoint: only moderators can delete comments
	router.DELETE("/comments/:id", AuthMiddleware("delete", "Comment"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Comment deleted successfully",
		})
	})

	logger.Info("Server running on http://localhost:8000")
	router.Run(":8000")
}

// registerUserHandler handles user registration.
func registerUserHandler(c *gin.Context) {
	var user UserIn
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create user in Permissio.io
	createUser := models.NewUserCreate(user.Email).
		SetEmail(user.Email).
		SetFirstName(user.FirstName).
		SetLastName(user.LastName)

	newUser, err := permisClient.Api.Users.SyncUser(c.Request.Context(), *createUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sync user"})
		return
	}

	if newUser == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	// Assign the Reader role to the user in the 'default' tenant
	roleAssignment, err := permisClient.Api.Users.AssignRole(
		c.Request.Context(),
		user.Email,
		"Reader",
		"default",
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign role"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":         "User registered and role assigned",
		"user":            newUser,
		"role_assignment": roleAssignment,
	})
}

// AuthMiddleware creates a middleware for checking permissions.
func AuthMiddleware(action, resourceType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user from header
		userKey := c.GetHeader("X-User")
		if userKey == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Missing permission header: X-User",
			})
			c.Abort()
			return
		}

		// Build user and resource for permission check
		user := enforcement.UserBuilder(userKey).Build()
		resource := enforcement.ResourceBuilder(resourceType).Build()

		// Check permission
		permitted, err := permisClient.Check(user, enforcement.Action(action), resource)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Permission check failed",
			})
			c.Abort()
			return
		}

		if !permitted {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "You are not authorized to " + action + " a " + resourceType,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// AuthMiddlewareWithTenant creates a middleware for checking permissions with tenant context.
func AuthMiddlewareWithTenant(action, resourceType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userKey := c.GetHeader("X-User")
		if userKey == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Missing permission header: X-User",
			})
			c.Abort()
			return
		}

		tenantKey := c.GetHeader("X-Tenant")

		user := enforcement.UserBuilder(userKey).Build()
		resource := enforcement.ResourceBuilder(resourceType).
			WithTenant(tenantKey).
			Build()

		permitted, err := permisClient.Check(user, enforcement.Action(action), resource)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Permission check failed",
			})
			c.Abort()
			return
		}

		if !permitted {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "You are not authorized to " + action + " a " + resourceType,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// AuthMiddlewareWithResource creates a middleware for checking permissions on specific resource instances.
func AuthMiddlewareWithResource(action, resourceType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userKey := c.GetHeader("X-User")
		if userKey == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Missing permission header: X-User",
			})
			c.Abort()
			return
		}

		// Get resource instance key from URL parameter
		resourceKey := c.Param("id")

		user := enforcement.UserBuilder(userKey).Build()
		resource := enforcement.ResourceBuilder(resourceType).
			WithKey(resourceKey).
			Build()

		permitted, err := permisClient.Check(user, enforcement.Action(action), resource)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Permission check failed",
			})
			c.Abort()
			return
		}

		if !permitted {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "You are not authorized to " + action + " this " + resourceType,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

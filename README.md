# Permissio.io Go SDK

The official Go SDK for [Permissio.io](https://permissio.io) - a powerful authorization service for managing roles, permissions, and access control.

## Installation

```bash
go get github.com/permissio/permissio-go
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/permissio/permissio-go/pkg/config"
	"github.com/permissio/permissio-go/pkg/enforcement"
	"github.com/permissio/permissio-go/pkg/permissio"
)

func main() {
	// Create a new Permis client
	cfg := config.NewConfigBuilder("permis_key_your_api_key_here").
		WithProjectID("your-project-id").
		WithEnvironmentID("your-environment-id").
		Build()

	client := permissio.New(cfg)

	// Check permissions
	user := enforcement.UserBuilder("user@example.com").Build()
	resource := enforcement.ResourceBuilder("document").Build()

	allowed, err := client.Check(user, enforcement.Action("read"), resource)
	if err != nil {
		log.Fatal(err)
	}

	if allowed {
		fmt.Println("Access granted!")
	} else {
		fmt.Println("Access denied!")
	}
}
```

## Auto-Scope Detection

If you don't provide `projectId` and `environmentId`, the SDK will automatically fetch them from your API key:

```go
cfg := config.NewConfigBuilder("permis_key_your_api_key_here").Build()
client := permissio.New(cfg)

// The SDK will automatically fetch projectId and environmentId
// from the /v1/api-key/scope endpoint on first API call
```

## API Management

The SDK provides full CRUD access to all Permissio.io resources:

```go
ctx := context.Background()

// Sync a user
user, err := client.Api.Users.SyncUser(ctx, models.NewUserCreate("user@example.com").
	SetEmail("user@example.com").
	SetFirstName("John").
	SetLastName("Doe"))

// Assign a role
_, err = client.Api.Users.AssignRole(ctx, "user@example.com", "admin", "default")

// List roles
roles, err := client.Api.Roles.List(ctx, nil)

// Create a tenant
tenant, err := client.Api.Tenants.Create(ctx, &models.TenantCreate{
	Key:  "acme-corp",
	Name: "Acme Corporation",
})
```

## Gin Middleware Example

```go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/permissio/permissio-go/pkg/config"
	"github.com/permissio/permissio-go/pkg/enforcement"
	"github.com/permissio/permissio-go/pkg/permissio"
)

var permisClient *permissio.Client

func AuthMiddleware(action, resourceType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userKey := c.GetHeader("X-User")
		if userKey == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing X-User header"})
			c.Abort()
			return
		}

		user := enforcement.UserBuilder(userKey).Build()
		resource := enforcement.ResourceBuilder(resourceType).Build()

		allowed, err := permisClient.Check(user, enforcement.Action(action), resource)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Permission check failed"})
			c.Abort()
			return
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func main() {
	cfg := config.NewConfigBuilder("permis_key_your_api_key_here").Build()
	permisClient = permissio.New(cfg)

	router := gin.Default()

	// Protected endpoint - only users with "create" permission on "Post" can access
	router.POST("/posts", AuthMiddleware("create", "Post"), func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{"message": "Post created successfully"})
	})

	router.Run(":8000")
}
```

## Configuration Options

| Option | Description | Default |
|--------|-------------|---------|
| `WithApiUrl(url)` | Set the API base URL | `https://api.permissio.io` |
| `WithProjectID(id)` | Set the project ID | Auto-fetched |
| `WithEnvironmentID(id)` | Set the environment ID | Auto-fetched |
| `WithTimeout(duration)` | Set request timeout | 30 seconds |
| `WithDebug(enabled)` | Enable debug logging | false |
| `WithRetryAttempts(n)` | Set retry attempts | 3 |
| `WithLogger(logger)` | Set custom zap logger | nil |
| `WithHTTPClient(client)` | Set custom HTTP client | http.DefaultClient |

## API Reference

### Permission Checking

- `Check(user, action, resource) (bool, error)` - Simple permission check
- `CheckWithContext(ctx, user, action, resource) (bool, error)` - Check with context
- `CheckWithDetails(user, action, resource) (*CheckResponse, error)` - Check with full details
- `BulkCheck(checks) (*BulkCheckResponse, error)` - Bulk permission checks
- `GetPermissions(user) (*PermissionsResponse, error)` - Get all permissions for a user

### Users API

- `Api.Users.List(ctx, params)` - List users
- `Api.Users.Get(ctx, userKey)` - Get a user
- `Api.Users.Create(ctx, user)` - Create a user
- `Api.Users.Update(ctx, userKey, data)` - Update a user
- `Api.Users.Delete(ctx, userKey)` - Delete a user
- `Api.Users.SyncUser(ctx, user)` - Create or update a user
- `Api.Users.AssignRole(ctx, userKey, role, tenant)` - Assign a role to a user
- `Api.Users.UnassignRole(ctx, userKey, role, tenant)` - Remove a role from a user

### Roles API

- `Api.Roles.List(ctx, params)` - List roles
- `Api.Roles.Get(ctx, roleKey)` - Get a role
- `Api.Roles.Create(ctx, role)` - Create a role
- `Api.Roles.Update(ctx, roleKey, data)` - Update a role
- `Api.Roles.Delete(ctx, roleKey)` - Delete a role

### Tenants API

- `Api.Tenants.List(ctx, params)` - List tenants
- `Api.Tenants.Get(ctx, tenantKey)` - Get a tenant
- `Api.Tenants.Create(ctx, tenant)` - Create a tenant
- `Api.Tenants.Update(ctx, tenantKey, data)` - Update a tenant
- `Api.Tenants.Delete(ctx, tenantKey)` - Delete a tenant

### Resources API

- `Api.Resources.List(ctx, params)` - List resources
- `Api.Resources.Get(ctx, resourceKey)` - Get a resource
- `Api.Resources.Create(ctx, resource)` - Create a resource
- `Api.Resources.Update(ctx, resourceKey, data)` - Update a resource
- `Api.Resources.Delete(ctx, resourceKey)` - Delete a resource

### Role Assignments API

- `Api.RoleAssignments.List(ctx, params)` - List role assignments
- `Api.RoleAssignments.Assign(ctx, assignment)` - Create a role assignment
- `Api.RoleAssignments.Unassign(ctx, user, role, tenant)` - Remove a role assignment
- `Api.RoleAssignments.BulkAssign(ctx, assignments)` - Bulk create role assignments

## License

MIT License - see [LICENSE](LICENSE) for details.

// Package main demonstrates basic usage of the Permissio.io Go SDK.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/permissio/permissio-go/pkg/config"
	"github.com/permissio/permissio-go/pkg/enforcement"
	"github.com/permissio/permissio-go/pkg/models"
	"github.com/permissio/permissio-go/pkg/permissio"
	"go.uber.org/zap"
)

func main() {
	// Get API key from environment
	apiKey := os.Getenv("PERMIS_API_KEY")
	if apiKey == "" {
		log.Fatal("PERMIS_API_KEY environment variable is required")
	}

	// Create logger for debug output
	logger, _ := zap.NewDevelopment()

	// Build configuration with options
	cfg := config.NewConfigBuilder(apiKey).
		WithApiUrl("http://localhost:3001").
		WithDebug(true).
		WithLogger(logger).
		// WithProjectID("your-project-id").      // Optional: auto-fetched from API key
		// WithEnvironmentID("your-environment-id"). // Optional: auto-fetched from API key
		Build()

	// Create client
	client := permissio.New(cfg)
	ctx := context.Background()

	// Initialize the client (fetches project/environment scope from API key)
	// This is required before using API methods when projectId/environmentId 
	// are not explicitly configured.
	if err := client.Init(ctx); err != nil {
		log.Fatalf("Failed to initialize client: %v", err)
	}

	// Example 1: Create or sync a user
	fmt.Println("=== Example 1: Create/Sync User ===")
	user, err := client.Api.Users.SyncUser(ctx, *models.NewUserCreate("john@example.com").
		SetEmail("john@example.com").
		SetFirstName("John").
		SetLastName("Doe"))
	if err != nil {
		log.Printf("Failed to sync user: %v", err)
	} else {
		fmt.Printf("User synced: %s (%s)\n", user.Key, user.ID)
	}

	// Example 2: Create a tenant
	fmt.Println("\n=== Example 2: Create Tenant ===")
	tenant, err := client.Api.Tenants.Create(ctx, &models.TenantCreate{
		Key:         "acme-corp",
		Name:        "Acme Corporation",
		Description: "Our main customer",
	})
	if err != nil {
		log.Printf("Failed to create tenant: %v", err)
	} else {
		fmt.Printf("Tenant created: %s (%s)\n", tenant.Key, tenant.ID)
	}

	// Example 3: Assign a role to the user
	fmt.Println("\n=== Example 3: Assign Role ===")
	assignment, err := client.Api.Users.AssignRole(ctx, "john@example.com", "admin", "acme-corp")
	if err != nil {
		log.Printf("Failed to assign role: %v", err)
	} else {
		fmt.Printf("Role assigned: %s -> %s in %s\n",
			assignment.User, assignment.Role, assignment.Tenant)
	}

	// Example 4: Check permission
	fmt.Println("\n=== Example 4: Check Permission ===")
	userEnf := enforcement.UserBuilder("john@example.com").Build()
	resource := enforcement.ResourceBuilder("document").Build()

	allowed, err := client.Check(userEnf, enforcement.Action("read"), resource)
	if err != nil {
		log.Printf("Permission check failed: %v", err)
	} else {
		fmt.Printf("User john@example.com can read documents: %v\n", allowed)
	}

	// Example 5: Check permission with details
	fmt.Println("\n=== Example 5: Check Permission with Details ===")
	response, err := client.CheckWithDetails(ctx, userEnf, enforcement.Action("write"), resource)
	if err != nil {
		log.Printf("Permission check failed: %v", err)
	} else {
		fmt.Printf("Allowed: %v\n", response.Allowed)
		fmt.Printf("Reason: %s\n", response.Reason)
		if response.Debug != nil {
			fmt.Printf("Matched Roles: %v\n", response.Debug.MatchedRoles)
		}
	}

	// Example 6: Get all permissions for a user
	fmt.Println("\n=== Example 6: Get User Permissions ===")
	permissions, err := client.GetPermissions(ctx, models.GetPermissionsRequest{
		User:   "john@example.com",
		Tenant: "acme-corp",
	})
	if err != nil {
		log.Printf("Failed to get permissions: %v", err)
	} else {
		fmt.Printf("Roles: %v\n", permissions.Roles)
		fmt.Printf("Permissions: %v\n", permissions.Permissions)
	}

	// Example 7: Bulk permission check
	fmt.Println("\n=== Example 7: Bulk Permission Check ===")
	checks := []models.CheckRequest{
		{User: "john@example.com", Action: "read", Resource: "document"},
		{User: "john@example.com", Action: "write", Resource: "document"},
		{User: "john@example.com", Action: "delete", Resource: "document"},
	}
	bulkResponse, err := client.BulkCheck(ctx, checks)
	if err != nil {
		log.Printf("Bulk check failed: %v", err)
	} else {
		for _, result := range bulkResponse.Results {
			fmt.Printf("Action '%s' on '%s': %v\n",
				result.Request.Action, result.Request.Resource, result.Response.Allowed)
		}
	}

	// Example 8: List users
	fmt.Println("\n=== Example 8: List Users ===")
	users, err := client.Api.Users.List(ctx, &models.UserListParams{
		ListParams: models.ListParams{Page: 1, PerPage: 10},
	})
	if err != nil {
		log.Printf("Failed to list users: %v", err)
	} else {
		fmt.Printf("Total users: %d\n", users.Total)
		for _, u := range users.Data {
			fmt.Printf("  - %s (%s %s)\n", u.Key, u.FirstName, u.LastName)
		}
	}

	// Example 9: List roles
	fmt.Println("\n=== Example 9: List Roles ===")
	roles, err := client.Api.Roles.List(ctx, nil)
	if err != nil {
		log.Printf("Failed to list roles: %v", err)
	} else {
		fmt.Printf("Total roles: %d\n", roles.Total)
		for _, r := range roles.Data {
			fmt.Printf("  - %s (permissions: %v)\n", r.Key, r.Permissions)
		}
	}

	// Example 10: Check and throw (for use in authorization middleware)
	fmt.Println("\n=== Example 10: Check and Throw ===")
	err = client.CheckAndThrow(ctx, userEnf, enforcement.Action("admin"), resource)
	if err != nil {
		fmt.Printf("Access denied (as expected): %v\n", err)
	} else {
		fmt.Println("Access granted!")
	}

	fmt.Println("\n=== Examples Complete ===")
}

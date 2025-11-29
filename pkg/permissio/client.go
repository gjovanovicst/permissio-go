// Package permissio provides the main Permissio.io SDK client.
package permissio

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/permissio/permissio-go/pkg/api"
	"github.com/permissio/permissio-go/pkg/config"
	"github.com/permissio/permissio-go/pkg/enforcement"
	"github.com/permissio/permissio-go/pkg/models"
	"go.uber.org/zap"
)

// Api contains all API clients.
type Api struct {
	Users           *api.UsersAPI
	Tenants         *api.TenantsAPI
	Roles           *api.RolesAPI
	Resources       *api.ResourcesAPI
	RoleAssignments *api.RoleAssignmentsAPI
}

// Client is the main Permissio.io SDK client.
type Client struct {
	// Api provides access to all API clients.
	Api *Api

	// config holds the SDK configuration.
	config *config.Config

	// scopeInitialized tracks if scope has been fetched.
	scopeInitialized bool

	// scopeMu protects scope initialization.
	scopeMu sync.Mutex
}

// New creates a new Permissio.io SDK client.
func New(cfg *config.Config) *Client {
	return &Client{
		config: cfg,
		Api: &Api{
			Users:           api.NewUsersAPI(cfg),
			Tenants:         api.NewTenantsAPI(cfg),
			Roles:           api.NewRolesAPI(cfg),
			Resources:       api.NewResourcesAPI(cfg),
			RoleAssignments: api.NewRoleAssignmentsAPI(cfg),
		},
	}
}

// Check performs a permission check.
// Returns true if the user is allowed to perform the action on the resource.
func (c *Client) Check(user enforcement.User, action enforcement.Action, resource enforcement.Resource) (bool, error) {
	return c.CheckWithContext(context.Background(), user, action, resource)
}

// CheckWithContext performs a permission check with context.
func (c *Client) CheckWithContext(ctx context.Context, user enforcement.User, action enforcement.Action, resource enforcement.Resource) (bool, error) {
	response, err := c.CheckWithDetails(ctx, user, action, resource)
	if err != nil {
		return false, err
	}
	return response.Allowed, nil
}

// CheckWithDetails performs a permission check and returns full response details.
// This performs client-side permission checking by:
// 1. Fetching user's role assignments
// 2. Fetching role definitions with permissions
// 3. Checking if any role grants the required permission
func (c *Client) CheckWithDetails(ctx context.Context, user enforcement.User, action enforcement.Action, resource enforcement.Resource) (*models.CheckResponse, error) {
	// Ensure scope is initialized
	if err := c.ensureScope(ctx); err != nil {
		return nil, err
	}

	userKey := user.Key
	resourceType := resource.Type
	requiredPermission := fmt.Sprintf("%s:%s", resourceType, string(action))

	if c.config.Debug && c.config.Logger != nil {
		c.config.Logger.Debug("Permission check",
			zap.String("user", userKey),
			zap.String("action", string(action)),
			zap.String("resource", resourceType),
			zap.String("requiredPermission", requiredPermission))
	}

	// 1. Get user's role assignments (filtered by tenant if provided)
	listParams := &models.RoleAssignmentListParams{
		User: userKey,
	}
	if resource.Tenant != "" {
		listParams.Tenant = resource.Tenant
	}

	assignments, err := c.Api.RoleAssignments.List(ctx, listParams)
	if err != nil {
		if c.config.ThrowOnError {
			return nil, err
		}
		return &models.CheckResponse{
			Allowed: false,
			Reason:  fmt.Sprintf("Error fetching role assignments: %v", err),
		}, nil
	}

	if c.config.Debug && c.config.Logger != nil {
		c.config.Logger.Debug("Role assignments fetched",
			zap.Int("count", len(assignments)))
	}

	if len(assignments) == 0 {
		return &models.CheckResponse{
			Allowed: false,
			Reason:  fmt.Sprintf("User %s has no role assignments", userKey),
		}, nil
	}

	// 2. Get unique role keys from assignments
	roleKeys := make(map[string]struct{})
	for _, assignment := range assignments {
		roleKeys[assignment.Role] = struct{}{}
	}

	if c.config.Debug && c.config.Logger != nil {
		keys := make([]string, 0, len(roleKeys))
		for k := range roleKeys {
			keys = append(keys, k)
		}
		c.config.Logger.Debug("User's role keys", zap.Strings("roles", keys))
	}

	// 3. Fetch all roles and build permission map (with role inheritance)
	rolesResponse, err := c.Api.Roles.List(ctx, &models.RoleListParams{
		ListParams: models.ListParams{PerPage: 100},
	})
	if err != nil {
		if c.config.ThrowOnError {
			return nil, err
		}
		return &models.CheckResponse{
			Allowed: false,
			Reason:  fmt.Sprintf("Error fetching roles: %v", err),
		}, nil
	}

	rolesMap := make(map[string]*models.RoleRead)
	for i := range rolesResponse.Data {
		role := &rolesResponse.Data[i]
		rolesMap[role.Key] = role
	}

	// 4. Check if any assigned role grants the required permission
	var matchedRoles []string
	var matchedPermissions []string

	for roleKey := range roleKeys {
		permissions := c.getRolePermissions(roleKey, rolesMap, make(map[string]struct{}))

		if c.config.Debug && c.config.Logger != nil {
			c.config.Logger.Debug("Role permissions",
				zap.String("role", roleKey),
				zap.Strings("permissions", permissions))
		}

		for _, perm := range permissions {
			if perm == requiredPermission ||
				perm == fmt.Sprintf("%s:*", resourceType) ||
				perm == "*:*" {
				matchedRoles = append(matchedRoles, roleKey)
				matchedPermissions = append(matchedPermissions, requiredPermission)
				break
			}
		}
	}

	allowed := len(matchedRoles) > 0

	if c.config.Debug && c.config.Logger != nil {
		c.config.Logger.Debug("Permission check result",
			zap.Bool("allowed", allowed),
			zap.Strings("matchedRoles", matchedRoles))
	}

	reason := fmt.Sprintf("No role grants permission %s", requiredPermission)
	if allowed {
		reason = fmt.Sprintf("Granted by role(s): %s", strings.Join(matchedRoles, ", "))
	}

	return &models.CheckResponse{
		Allowed: allowed,
		Reason:  reason,
		Debug: &models.CheckDebugInfo{
			MatchedRoles:       matchedRoles,
			MatchedPermissions: matchedPermissions,
		},
	}, nil
}

// getRolePermissions returns all permissions for a role, including inherited ones.
func (c *Client) getRolePermissions(roleKey string, rolesMap map[string]*models.RoleRead, visited map[string]struct{}) []string {
	// Prevent circular inheritance
	if _, ok := visited[roleKey]; ok {
		return nil
	}
	visited[roleKey] = struct{}{}

	role, ok := rolesMap[roleKey]
	if !ok {
		return nil
	}

	permissions := make([]string, len(role.Permissions))
	copy(permissions, role.Permissions)

	// Add inherited permissions from parent roles
	for _, parentRoleKey := range role.Extends {
		parentPermissions := c.getRolePermissions(parentRoleKey, rolesMap, visited)
		permissions = append(permissions, parentPermissions...)
	}

	// Remove duplicates
	seen := make(map[string]struct{})
	unique := permissions[:0]
	for _, perm := range permissions {
		if _, ok := seen[perm]; !ok {
			seen[perm] = struct{}{}
			unique = append(unique, perm)
		}
	}

	return unique
}

// BulkCheck performs multiple permission checks at once.
func (c *Client) BulkCheck(ctx context.Context, checks []models.CheckRequest) (*models.BulkCheckResponse, error) {
	results := make([]models.BulkCheckResult, len(checks))

	for i, check := range checks {
		user := enforcement.User{Key: check.User.(string)}
		action := enforcement.Action(check.Action)

		var resource enforcement.Resource
		switch r := check.Resource.(type) {
		case string:
			resource = enforcement.Resource{Type: r}
		case map[string]interface{}:
			if t, ok := r["type"].(string); ok {
				resource.Type = t
			}
			if k, ok := r["key"].(string); ok {
				resource.Key = k
			}
			if t, ok := r["tenant"].(string); ok {
				resource.Tenant = t
			}
		}

		if check.Tenant != "" {
			resource.Tenant = check.Tenant
		}

		response, err := c.CheckWithDetails(ctx, user, action, resource)
		if err != nil {
			results[i] = models.BulkCheckResult{
				Request:  check,
				Response: models.CheckResponse{Allowed: false, Reason: err.Error()},
			}
		} else {
			results[i] = models.BulkCheckResult{
				Request:  check,
				Response: *response,
			}
		}
	}

	return &models.BulkCheckResponse{Results: results}, nil
}

// GetPermissions returns all permissions for a user.
func (c *Client) GetPermissions(ctx context.Context, request models.GetPermissionsRequest) (*models.GetPermissionsResponse, error) {
	// Ensure scope is initialized
	if err := c.ensureScope(ctx); err != nil {
		return nil, err
	}

	// 1. Get user's role assignments
	listParams := &models.RoleAssignmentListParams{
		User: request.User,
	}
	if request.Tenant != "" {
		listParams.Tenant = request.Tenant
	}
	if request.Resource != "" {
		listParams.Resource = request.Resource
	}

	assignments, err := c.Api.RoleAssignments.List(ctx, listParams)
	if err != nil {
		if c.config.ThrowOnError {
			return nil, err
		}
		return &models.GetPermissionsResponse{Roles: []string{}, Permissions: []string{}}, nil
	}

	if len(assignments) == 0 {
		return &models.GetPermissionsResponse{Roles: []string{}, Permissions: []string{}}, nil
	}

	// 2. Get unique role keys
	roleKeys := make(map[string]struct{})
	for _, assignment := range assignments {
		roleKeys[assignment.Role] = struct{}{}
	}

	// 3. Fetch all roles
	rolesResponse, err := c.Api.Roles.List(ctx, &models.RoleListParams{
		ListParams: models.ListParams{PerPage: 100},
	})
	if err != nil {
		if c.config.ThrowOnError {
			return nil, err
		}
		return &models.GetPermissionsResponse{Roles: []string{}, Permissions: []string{}}, nil
	}

	rolesMap := make(map[string]*models.RoleRead)
	for i := range rolesResponse.Data {
		role := &rolesResponse.Data[i]
		rolesMap[role.Key] = role
	}

	// 4. Collect all permissions from assigned roles
	allPermissions := make(map[string]struct{})
	roles := make([]string, 0, len(roleKeys))

	for roleKey := range roleKeys {
		roles = append(roles, roleKey)
		permissions := c.getRolePermissions(roleKey, rolesMap, make(map[string]struct{}))
		for _, perm := range permissions {
			allPermissions[perm] = struct{}{}
		}
	}

	permissions := make([]string, 0, len(allPermissions))
	for perm := range allPermissions {
		permissions = append(permissions, perm)
	}

	return &models.GetPermissionsResponse{
		Roles:       roles,
		Permissions: permissions,
	}, nil
}

// SyncUser creates or updates a user and optionally assigns roles.
func (c *Client) SyncUser(ctx context.Context, user models.UserCreate, roles []models.RoleAssignmentCreate) (*models.UserRead, error) {
	// Ensure scope is initialized
	if err := c.ensureScope(ctx); err != nil {
		return nil, err
	}

	// Sync user
	result, err := c.Api.Users.SyncUser(ctx, user)
	if err != nil {
		return nil, err
	}

	// Assign roles if provided
	if len(roles) > 0 {
		for _, role := range roles {
			role.User = user.Key
			_, err := c.Api.RoleAssignments.Assign(ctx, &role)
			if err != nil {
				// Log but don't fail if role assignment fails
				if c.config.Debug && c.config.Logger != nil {
					c.config.Logger.Warn("Failed to assign role",
						zap.String("user", user.Key),
						zap.String("role", role.Role),
						zap.Error(err))
				}
			}
		}
	}

	return result, nil
}

// CheckAndThrow performs a permission check and returns an error if not allowed.
func (c *Client) CheckAndThrow(ctx context.Context, user enforcement.User, action enforcement.Action, resource enforcement.Resource) error {
	response, err := c.CheckWithDetails(ctx, user, action, resource)
	if err != nil {
		return err
	}

	if !response.Allowed {
		return api.AccessDeniedError(fmt.Sprintf(
			"Access denied: User %s is not allowed to perform %s on %s",
			user.Key, string(action), resource.Type))
	}

	return nil
}

// GetConfig returns the current configuration.
func (c *Client) GetConfig() *config.Config {
	return c.config
}

// Init initializes the client by fetching the API key scope if not already configured.
// This method should be called before using any API methods if you're relying on
// auto-fetching the projectId and environmentId from the API key.
// Returns an error if the scope cannot be determined.
func (c *Client) Init(ctx context.Context) error {
	return c.ensureScope(ctx)
}

// GetScope returns the project and environment IDs (fetches from API if needed).
func (c *Client) GetScope(ctx context.Context) (projectID, environmentID string, err error) {
	if err := c.ensureScope(ctx); err != nil {
		return "", "", err
	}
	return c.config.ProjectID, c.config.EnvironmentID, nil
}

// ensureScope ensures that projectId and environmentId are available.
func (c *Client) ensureScope(ctx context.Context) error {
	// Fast path: already initialized or has scope from config
	if c.scopeInitialized || c.config.HasScope() {
		return nil
	}

	c.scopeMu.Lock()
	defer c.scopeMu.Unlock()

	// Double-check after acquiring lock
	if c.scopeInitialized || c.config.HasScope() {
		return nil
	}

	// Fetch scope from API
	if err := c.fetchAndSetScope(ctx); err != nil {
		return err
	}

	c.scopeInitialized = true
	return nil
}

// fetchAndSetScope fetches scope from the API key scope endpoint.
func (c *Client) fetchAndSetScope(ctx context.Context) error {
	url := fmt.Sprintf("%s/v1/api-key/scope", c.config.ApiURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.config.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.config.HTTPClient.Do(req)
	if err != nil {
		if !c.config.HasScope() {
			return fmt.Errorf("failed to fetch API key scope: %w. "+
				"Either provide projectId and environmentId in config, "+
				"or ensure the API key has valid scope", err)
		}
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		if !c.config.HasScope() {
			return fmt.Errorf("failed to fetch API key scope: status %d, body: %s. "+
				"Either provide projectId and environmentId in config, "+
				"or ensure the API key has valid scope", resp.StatusCode, string(body))
		}
		return nil
	}

	var scope models.APIKeyScope
	if err := json.NewDecoder(resp.Body).Decode(&scope); err != nil {
		if !c.config.HasScope() {
			return fmt.Errorf("failed to decode API key scope: %w", err)
		}
		return nil
	}

	c.config.UpdateScope(scope.ProjectID, scope.EnvironmentID)

	if c.config.Debug && c.config.Logger != nil {
		c.config.Logger.Debug("Auto-fetched scope",
			zap.String("projectId", scope.ProjectID),
			zap.String("environmentId", scope.EnvironmentID))
	}

	return nil
}

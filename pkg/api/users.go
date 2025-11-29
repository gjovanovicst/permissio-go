package api

import (
	"context"
	"fmt"

	"github.com/permissio/permissio-go/pkg/config"
	"github.com/permissio/permissio-go/pkg/models"
)

// UsersAPI provides methods for managing users.
type UsersAPI struct {
	*BaseClient
}

// NewUsersAPI creates a new UsersAPI.
func NewUsersAPI(cfg *config.Config) *UsersAPI {
	return &UsersAPI{
		BaseClient: NewBaseClient(cfg),
	}
}

// List returns a paginated list of users.
func (a *UsersAPI) List(ctx context.Context, params *models.UserListParams) (*models.UserList, error) {
	url := a.BuildFactsURL("/users")

	if params != nil {
		queryParams := ListParamsToMap(params.Page, params.PerPage, map[string]string{
			"search": params.Search,
			"role":   params.Role,
			"tenant": params.Tenant,
		})
		url = BuildQueryParams(url, queryParams)
	}

	var result models.UserList
	if err := a.BaseClient.Get(ctx, url, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get retrieves a user by key.
func (a *UsersAPI) Get(ctx context.Context, userKey string) (*models.UserRead, error) {
	url := a.BuildFactsURL(fmt.Sprintf("/users/%s", userKey))

	var result models.UserRead
	if err := a.BaseClient.Get(ctx, url, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Create creates a new user.
func (a *UsersAPI) Create(ctx context.Context, user *models.UserCreate) (*models.UserRead, error) {
	url := a.BuildFactsURL("/users")

	var result models.UserRead
	if err := a.Post(ctx, url, user, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Update updates an existing user.
func (a *UsersAPI) Update(ctx context.Context, userKey string, data *models.UserUpdate) (*models.UserRead, error) {
	url := a.BuildFactsURL(fmt.Sprintf("/users/%s", userKey))

	var result models.UserRead
	if err := a.Patch(ctx, url, data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete deletes a user.
func (a *UsersAPI) Delete(ctx context.Context, userKey string) error {
	url := a.BuildFactsURL(fmt.Sprintf("/users/%s", userKey))
	return a.BaseClient.Delete(ctx, url, nil)
}

// SyncUser creates or updates a user (upsert).
// Uses PUT to replace/create the user with the given key.
func (a *UsersAPI) SyncUser(ctx context.Context, user models.UserCreate) (*models.UserRead, error) {
	url := a.BuildFactsURL(fmt.Sprintf("/users/%s", user.Key))

	var result models.UserRead
	if err := a.Put(ctx, url, user, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// AssignRole assigns a role to a user.
func (a *UsersAPI) AssignRole(ctx context.Context, userKey, role, tenant string) (*models.RoleAssignmentRead, error) {
	url := a.BuildFactsURL(fmt.Sprintf("/users/%s/roles", userKey))

	body := map[string]string{
		"role":   role,
		"tenant": tenant,
	}

	var result models.RoleAssignmentRead
	if err := a.Post(ctx, url, body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UnassignRole removes a role from a user.
func (a *UsersAPI) UnassignRole(ctx context.Context, userKey, role, tenant string) error {
	url := a.BuildFactsURL(fmt.Sprintf("/users/%s/roles/%s", userKey, role))
	if tenant != "" {
		url = BuildQueryParams(url, map[string]string{"tenant": tenant})
	}
	return a.BaseClient.Delete(ctx, url, nil)
}

// GetRoles returns the roles assigned to a user.
func (a *UsersAPI) GetRoles(ctx context.Context, userKey string, tenant string) ([]string, error) {
	url := a.BuildFactsURL(fmt.Sprintf("/users/%s/roles", userKey))
	if tenant != "" {
		url = BuildQueryParams(url, map[string]string{"tenant": tenant})
	}

	var result struct {
		Roles []string `json:"roles"`
	}
	if err := a.BaseClient.Get(ctx, url, &result); err != nil {
		return nil, err
	}
	return result.Roles, nil
}

// AddTenant adds a user to a tenant.
func (a *UsersAPI) AddTenant(ctx context.Context, userKey, tenantKey string) error {
	url := a.BuildFactsURL(fmt.Sprintf("/users/%s/tenants", userKey))
	body := map[string]string{"tenant": tenantKey}
	return a.Post(ctx, url, body, nil)
}

// RemoveTenant removes a user from a tenant.
func (a *UsersAPI) RemoveTenant(ctx context.Context, userKey, tenantKey string) error {
	url := a.BuildFactsURL(fmt.Sprintf("/users/%s/tenants/%s", userKey, tenantKey))
	return a.BaseClient.Delete(ctx, url, nil)
}

// GetTenants returns the tenants a user belongs to.
func (a *UsersAPI) GetTenants(ctx context.Context, userKey string) ([]string, error) {
	url := a.BuildFactsURL(fmt.Sprintf("/users/%s/tenants", userKey))

	var result struct {
		Tenants []string `json:"tenants"`
	}
	if err := a.BaseClient.Get(ctx, url, &result); err != nil {
		return nil, err
	}
	return result.Tenants, nil
}

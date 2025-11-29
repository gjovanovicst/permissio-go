package api

import (
	"context"
	"fmt"

	"github.com/permissio/permissio-go/pkg/config"
	"github.com/permissio/permissio-go/pkg/models"
)

// RolesAPI provides methods for managing roles.
type RolesAPI struct {
	*BaseClient
}

// NewRolesAPI creates a new RolesAPI.
func NewRolesAPI(cfg *config.Config) *RolesAPI {
	return &RolesAPI{
		BaseClient: NewBaseClient(cfg),
	}
}

// List returns a paginated list of roles.
func (a *RolesAPI) List(ctx context.Context, params *models.RoleListParams) (*models.RoleList, error) {
	url := a.BuildSchemaURL("/roles")

	if params != nil {
		queryParams := ListParamsToMap(params.Page, params.PerPage, map[string]string{
			"search": params.Search,
		})
		url = BuildQueryParams(url, queryParams)
	}

	var result models.RoleList
	if err := a.BaseClient.Get(ctx, url, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get retrieves a role by key.
func (a *RolesAPI) Get(ctx context.Context, roleKey string) (*models.RoleRead, error) {
	url := a.BuildSchemaURL(fmt.Sprintf("/roles/%s", roleKey))

	var result models.RoleRead
	if err := a.BaseClient.Get(ctx, url, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Create creates a new role.
func (a *RolesAPI) Create(ctx context.Context, role *models.RoleCreate) (*models.RoleRead, error) {
	url := a.BuildSchemaURL("/roles")

	var result models.RoleRead
	if err := a.Post(ctx, url, role, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Update updates an existing role.
func (a *RolesAPI) Update(ctx context.Context, roleKey string, data *models.RoleUpdate) (*models.RoleRead, error) {
	url := a.BuildSchemaURL(fmt.Sprintf("/roles/%s", roleKey))

	var result models.RoleRead
	if err := a.Patch(ctx, url, data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete deletes a role.
func (a *RolesAPI) Delete(ctx context.Context, roleKey string) error {
	url := a.BuildSchemaURL(fmt.Sprintf("/roles/%s", roleKey))
	return a.BaseClient.Delete(ctx, url, nil)
}

// Sync creates or updates a role (upsert).
func (a *RolesAPI) Sync(ctx context.Context, role *models.RoleCreate) (*models.RoleRead, error) {
	url := a.BuildSchemaURL("/roles")

	var result models.RoleRead
	if err := a.Put(ctx, url, role, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPermissions returns the permissions for a role.
func (a *RolesAPI) GetPermissions(ctx context.Context, roleKey string) ([]string, error) {
	url := a.BuildSchemaURL(fmt.Sprintf("/roles/%s/permissions", roleKey))

	var result struct {
		Permissions []string `json:"permissions"`
	}
	if err := a.BaseClient.Get(ctx, url, &result); err != nil {
		return nil, err
	}
	return result.Permissions, nil
}

// AddPermission adds a permission to a role.
func (a *RolesAPI) AddPermission(ctx context.Context, roleKey, permission string) error {
	url := a.BuildSchemaURL(fmt.Sprintf("/roles/%s/permissions", roleKey))
	body := map[string]string{"permission": permission}
	return a.Post(ctx, url, body, nil)
}

// RemovePermission removes a permission from a role.
func (a *RolesAPI) RemovePermission(ctx context.Context, roleKey, permission string) error {
	url := a.BuildSchemaURL(fmt.Sprintf("/roles/%s/permissions/%s", roleKey, permission))
	return a.BaseClient.Delete(ctx, url, nil)
}

// GetExtends returns the roles that this role extends.
func (a *RolesAPI) GetExtends(ctx context.Context, roleKey string) ([]string, error) {
	url := a.BuildSchemaURL(fmt.Sprintf("/roles/%s/extends", roleKey))

	var result struct {
		Extends []string `json:"extends"`
	}
	if err := a.BaseClient.Get(ctx, url, &result); err != nil {
		return nil, err
	}
	return result.Extends, nil
}

// AddExtends adds a parent role to extend from.
func (a *RolesAPI) AddExtends(ctx context.Context, roleKey, parentRoleKey string) error {
	url := a.BuildSchemaURL(fmt.Sprintf("/roles/%s/extends", roleKey))
	body := map[string]string{"role": parentRoleKey}
	return a.Post(ctx, url, body, nil)
}

// RemoveExtends removes a parent role from the extends list.
func (a *RolesAPI) RemoveExtends(ctx context.Context, roleKey, parentRoleKey string) error {
	url := a.BuildSchemaURL(fmt.Sprintf("/roles/%s/extends/%s", roleKey, parentRoleKey))
	return a.BaseClient.Delete(ctx, url, nil)
}

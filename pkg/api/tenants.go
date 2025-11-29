package api

import (
	"context"
	"fmt"

	"github.com/permissio/permissio-go/pkg/config"
	"github.com/permissio/permissio-go/pkg/models"
)

// TenantsAPI provides methods for managing tenants.
type TenantsAPI struct {
	*BaseClient
}

// NewTenantsAPI creates a new TenantsAPI.
func NewTenantsAPI(cfg *config.Config) *TenantsAPI {
	return &TenantsAPI{
		BaseClient: NewBaseClient(cfg),
	}
}

// List returns a paginated list of tenants.
func (a *TenantsAPI) List(ctx context.Context, params *models.TenantListParams) (*models.TenantList, error) {
	url := a.BuildFactsURL("/tenants")

	if params != nil {
		queryParams := ListParamsToMap(params.Page, params.PerPage, map[string]string{
			"search": params.Search,
		})
		url = BuildQueryParams(url, queryParams)
	}

	var result models.TenantList
	if err := a.BaseClient.Get(ctx, url, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get retrieves a tenant by key.
func (a *TenantsAPI) Get(ctx context.Context, tenantKey string) (*models.TenantRead, error) {
	url := a.BuildFactsURL(fmt.Sprintf("/tenants/%s", tenantKey))

	var result models.TenantRead
	if err := a.BaseClient.Get(ctx, url, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Create creates a new tenant.
func (a *TenantsAPI) Create(ctx context.Context, tenant *models.TenantCreate) (*models.TenantRead, error) {
	url := a.BuildFactsURL("/tenants")

	var result models.TenantRead
	if err := a.Post(ctx, url, tenant, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Update updates an existing tenant.
func (a *TenantsAPI) Update(ctx context.Context, tenantKey string, data *models.TenantUpdate) (*models.TenantRead, error) {
	url := a.BuildFactsURL(fmt.Sprintf("/tenants/%s", tenantKey))

	var result models.TenantRead
	if err := a.Patch(ctx, url, data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete deletes a tenant.
func (a *TenantsAPI) Delete(ctx context.Context, tenantKey string) error {
	url := a.BuildFactsURL(fmt.Sprintf("/tenants/%s", tenantKey))
	return a.BaseClient.Delete(ctx, url, nil)
}

// Sync creates or updates a tenant (upsert).
func (a *TenantsAPI) Sync(ctx context.Context, tenant *models.TenantCreate) (*models.TenantRead, error) {
	url := a.BuildFactsURL("/tenants")

	var result models.TenantRead
	if err := a.Put(ctx, url, tenant, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// AddUser adds a user to a tenant.
func (a *TenantsAPI) AddUser(ctx context.Context, tenantKey, userKey string) error {
	url := a.BuildFactsURL(fmt.Sprintf("/tenants/%s/users", tenantKey))
	body := map[string]string{"user": userKey}
	return a.Post(ctx, url, body, nil)
}

// RemoveUser removes a user from a tenant.
func (a *TenantsAPI) RemoveUser(ctx context.Context, tenantKey, userKey string) error {
	url := a.BuildFactsURL(fmt.Sprintf("/tenants/%s/users/%s", tenantKey, userKey))
	return a.BaseClient.Delete(ctx, url, nil)
}

// GetUsers returns the users in a tenant.
func (a *TenantsAPI) GetUsers(ctx context.Context, tenantKey string) ([]string, error) {
	url := a.BuildFactsURL(fmt.Sprintf("/tenants/%s/users", tenantKey))

	var result struct {
		Users []string `json:"users"`
	}
	if err := a.BaseClient.Get(ctx, url, &result); err != nil {
		return nil, err
	}
	return result.Users, nil
}

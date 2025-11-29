package api

import (
	"context"
	"fmt"

	"github.com/permissio/permissio-go/pkg/config"
	"github.com/permissio/permissio-go/pkg/models"
)

// ResourcesAPI provides methods for managing resources.
type ResourcesAPI struct {
	*BaseClient
}

// NewResourcesAPI creates a new ResourcesAPI.
func NewResourcesAPI(cfg *config.Config) *ResourcesAPI {
	return &ResourcesAPI{
		BaseClient: NewBaseClient(cfg),
	}
}

// List returns a paginated list of resources.
func (a *ResourcesAPI) List(ctx context.Context, params *models.ResourceListParams) (*models.ResourceList, error) {
	url := a.BuildSchemaURL("/resources")

	if params != nil {
		queryParams := ListParamsToMap(params.Page, params.PerPage, map[string]string{
			"search": params.Search,
		})
		url = BuildQueryParams(url, queryParams)
	}

	var result models.ResourceList
	if err := a.BaseClient.Get(ctx, url, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get retrieves a resource by key.
func (a *ResourcesAPI) Get(ctx context.Context, resourceKey string) (*models.ResourceRead, error) {
	url := a.BuildSchemaURL(fmt.Sprintf("/resources/%s", resourceKey))

	var result models.ResourceRead
	if err := a.BaseClient.Get(ctx, url, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Create creates a new resource.
func (a *ResourcesAPI) Create(ctx context.Context, resource *models.ResourceCreate) (*models.ResourceRead, error) {
	url := a.BuildSchemaURL("/resources")

	var result models.ResourceRead
	if err := a.Post(ctx, url, resource, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Update updates an existing resource.
func (a *ResourcesAPI) Update(ctx context.Context, resourceKey string, data *models.ResourceUpdate) (*models.ResourceRead, error) {
	url := a.BuildSchemaURL(fmt.Sprintf("/resources/%s", resourceKey))

	var result models.ResourceRead
	if err := a.Patch(ctx, url, data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete deletes a resource.
func (a *ResourcesAPI) Delete(ctx context.Context, resourceKey string) error {
	url := a.BuildSchemaURL(fmt.Sprintf("/resources/%s", resourceKey))
	return a.BaseClient.Delete(ctx, url, nil)
}

// Sync creates or updates a resource (upsert).
func (a *ResourcesAPI) Sync(ctx context.Context, resource *models.ResourceCreate) (*models.ResourceRead, error) {
	url := a.BuildSchemaURL("/resources")

	var result models.ResourceRead
	if err := a.Put(ctx, url, resource, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetActions returns the actions for a resource.
func (a *ResourcesAPI) GetActions(ctx context.Context, resourceKey string) ([]string, error) {
	url := a.BuildSchemaURL(fmt.Sprintf("/resources/%s/actions", resourceKey))

	var result struct {
		Actions []string `json:"actions"`
	}
	if err := a.BaseClient.Get(ctx, url, &result); err != nil {
		return nil, err
	}
	return result.Actions, nil
}

// AddAction adds an action to a resource.
func (a *ResourcesAPI) AddAction(ctx context.Context, resourceKey, action string) error {
	url := a.BuildSchemaURL(fmt.Sprintf("/resources/%s/actions", resourceKey))
	body := map[string]string{"action": action}
	return a.Post(ctx, url, body, nil)
}

// RemoveAction removes an action from a resource.
func (a *ResourcesAPI) RemoveAction(ctx context.Context, resourceKey, action string) error {
	url := a.BuildSchemaURL(fmt.Sprintf("/resources/%s/actions/%s", resourceKey, action))
	return a.BaseClient.Delete(ctx, url, nil)
}

// CreateInstance creates a resource instance.
func (a *ResourcesAPI) CreateInstance(ctx context.Context, resourceKey string, instance *models.ResourceInstanceCreate) (*models.ResourceInstanceRead, error) {
	url := a.BuildFactsURL(fmt.Sprintf("/resources/%s/instances", resourceKey))

	var result models.ResourceInstanceRead
	if err := a.Post(ctx, url, instance, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetInstance retrieves a resource instance.
func (a *ResourcesAPI) GetInstance(ctx context.Context, resourceKey, instanceKey string) (*models.ResourceInstanceRead, error) {
	url := a.BuildFactsURL(fmt.Sprintf("/resources/%s/instances/%s", resourceKey, instanceKey))

	var result models.ResourceInstanceRead
	if err := a.BaseClient.Get(ctx, url, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteInstance deletes a resource instance.
func (a *ResourcesAPI) DeleteInstance(ctx context.Context, resourceKey, instanceKey string) error {
	url := a.BuildFactsURL(fmt.Sprintf("/resources/%s/instances/%s", resourceKey, instanceKey))
	return a.BaseClient.Delete(ctx, url, nil)
}

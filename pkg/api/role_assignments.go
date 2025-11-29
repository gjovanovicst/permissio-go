package api

import (
	"context"
	"fmt"

	"github.com/permissio/permissio-go/pkg/config"
	"github.com/permissio/permissio-go/pkg/models"
)

// RoleAssignmentsAPI provides methods for managing role assignments.
type RoleAssignmentsAPI struct {
	*BaseClient
}

// NewRoleAssignmentsAPI creates a new RoleAssignmentsAPI.
func NewRoleAssignmentsAPI(cfg *config.Config) *RoleAssignmentsAPI {
	return &RoleAssignmentsAPI{
		BaseClient: NewBaseClient(cfg),
	}
}

// List returns a list of role assignments.
func (a *RoleAssignmentsAPI) List(ctx context.Context, params *models.RoleAssignmentListParams) (models.RoleAssignmentList, error) {
	url := a.BuildFactsURL("/role_assignments")

	if params != nil {
		queryParams := ListParamsToMap(params.Page, params.PerPage, map[string]string{
			"user":              params.User,
			"role":              params.Role,
			"tenant":            params.Tenant,
			"resource":          params.Resource,
			"resource_instance": params.ResourceInstance,
		})
		url = BuildQueryParams(url, queryParams)
	}

	var result models.RoleAssignmentList
	if err := a.Get(ctx, url, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// ListByUser returns role assignments for a specific user.
func (a *RoleAssignmentsAPI) ListByUser(ctx context.Context, userKey string, params *models.RoleAssignmentListParams) (models.RoleAssignmentList, error) {
	if params == nil {
		params = &models.RoleAssignmentListParams{}
	}
	params.User = userKey
	return a.List(ctx, params)
}

// ListByTenant returns role assignments for a specific tenant.
func (a *RoleAssignmentsAPI) ListByTenant(ctx context.Context, tenantKey string, params *models.RoleAssignmentListParams) (models.RoleAssignmentList, error) {
	if params == nil {
		params = &models.RoleAssignmentListParams{}
	}
	params.Tenant = tenantKey
	return a.List(ctx, params)
}

// ListByResource returns role assignments for a specific resource.
func (a *RoleAssignmentsAPI) ListByResource(ctx context.Context, resourceType, instanceKey string, params *models.RoleAssignmentListParams) (models.RoleAssignmentList, error) {
	if params == nil {
		params = &models.RoleAssignmentListParams{}
	}
	params.Resource = resourceType
	params.ResourceInstance = instanceKey
	return a.List(ctx, params)
}

// Assign creates a new role assignment.
func (a *RoleAssignmentsAPI) Assign(ctx context.Context, assignment *models.RoleAssignmentCreate) (*models.RoleAssignmentRead, error) {
	url := a.BuildFactsURL("/role_assignments")

	var result models.RoleAssignmentRead
	if err := a.Post(ctx, url, assignment, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Unassign removes a role assignment.
func (a *RoleAssignmentsAPI) Unassign(ctx context.Context, user, role, tenant string) error {
	url := a.BuildFactsURL("/role_assignments")

	body := map[string]string{
		"user":   user,
		"role":   role,
		"tenant": tenant,
	}

	return a.DeleteWithBody(ctx, url, body, nil)
}

// UnassignWithResource removes a role assignment with resource context.
func (a *RoleAssignmentsAPI) UnassignWithResource(ctx context.Context, user, role, tenant, resource, resourceInstance string) error {
	url := a.BuildFactsURL("/role_assignments")

	body := map[string]string{
		"user":              user,
		"role":              role,
		"tenant":            tenant,
		"resource":          resource,
		"resource_instance": resourceInstance,
	}

	return a.DeleteWithBody(ctx, url, body, nil)
}

// BulkAssign creates multiple role assignments at once.
func (a *RoleAssignmentsAPI) BulkAssign(ctx context.Context, assignments []models.RoleAssignmentCreate) (*models.BulkRoleAssignmentResponse, error) {
	url := a.BuildFactsURL("/role_assignments/bulk")

	body := models.BulkRoleAssignmentRequest{
		Assignments: assignments,
	}

	var result models.BulkRoleAssignmentResponse
	if err := a.Post(ctx, url, body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// BulkUnassign removes multiple role assignments at once.
func (a *RoleAssignmentsAPI) BulkUnassign(ctx context.Context, assignments []models.RoleAssignmentCreate) (*models.BulkRoleAssignmentResponse, error) {
	url := a.BuildFactsURL("/role_assignments/bulk")

	body := models.BulkRoleAssignmentRequest{
		Assignments: assignments,
	}

	var result models.BulkRoleAssignmentResponse
	if err := a.DeleteWithBody(ctx, url, body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// HasRole checks if a user has a specific role.
func (a *RoleAssignmentsAPI) HasRole(ctx context.Context, userKey, roleKey string, options *HasRoleOptions) (bool, error) {
	params := &models.RoleAssignmentListParams{
		User: userKey,
		Role: roleKey,
	}

	if options != nil {
		params.Tenant = options.Tenant
		params.Resource = options.Resource
		params.ResourceInstance = options.ResourceInstance
	}

	result, err := a.List(ctx, params)
	if err != nil {
		return false, err
	}

	return len(result) > 0, nil
}

// HasRoleOptions contains optional parameters for HasRole.
type HasRoleOptions struct {
	Tenant           string
	Resource         string
	ResourceInstance string
}

// GetUserRoles returns all roles assigned to a user.
func (a *RoleAssignmentsAPI) GetUserRoles(ctx context.Context, userKey string, options *GetUserRolesOptions) ([]string, error) {
	params := &models.RoleAssignmentListParams{
		User: userKey,
	}

	if options != nil {
		params.Tenant = options.Tenant
		params.Resource = options.Resource
		params.ResourceInstance = options.ResourceInstance
	}

	result, err := a.List(ctx, params)
	if err != nil {
		return nil, err
	}

	roleSet := make(map[string]struct{})
	for _, assignment := range result {
		roleSet[assignment.Role] = struct{}{}
	}

	roles := make([]string, 0, len(roleSet))
	for role := range roleSet {
		roles = append(roles, role)
	}

	return roles, nil
}

// GetUserRolesOptions contains optional parameters for GetUserRoles.
type GetUserRolesOptions struct {
	Tenant           string
	Resource         string
	ResourceInstance string
}

// GetRoleUsers returns all users with a specific role.
func (a *RoleAssignmentsAPI) GetRoleUsers(ctx context.Context, roleKey string, options *GetRoleUsersOptions) ([]string, error) {
	params := &models.RoleAssignmentListParams{
		Role: roleKey,
	}

	if options != nil {
		params.Tenant = options.Tenant
		params.Resource = options.Resource
		params.ResourceInstance = options.ResourceInstance
	}

	result, err := a.List(ctx, params)
	if err != nil {
		return nil, err
	}

	userSet := make(map[string]struct{})
	for _, assignment := range result {
		userSet[assignment.User] = struct{}{}
	}

	users := make([]string, 0, len(userSet))
	for user := range userSet {
		users = append(users, user)
	}

	return users, nil
}

// GetRoleUsersOptions contains optional parameters for GetRoleUsers.
type GetRoleUsersOptions struct {
	Tenant           string
	Resource         string
	ResourceInstance string
}

// ListDetailed returns detailed role assignments with expanded information.
func (a *RoleAssignmentsAPI) ListDetailed(ctx context.Context, params *models.RoleAssignmentListParams) (*models.RoleAssignmentList, error) {
	url := a.BuildFactsURL("/role_assignments/detailed")

	if params != nil {
		queryParams := ListParamsToMap(params.Page, params.PerPage, map[string]string{
			"user":              params.User,
			"role":              params.Role,
			"tenant":            params.Tenant,
			"resource":          params.Resource,
			"resource_instance": params.ResourceInstance,
		})
		url = BuildQueryParams(url, queryParams)
	}

	var result models.RoleAssignmentList
	if err := a.Get(ctx, url, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetByID retrieves a role assignment by ID.
func (a *RoleAssignmentsAPI) GetByID(ctx context.Context, id string) (*models.RoleAssignmentRead, error) {
	url := a.BuildFactsURL(fmt.Sprintf("/role_assignments/%s", id))

	var result models.RoleAssignmentRead
	if err := a.BaseClient.Get(ctx, url, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

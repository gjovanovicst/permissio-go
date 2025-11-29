package models

// RoleAssignmentCreate represents the data for creating a role assignment.
type RoleAssignmentCreate struct {
	User             string `json:"user"`
	Role             string `json:"role"`
	Tenant           string `json:"tenant,omitempty"`
	Resource         string `json:"resource,omitempty"`
	ResourceInstance string `json:"resource_instance,omitempty"`
}

// NewRoleAssignmentCreate creates a new RoleAssignmentCreate.
func NewRoleAssignmentCreate(user, role string) *RoleAssignmentCreate {
	return &RoleAssignmentCreate{
		User: user,
		Role: role,
	}
}

// SetTenant sets the tenant for the role assignment.
func (r *RoleAssignmentCreate) SetTenant(tenant string) *RoleAssignmentCreate {
	r.Tenant = tenant
	return r
}

// SetResource sets the resource for the role assignment.
func (r *RoleAssignmentCreate) SetResource(resource string) *RoleAssignmentCreate {
	r.Resource = resource
	return r
}

// SetResourceInstance sets the resource instance for the role assignment.
func (r *RoleAssignmentCreate) SetResourceInstance(resourceInstance string) *RoleAssignmentCreate {
	r.ResourceInstance = resourceInstance
	return r
}

// RoleAssignmentRead represents a role assignment returned from the API.
type RoleAssignmentRead struct {
	ID               string `json:"id"`
	User             string `json:"user"`
	Role             string `json:"role"`
	Tenant           string `json:"tenant,omitempty"`
	Resource         string `json:"resource,omitempty"`
	ResourceInstance string `json:"resource_instance,omitempty"`
	UserID           string `json:"user_id,omitempty"`
	RoleID           string `json:"role_id,omitempty"`
	TenantID         string `json:"tenant_id,omitempty"`
	OrganizationID   string `json:"organization_id,omitempty"`
	ProjectID        string `json:"project_id,omitempty"`
	EnvironmentID    string `json:"environment_id,omitempty"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at,omitempty"`
}

// RoleAssignmentList represents a list of role assignments.
// Note: The API returns an array directly, not a paginated object.
type RoleAssignmentList []RoleAssignmentRead

// RoleAssignmentListParams represents parameters for listing role assignments.
type RoleAssignmentListParams struct {
	ListParams
	User             string `json:"user,omitempty"`
	Role             string `json:"role,omitempty"`
	Tenant           string `json:"tenant,omitempty"`
	Resource         string `json:"resource,omitempty"`
	ResourceInstance string `json:"resource_instance,omitempty"`
}

// BulkRoleAssignmentRequest represents a bulk role assignment request.
type BulkRoleAssignmentRequest struct {
	Assignments []RoleAssignmentCreate `json:"assignments"`
}

// BulkRoleAssignmentResponse represents a bulk role assignment response.
type BulkRoleAssignmentResponse struct {
	Created int                        `json:"created"`
	Failed  int                        `json:"failed"`
	Errors  []BulkRoleAssignmentError  `json:"errors,omitempty"`
}

// BulkRoleAssignmentError represents an error in a bulk role assignment.
type BulkRoleAssignmentError struct {
	Assignment RoleAssignmentCreate `json:"assignment"`
	Error      string               `json:"error"`
}

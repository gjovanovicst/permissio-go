package models

// RoleCreate represents the data for creating a new role.
type RoleCreate struct {
	Key         string                 `json:"key"`
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Permissions []string               `json:"permissions,omitempty"`
	Extends     []string               `json:"extends,omitempty"`
	Attributes  map[string]interface{} `json:"attributes,omitempty"`
}

// NewRoleCreate creates a new RoleCreate with the given key.
func NewRoleCreate(key string) *RoleCreate {
	return &RoleCreate{Key: key}
}

// SetName sets the name for the role.
func (r *RoleCreate) SetName(name string) *RoleCreate {
	r.Name = name
	return r
}

// SetDescription sets the description for the role.
func (r *RoleCreate) SetDescription(description string) *RoleCreate {
	r.Description = description
	return r
}

// SetPermissions sets the permissions for the role.
func (r *RoleCreate) SetPermissions(permissions []string) *RoleCreate {
	r.Permissions = permissions
	return r
}

// AddPermission adds a permission to the role.
func (r *RoleCreate) AddPermission(permission string) *RoleCreate {
	r.Permissions = append(r.Permissions, permission)
	return r
}

// SetExtends sets the roles this role extends.
func (r *RoleCreate) SetExtends(extends []string) *RoleCreate {
	r.Extends = extends
	return r
}

// SetAttributes sets custom attributes for the role.
func (r *RoleCreate) SetAttributes(attributes map[string]interface{}) *RoleCreate {
	r.Attributes = attributes
	return r
}

// RoleUpdate represents the data for updating a role.
type RoleUpdate struct {
	Name        *string                `json:"name,omitempty"`
	Description *string                `json:"description,omitempty"`
	Permissions []string               `json:"permissions,omitempty"`
	Extends     []string               `json:"extends,omitempty"`
	Attributes  map[string]interface{} `json:"attributes,omitempty"`
}

// RoleRead represents a role returned from the API.
type RoleRead struct {
	ID          string                 `json:"id"`
	Key         string                 `json:"key"`
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Permissions []string               `json:"permissions,omitempty"`
	Extends     []string               `json:"extends,omitempty"`
	Attributes  map[string]interface{} `json:"attributes,omitempty"`
	CreatedAt   string                 `json:"created_at"`
	UpdatedAt   string                 `json:"updated_at"`
}

// RoleList represents a paginated list of roles.
type RoleList struct {
	Data []RoleRead `json:"data"`
	PaginatedResponse
}

// RoleListParams represents parameters for listing roles.
type RoleListParams struct {
	ListParams
	Search string `json:"search,omitempty"`
}

package models

// TenantCreate represents the data for creating a new tenant.
type TenantCreate struct {
	Key         string                 `json:"key"`
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Attributes  map[string]interface{} `json:"attributes,omitempty"`
}

// NewTenantCreate creates a new TenantCreate with the given key.
func NewTenantCreate(key string) *TenantCreate {
	return &TenantCreate{Key: key}
}

// SetName sets the name for the tenant.
func (t *TenantCreate) SetName(name string) *TenantCreate {
	t.Name = name
	return t
}

// SetDescription sets the description for the tenant.
func (t *TenantCreate) SetDescription(description string) *TenantCreate {
	t.Description = description
	return t
}

// SetAttributes sets custom attributes for the tenant.
func (t *TenantCreate) SetAttributes(attributes map[string]interface{}) *TenantCreate {
	t.Attributes = attributes
	return t
}

// TenantUpdate represents the data for updating a tenant.
type TenantUpdate struct {
	Name        *string                `json:"name,omitempty"`
	Description *string                `json:"description,omitempty"`
	Attributes  map[string]interface{} `json:"attributes,omitempty"`
}

// TenantRead represents a tenant returned from the API.
type TenantRead struct {
	ID          string                 `json:"id"`
	Key         string                 `json:"key"`
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Attributes  map[string]interface{} `json:"attributes,omitempty"`
	CreatedAt   string                 `json:"created_at"`
	UpdatedAt   string                 `json:"updated_at"`
}

// TenantList represents a paginated list of tenants.
type TenantList struct {
	Data []TenantRead `json:"data"`
	PaginatedResponse
}

// TenantListParams represents parameters for listing tenants.
type TenantListParams struct {
	ListParams
	Search string `json:"search,omitempty"`
}

package models

// ResourceCreate represents the data for creating a new resource type.
type ResourceCreate struct {
	Key         string                 `json:"key"`
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Actions     []string               `json:"actions,omitempty"`
	Attributes  map[string]interface{} `json:"attributes,omitempty"`
}

// NewResourceCreate creates a new ResourceCreate with the given key.
func NewResourceCreate(key string) *ResourceCreate {
	return &ResourceCreate{Key: key}
}

// SetName sets the name for the resource.
func (r *ResourceCreate) SetName(name string) *ResourceCreate {
	r.Name = name
	return r
}

// SetDescription sets the description for the resource.
func (r *ResourceCreate) SetDescription(description string) *ResourceCreate {
	r.Description = description
	return r
}

// SetActions sets the actions for the resource.
func (r *ResourceCreate) SetActions(actions []string) *ResourceCreate {
	r.Actions = actions
	return r
}

// AddAction adds an action to the resource.
func (r *ResourceCreate) AddAction(action string) *ResourceCreate {
	r.Actions = append(r.Actions, action)
	return r
}

// SetAttributes sets custom attributes for the resource.
func (r *ResourceCreate) SetAttributes(attributes map[string]interface{}) *ResourceCreate {
	r.Attributes = attributes
	return r
}

// ResourceUpdate represents the data for updating a resource.
type ResourceUpdate struct {
	Name        *string                `json:"name,omitempty"`
	Description *string                `json:"description,omitempty"`
	Actions     []string               `json:"actions,omitempty"`
	Attributes  map[string]interface{} `json:"attributes,omitempty"`
}

// ResourceRead represents a resource returned from the API.
type ResourceRead struct {
	ID          string                 `json:"id"`
	Key         string                 `json:"key"`
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Actions     []string               `json:"actions,omitempty"`
	Attributes  map[string]interface{} `json:"attributes,omitempty"`
	CreatedAt   string                 `json:"created_at"`
	UpdatedAt   string                 `json:"updated_at"`
}

// ResourceList represents a paginated list of resources.
type ResourceList struct {
	Data []ResourceRead `json:"data"`
	PaginatedResponse
}

// ResourceListParams represents parameters for listing resources.
type ResourceListParams struct {
	ListParams
	Search string `json:"search,omitempty"`
}

// ResourceInstanceCreate represents the data for creating a resource instance.
type ResourceInstanceCreate struct {
	Key          string                 `json:"key"`
	ResourceType string                 `json:"resource_type"`
	Tenant       string                 `json:"tenant,omitempty"`
	Attributes   map[string]interface{} `json:"attributes,omitempty"`
}

// ResourceInstanceRead represents a resource instance returned from the API.
type ResourceInstanceRead struct {
	ID           string                 `json:"id"`
	Key          string                 `json:"key"`
	ResourceType string                 `json:"resource_type"`
	Tenant       string                 `json:"tenant,omitempty"`
	Attributes   map[string]interface{} `json:"attributes,omitempty"`
	CreatedAt    string                 `json:"created_at"`
	UpdatedAt    string                 `json:"updated_at"`
}

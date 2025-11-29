// Package enforcement provides builder types for constructing permission check requests.
package enforcement

// Action represents an action to check permission for.
type Action string

// User represents a user in a permission check.
type User struct {
	Key        string                 `json:"key"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}

// userBuilder provides a fluent interface for building User.
type userBuilder struct {
	user User
}

// UserBuilder creates a new userBuilder with the given user key.
func UserBuilder(key string) *userBuilder {
	return &userBuilder{
		user: User{
			Key:        key,
			Attributes: make(map[string]interface{}),
		},
	}
}

// WithAttribute adds an attribute to the user.
func (b *userBuilder) WithAttribute(key string, value interface{}) *userBuilder {
	b.user.Attributes[key] = value
	return b
}

// WithAttributes sets multiple attributes for the user.
func (b *userBuilder) WithAttributes(attributes map[string]interface{}) *userBuilder {
	for k, v := range attributes {
		b.user.Attributes[k] = v
	}
	return b
}

// Build returns the built User.
func (b *userBuilder) Build() User {
	return b.user
}

// Resource represents a resource in a permission check.
type Resource struct {
	Type       string                 `json:"type"`
	Key        string                 `json:"key,omitempty"`
	Tenant     string                 `json:"tenant,omitempty"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}

// resourceBuilder provides a fluent interface for building Resource.
type resourceBuilder struct {
	resource Resource
}

// ResourceBuilder creates a new resourceBuilder with the given resource type.
func ResourceBuilder(resourceType string) *resourceBuilder {
	return &resourceBuilder{
		resource: Resource{
			Type:       resourceType,
			Attributes: make(map[string]interface{}),
		},
	}
}

// WithKey sets the resource instance key.
func (b *resourceBuilder) WithKey(key string) *resourceBuilder {
	b.resource.Key = key
	return b
}

// WithTenant sets the tenant for the resource.
func (b *resourceBuilder) WithTenant(tenant string) *resourceBuilder {
	b.resource.Tenant = tenant
	return b
}

// WithAttribute adds an attribute to the resource.
func (b *resourceBuilder) WithAttribute(key string, value interface{}) *resourceBuilder {
	b.resource.Attributes[key] = value
	return b
}

// WithAttributes sets multiple attributes for the resource.
func (b *resourceBuilder) WithAttributes(attributes map[string]interface{}) *resourceBuilder {
	for k, v := range attributes {
		b.resource.Attributes[k] = v
	}
	return b
}

// Build returns the built Resource.
func (b *resourceBuilder) Build() Resource {
	return b.resource
}

// Context represents additional context for a permission check.
type Context struct {
	data map[string]interface{}
}

// contextBuilder provides a fluent interface for building Context.
type contextBuilder struct {
	context Context
}

// ContextBuilder creates a new contextBuilder.
func ContextBuilder() *contextBuilder {
	return &contextBuilder{
		context: Context{
			data: make(map[string]interface{}),
		},
	}
}

// With adds a key-value pair to the context.
func (b *contextBuilder) With(key string, value interface{}) *contextBuilder {
	b.context.data[key] = value
	return b
}

// WithData sets multiple key-value pairs in the context.
func (b *contextBuilder) WithData(data map[string]interface{}) *contextBuilder {
	for k, v := range data {
		b.context.data[k] = v
	}
	return b
}

// Build returns the built Context.
func (b *contextBuilder) Build() Context {
	return b.context
}

// Data returns the context data as a map.
func (c Context) Data() map[string]interface{} {
	return c.data
}

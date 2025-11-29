package models

// UserCreate represents the data for creating a new user.
type UserCreate struct {
	Key        string                 `json:"key"`
	Email      string                 `json:"email,omitempty"`
	FirstName  string                 `json:"first_name,omitempty"`
	LastName   string                 `json:"last_name,omitempty"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}

// NewUserCreate creates a new UserCreate with the given key.
func NewUserCreate(key string) *UserCreate {
	return &UserCreate{Key: key}
}

// SetEmail sets the email for the user.
func (u *UserCreate) SetEmail(email string) *UserCreate {
	u.Email = email
	return u
}

// SetFirstName sets the first name for the user.
func (u *UserCreate) SetFirstName(firstName string) *UserCreate {
	u.FirstName = firstName
	return u
}

// SetLastName sets the last name for the user.
func (u *UserCreate) SetLastName(lastName string) *UserCreate {
	u.LastName = lastName
	return u
}

// SetAttributes sets custom attributes for the user.
func (u *UserCreate) SetAttributes(attributes map[string]interface{}) *UserCreate {
	u.Attributes = attributes
	return u
}

// UserUpdate represents the data for updating a user.
type UserUpdate struct {
	Email      *string                `json:"email,omitempty"`
	FirstName  *string                `json:"first_name,omitempty"`
	LastName   *string                `json:"last_name,omitempty"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}

// UserRead represents a user returned from the API.
type UserRead struct {
	ID         string                 `json:"id"`
	Key        string                 `json:"key"`
	Email      string                 `json:"email,omitempty"`
	FirstName  string                 `json:"first_name,omitempty"`
	LastName   string                 `json:"last_name,omitempty"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
	CreatedAt  string                 `json:"created_at"`
	UpdatedAt  string                 `json:"updated_at"`
}

// UserList represents a paginated list of users.
type UserList struct {
	Data []UserRead `json:"data"`
	PaginatedResponse
}

// UserListParams represents parameters for listing users.
type UserListParams struct {
	ListParams
	Search string `json:"search,omitempty"`
	Role   string `json:"role,omitempty"`
	Tenant string `json:"tenant,omitempty"`
}

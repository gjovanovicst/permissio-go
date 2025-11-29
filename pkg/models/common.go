// Package models provides data types for the Permissio.io SDK.
package models

// ListParams represents common pagination parameters for list operations.
type ListParams struct {
	Page    int `json:"page,omitempty"`
	PerPage int `json:"perPage,omitempty"`
}

// PaginatedResponse represents paginated API response metadata.
type PaginatedResponse struct {
	Page       int `json:"page"`
	PerPage    int `json:"perPage"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}

// APIKeyScope represents the scope information from an API key.
type APIKeyScope struct {
	ProjectID     string `json:"project_id"`
	EnvironmentID string `json:"environment_id"`
}

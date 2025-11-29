package models

// CheckRequest represents a permission check request.
type CheckRequest struct {
	User     interface{}            `json:"user"`
	Action   string                 `json:"action"`
	Resource interface{}            `json:"resource"`
	Tenant   string                 `json:"tenant,omitempty"`
	Context  map[string]interface{} `json:"context,omitempty"`
}

// CheckResponse represents a permission check response.
type CheckResponse struct {
	Allowed bool              `json:"allowed"`
	Reason  string            `json:"reason,omitempty"`
	Debug   *CheckDebugInfo   `json:"debug,omitempty"`
}

// CheckDebugInfo contains debug information from a permission check.
type CheckDebugInfo struct {
	MatchedRoles       []string `json:"matchedRoles,omitempty"`
	MatchedPermissions []string `json:"matchedPermissions,omitempty"`
	EvaluationTime     int64    `json:"evaluationTime,omitempty"`
}

// BulkCheckRequest represents a bulk permission check request.
type BulkCheckRequest struct {
	Checks []CheckRequest `json:"checks"`
}

// BulkCheckResult represents a single result in a bulk check response.
type BulkCheckResult struct {
	Request  CheckRequest  `json:"request"`
	Response CheckResponse `json:"response"`
}

// BulkCheckResponse represents a bulk permission check response.
type BulkCheckResponse struct {
	Results []BulkCheckResult `json:"results"`
}

// GetPermissionsRequest represents a request to get all permissions for a user.
type GetPermissionsRequest struct {
	User     string `json:"user"`
	Tenant   string `json:"tenant,omitempty"`
	Resource string `json:"resource,omitempty"`
}

// GetPermissionsResponse represents the response containing user permissions.
type GetPermissionsResponse struct {
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}

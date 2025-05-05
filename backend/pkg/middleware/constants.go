package middleware

import "github.com/google/uuid"

// DefaultOrganizationID is the default organization ID used for testing
// In a real application, this would be determined based on the request
var DefaultOrganizationID = uuid.MustParse("00000000-0000-0000-0000-000000000001")

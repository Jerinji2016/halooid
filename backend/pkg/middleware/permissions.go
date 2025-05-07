package middleware

// Permission constants for the RBAC system
const (
	// User permissions
	PermissionReadUsers   = "read:users"
	PermissionWriteUsers  = "write:users"
	PermissionDeleteUsers = "delete:users"

	// Organization permissions
	PermissionReadOrganizations   = "read:organizations"
	PermissionWriteOrganizations  = "write:organizations"
	PermissionDeleteOrganizations = "delete:organizations"

	// Role permissions
	PermissionReadRoles   = "read:roles"
	PermissionWriteRoles  = "write:roles"
	PermissionDeleteRoles = "delete:roles"
	PermissionAssignRoles = "assign:roles"

	// Permission permissions
	PermissionReadPermissions   = "read:permissions"
	PermissionWritePermissions  = "write:permissions"
	PermissionDeletePermissions = "delete:permissions"
	PermissionAssignPermissions = "assign:permissions"

	// Admin permissions
	PermissionAdminAccess = "admin:access"

	// User-specific permissions
	PermissionUserRead   = "user:read"
	PermissionUserWrite  = "user:write"
	PermissionUserDelete = "user:delete"

	// Product-specific permissions - Taskodex
	PermissionTaskodexAccess = "taskodex:access"
	PermissionTaskRead       = "task:read"
	PermissionTaskWrite      = "task:write"
	PermissionTaskDelete     = "task:delete"
	PermissionProjectRead    = "project:read"
	PermissionProjectWrite   = "project:write"
	PermissionProjectDelete  = "project:delete"

	// Product-specific permissions - Qultrix
	PermissionQultrixAccess      = "qultrix:access"
	PermissionEmployeeRead       = "employee:read"
	PermissionEmployeeWrite      = "employee:write"
	PermissionEmployeeDelete     = "employee:delete"
	PermissionTimeOffRead        = "timeoff:read"
	PermissionTimeOffWrite       = "timeoff:write"
	PermissionTimeOffDelete      = "timeoff:delete"
	PermissionPerformanceRead    = "performance:read"
	PermissionPerformanceWrite   = "performance:write"
	PermissionPerformanceDelete  = "performance:delete"
	PermissionRecruitmentRead    = "recruitment:read"
	PermissionRecruitmentWrite   = "recruitment:write"
	PermissionRecruitmentDelete  = "recruitment:delete"

	// Product-specific permissions - AdminHub
	PermissionAdminHubAccess     = "adminhub:access"
	PermissionMonitoringRead     = "monitoring:read"
	PermissionMonitoringWrite    = "monitoring:write"
	PermissionSecurityRead       = "security:read"
	PermissionSecurityWrite      = "security:write"

	// Product-specific permissions - CustomerConnect
	PermissionCustomerConnectAccess = "customerconnect:access"
	PermissionContactRead           = "contact:read"
	PermissionContactWrite          = "contact:write"
	PermissionContactDelete         = "contact:delete"
	PermissionLeadRead              = "lead:read"
	PermissionLeadWrite             = "lead:write"
	PermissionLeadDelete            = "lead:delete"
	PermissionOpportunityRead       = "opportunity:read"
	PermissionOpportunityWrite      = "opportunity:write"
	PermissionOpportunityDelete     = "opportunity:delete"
	PermissionCaseRead              = "case:read"
	PermissionCaseWrite             = "case:write"
	PermissionCaseDelete            = "case:delete"

	// Product-specific permissions - Invantray
	PermissionInvantrayAccess    = "invantray:access"
	PermissionInventoryRead      = "inventory:read"
	PermissionInventoryWrite     = "inventory:write"
	PermissionInventoryDelete    = "inventory:delete"
	PermissionWarehouseRead      = "warehouse:read"
	PermissionWarehouseWrite     = "warehouse:write"
	PermissionWarehouseDelete    = "warehouse:delete"
	PermissionAssetRead          = "asset:read"
	PermissionAssetWrite         = "asset:write"
	PermissionAssetDelete        = "asset:delete"
	PermissionProcurementRead    = "procurement:read"
	PermissionProcurementWrite   = "procurement:write"
	PermissionProcurementDelete  = "procurement:delete"
)

// RoleAdmin is the name of the admin role
const RoleAdmin = "admin"

// RoleUser is the name of the user role
const RoleUser = "user"

// GetProductPermissions returns all permissions for a specific product
func GetProductPermissions(product string) []string {
	switch product {
	case "taskodex":
		return []string{
			PermissionTaskodexAccess,
			PermissionTaskRead,
			PermissionTaskWrite,
			PermissionTaskDelete,
			PermissionProjectRead,
			PermissionProjectWrite,
			PermissionProjectDelete,
		}
	case "qultrix":
		return []string{
			PermissionQultrixAccess,
			PermissionEmployeeRead,
			PermissionEmployeeWrite,
			PermissionEmployeeDelete,
			PermissionTimeOffRead,
			PermissionTimeOffWrite,
			PermissionTimeOffDelete,
			PermissionPerformanceRead,
			PermissionPerformanceWrite,
			PermissionPerformanceDelete,
			PermissionRecruitmentRead,
			PermissionRecruitmentWrite,
			PermissionRecruitmentDelete,
		}
	case "adminhub":
		return []string{
			PermissionAdminHubAccess,
			PermissionMonitoringRead,
			PermissionMonitoringWrite,
			PermissionSecurityRead,
			PermissionSecurityWrite,
			PermissionReadUsers,
			PermissionWriteUsers,
			PermissionDeleteUsers,
			PermissionReadOrganizations,
			PermissionWriteOrganizations,
			PermissionDeleteOrganizations,
			PermissionReadRoles,
			PermissionWriteRoles,
			PermissionDeleteRoles,
			PermissionAssignRoles,
			PermissionReadPermissions,
			PermissionWritePermissions,
			PermissionDeletePermissions,
			PermissionAssignPermissions,
		}
	case "customerconnect":
		return []string{
			PermissionCustomerConnectAccess,
			PermissionContactRead,
			PermissionContactWrite,
			PermissionContactDelete,
			PermissionLeadRead,
			PermissionLeadWrite,
			PermissionLeadDelete,
			PermissionOpportunityRead,
			PermissionOpportunityWrite,
			PermissionOpportunityDelete,
			PermissionCaseRead,
			PermissionCaseWrite,
			PermissionCaseDelete,
		}
	case "invantray":
		return []string{
			PermissionInvantrayAccess,
			PermissionInventoryRead,
			PermissionInventoryWrite,
			PermissionInventoryDelete,
			PermissionWarehouseRead,
			PermissionWarehouseWrite,
			PermissionWarehouseDelete,
			PermissionAssetRead,
			PermissionAssetWrite,
			PermissionAssetDelete,
			PermissionProcurementRead,
			PermissionProcurementWrite,
			PermissionProcurementDelete,
		}
	default:
		return []string{}
	}
}

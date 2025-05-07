# Role-Based Access Control (RBAC)

The Halooid platform implements a comprehensive Role-Based Access Control (RBAC) system to manage user permissions across all products. This document describes the architecture and usage of the RBAC system.

## Overview

The RBAC system is based on three main entities:

1. **Users**: Individuals who access the system
2. **Roles**: Named collections of permissions
3. **Permissions**: Fine-grained access controls for specific actions

Users are assigned roles, and roles are assigned permissions. This creates a flexible system where permissions can be managed at scale by grouping them into roles.

## Database Schema

The RBAC system uses the following database tables:

- `roles`: Defines roles in the system
- `permissions`: Defines permissions in the system
- `role_permissions`: Maps permissions to roles
- `user_roles`: Maps roles to users within an organization context

## Core Components

### RBAC Service

The RBAC service (`backend/internal/rbac/service.go`) provides the core functionality:

- Creating, retrieving, updating, and deleting roles
- Creating, retrieving, updating, and deleting permissions
- Assigning permissions to roles
- Assigning roles to users
- Checking if a user has a specific permission

### RBAC Middleware

The RBAC middleware (`backend/pkg/middleware/rbac.go`) provides HTTP handlers for permission checking:

- `RequirePermission`: Requires a specific permission
- `RequireAnyPermission`: Requires any of a set of permissions
- `RequireAllPermissions`: Requires all of a set of permissions

### JWT Integration

The RBAC system is integrated with the JWT authentication system:

1. When a user logs in, their roles and permissions are included in the JWT token
2. The RBAC middleware first checks the token for permissions before making database calls
3. If permissions are updated after a token is issued, the system falls back to database checks

## Permission Naming Convention

Permissions follow a `resource:action` naming convention:

- `read:users` - Permission to read user data
- `write:users` - Permission to create or update user data
- `delete:users` - Permission to delete user data

Product-specific permissions include a product prefix:

- `taskodex:access` - Permission to access the Taskodex product
- `qultrix:access` - Permission to access the Qultrix product

## Standard Roles

The system includes several standard roles:

- `admin`: Has all permissions
- `user`: Basic user with limited permissions
- Product-specific roles (e.g., `taskodex_user`, `qultrix_manager`)

## Usage Examples

### Protecting an API Endpoint

```go
// Require a specific permission
router.GET("/users", userHandler.GetUsers, rbacMiddleware.RequirePermission("read:users"))

// Require any of a set of permissions
router.PUT("/users/:id", userHandler.UpdateUser, rbacMiddleware.RequireAnyPermission("write:users", "admin:access"))

// Require all of a set of permissions
router.DELETE("/users/:id", userHandler.DeleteUser, rbacMiddleware.RequireAllPermissions("delete:users", "admin:access"))
```

### Checking Permissions in Business Logic

```go
// Check if a user has a permission
hasPermission, err := rbacService.HasPermission(ctx, userID, orgID, "read:users")
if err != nil {
    // Handle error
}

if !hasPermission {
    // Handle permission denied
}
```

## Best Practices

1. **Least Privilege**: Assign the minimum permissions necessary
2. **Role Hierarchy**: Create a hierarchy of roles for easier management
3. **Regular Audits**: Regularly review role assignments
4. **Permission Granularity**: Create fine-grained permissions for better control
5. **Context-Aware Permissions**: Consider the organization context when checking permissions

## Future Enhancements

1. **Resource-Level Permissions**: Add support for permissions on specific resources
2. **Permission Inheritance**: Implement a permission inheritance system
3. **Dynamic Permissions**: Support for dynamically calculated permissions
4. **Delegation**: Allow users to delegate permissions to others
5. **Audit Logging**: Log all permission checks for audit purposes

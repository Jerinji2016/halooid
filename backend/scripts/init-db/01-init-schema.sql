-- Initialize database schema for Halooid platform

-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create organizations table
CREATE TABLE IF NOT EXISTS organizations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create organization_users table (many-to-many relationship)
CREATE TABLE IF NOT EXISTS organization_users (
    organization_id UUID NOT NULL REFERENCES organizations(id),
    user_id UUID NOT NULL REFERENCES users(id),
    PRIMARY KEY (organization_id, user_id)
);

-- Create roles table
CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create permissions table
CREATE TABLE IF NOT EXISTS permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create role_permissions table (many-to-many relationship)
CREATE TABLE IF NOT EXISTS role_permissions (
    role_id UUID NOT NULL REFERENCES roles(id),
    permission_id UUID NOT NULL REFERENCES permissions(id),
    PRIMARY KEY (role_id, permission_id)
);

-- Create user_roles table (many-to-many relationship)
CREATE TABLE IF NOT EXISTS user_roles (
    user_id UUID NOT NULL REFERENCES users(id),
    role_id UUID NOT NULL REFERENCES roles(id),
    organization_id UUID NOT NULL REFERENCES organizations(id),
    PRIMARY KEY (user_id, role_id, organization_id)
);

-- Create indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_organizations_name ON organizations(name);
CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX idx_user_roles_organization_id ON user_roles(organization_id);
CREATE INDEX idx_role_permissions_role_id ON role_permissions(role_id);

-- Insert default roles
INSERT INTO roles (name, description) VALUES
    ('admin', 'Administrator with full access'),
    ('user', 'Regular user with limited access')
ON CONFLICT DO NOTHING;

-- Insert default permissions
INSERT INTO permissions (name, description) VALUES
    -- User permissions
    ('read:users', 'Can read user information'),
    ('write:users', 'Can create and update user information'),
    ('delete:users', 'Can delete users'),

    -- Organization permissions
    ('read:organizations', 'Can read organization information'),
    ('write:organizations', 'Can create and update organization information'),
    ('delete:organizations', 'Can delete organizations'),

    -- Role permissions
    ('read:roles', 'Can read role information'),
    ('write:roles', 'Can create and update roles'),
    ('delete:roles', 'Can delete roles'),
    ('assign:roles', 'Can assign roles to users'),

    -- Permission permissions
    ('read:permissions', 'Can read permission information'),
    ('write:permissions', 'Can create and update permissions'),
    ('delete:permissions', 'Can delete permissions'),
    ('assign:permissions', 'Can assign permissions to roles'),

    -- Admin permissions
    ('admin:access', 'Can access admin features'),

    -- User-specific permissions
    ('user:read', 'Can read user data'),
    ('user:write', 'Can write user data'),
    ('user:delete', 'Can delete user data')
ON CONFLICT DO NOTHING;

-- Assign permissions to roles

-- Assign all permissions to admin role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'admin'
ON CONFLICT DO NOTHING;

-- Assign limited permissions to user role
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'user' AND p.name IN (
    'read:users',
    'read:organizations',
    'user:read'
)
ON CONFLICT DO NOTHING;

-- Insert default organization
INSERT INTO organizations (id, name, description)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'Default Organization',
    'Default organization for testing'
)
ON CONFLICT DO NOTHING;

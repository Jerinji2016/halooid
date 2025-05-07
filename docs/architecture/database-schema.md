# Database Schema

This document describes the database schema for the Halooid platform.

## Overview

The Halooid platform uses PostgreSQL as its primary database. The schema is organized into several parts:

1. **Core Schema**: Contains shared tables used across all products
2. **Product-Specific Schemas**: Each product has its own schema with tables specific to that product

## Core Schema

The core schema contains tables that are shared across all products.

### Users

The `users` table stores information about users of the platform.

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| email | VARCHAR(255) | User's email address (unique) |
| password_hash | VARCHAR(255) | Hashed password |
| first_name | VARCHAR(100) | User's first name |
| last_name | VARCHAR(100) | User's last name |
| is_active | BOOLEAN | Whether the user is active |
| created_at | TIMESTAMP | When the user was created |
| updated_at | TIMESTAMP | When the user was last updated |

### Organizations

The `organizations` table stores information about organizations.

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| name | VARCHAR(255) | Organization name |
| description | TEXT | Organization description |
| is_active | BOOLEAN | Whether the organization is active |
| created_at | TIMESTAMP | When the organization was created |
| updated_at | TIMESTAMP | When the organization was last updated |

### Organization Users

The `organization_users` table maps users to organizations.

| Column | Type | Description |
|--------|------|-------------|
| organization_id | UUID | Foreign key to organizations |
| user_id | UUID | Foreign key to users |

### Roles

The `roles` table defines roles that can be assigned to users.

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| name | VARCHAR(100) | Role name |
| description | TEXT | Role description |
| created_at | TIMESTAMP | When the role was created |
| updated_at | TIMESTAMP | When the role was last updated |

### Permissions

The `permissions` table defines permissions that can be assigned to roles.

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| name | VARCHAR(100) | Permission name |
| description | TEXT | Permission description |
| created_at | TIMESTAMP | When the permission was created |
| updated_at | TIMESTAMP | When the permission was last updated |

### Role Permissions

The `role_permissions` table maps permissions to roles.

| Column | Type | Description |
|--------|------|-------------|
| role_id | UUID | Foreign key to roles |
| permission_id | UUID | Foreign key to permissions |

### User Roles

The `user_roles` table maps roles to users within an organization context.

| Column | Type | Description |
|--------|------|-------------|
| user_id | UUID | Foreign key to users |
| role_id | UUID | Foreign key to roles |
| organization_id | UUID | Foreign key to organizations |

## Product-Specific Schemas

### Taskodex Schema

The Taskodex schema contains tables related to task and project management.

#### Projects

The `taskodex.projects` table stores information about projects.

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| organization_id | UUID | Foreign key to organizations |
| name | VARCHAR(255) | Project name |
| description | TEXT | Project description |
| status | VARCHAR(50) | Project status |
| start_date | DATE | Project start date |
| end_date | DATE | Project end date |
| created_by | UUID | Foreign key to users |
| created_at | TIMESTAMP | When the project was created |
| updated_at | TIMESTAMP | When the project was last updated |

#### Tasks

The `taskodex.tasks` table stores information about tasks.

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| project_id | UUID | Foreign key to projects |
| title | VARCHAR(255) | Task title |
| description | TEXT | Task description |
| status | VARCHAR(50) | Task status |
| priority | VARCHAR(50) | Task priority |
| due_date | TIMESTAMP | Task due date |
| created_by | UUID | Foreign key to users |
| assigned_to | UUID | Foreign key to users |
| created_at | TIMESTAMP | When the task was created |
| updated_at | TIMESTAMP | When the task was last updated |

#### Task Comments

The `taskodex.task_comments` table stores comments on tasks.

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| task_id | UUID | Foreign key to tasks |
| user_id | UUID | Foreign key to users |
| content | TEXT | Comment content |
| created_at | TIMESTAMP | When the comment was created |
| updated_at | TIMESTAMP | When the comment was last updated |

### Qultrix Schema

The Qultrix schema contains tables related to HR and employee management.

#### Employees

The `qultrix.employees` table stores information about employees.

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| organization_id | UUID | Foreign key to organizations |
| user_id | UUID | Foreign key to users |
| employee_id | VARCHAR(50) | Employee ID |
| department | VARCHAR(100) | Employee department |
| position | VARCHAR(100) | Employee position |
| hire_date | DATE | Employee hire date |
| manager_id | UUID | Foreign key to employees |
| salary | DECIMAL(12, 2) | Employee salary |
| is_active | BOOLEAN | Whether the employee is active |
| created_at | TIMESTAMP | When the employee was created |
| updated_at | TIMESTAMP | When the employee was last updated |

#### Time Off Requests

The `qultrix.time_off_requests` table stores time off requests.

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| employee_id | UUID | Foreign key to employees |
| start_date | DATE | Start date of time off |
| end_date | DATE | End date of time off |
| type | VARCHAR(50) | Type of time off |
| status | VARCHAR(50) | Status of the request |
| reason | TEXT | Reason for the request |
| approved_by | UUID | Foreign key to employees |
| created_at | TIMESTAMP | When the request was created |
| updated_at | TIMESTAMP | When the request was last updated |

#### Performance Reviews

The `qultrix.performance_reviews` table stores performance reviews.

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| employee_id | UUID | Foreign key to employees |
| reviewer_id | UUID | Foreign key to employees |
| review_date | DATE | Date of the review |
| rating | INTEGER | Review rating |
| comments | TEXT | Review comments |
| created_at | TIMESTAMP | When the review was created |
| updated_at | TIMESTAMP | When the review was last updated |

### AdminHub Schema

The AdminHub schema contains tables related to administration and system management.

#### System Logs

The `adminhub.system_logs` table stores system logs.

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| organization_id | UUID | Foreign key to organizations |
| user_id | UUID | Foreign key to users |
| action | VARCHAR(100) | Action performed |
| resource_type | VARCHAR(100) | Type of resource |
| resource_id | UUID | ID of the resource |
| details | JSONB | Additional details |
| ip_address | VARCHAR(45) | IP address |
| created_at | TIMESTAMP | When the log was created |

#### System Settings

The `adminhub.system_settings` table stores system settings.

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| organization_id | UUID | Foreign key to organizations |
| key | VARCHAR(100) | Setting key |
| value | TEXT | Setting value |
| created_at | TIMESTAMP | When the setting was created |
| updated_at | TIMESTAMP | When the setting was last updated |

### CustomerConnect Schema

The CustomerConnect schema contains tables related to CRM and customer engagement.

#### Contacts

The `customerconnect.contacts` table stores contact information.

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| organization_id | UUID | Foreign key to organizations |
| first_name | VARCHAR(100) | Contact's first name |
| last_name | VARCHAR(100) | Contact's last name |
| email | VARCHAR(255) | Contact's email |
| phone | VARCHAR(50) | Contact's phone |
| company | VARCHAR(100) | Contact's company |
| title | VARCHAR(100) | Contact's title |
| address | TEXT | Contact's address |
| notes | TEXT | Notes about the contact |
| created_by | UUID | Foreign key to users |
| created_at | TIMESTAMP | When the contact was created |
| updated_at | TIMESTAMP | When the contact was last updated |

#### Leads

The `customerconnect.leads` table stores lead information.

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| organization_id | UUID | Foreign key to organizations |
| contact_id | UUID | Foreign key to contacts |
| status | VARCHAR(50) | Lead status |
| source | VARCHAR(100) | Lead source |
| estimated_value | DECIMAL(12, 2) | Estimated value |
| assigned_to | UUID | Foreign key to users |
| notes | TEXT | Notes about the lead |
| created_at | TIMESTAMP | When the lead was created |
| updated_at | TIMESTAMP | When the lead was last updated |

#### Opportunities

The `customerconnect.opportunities` table stores opportunity information.

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| organization_id | UUID | Foreign key to organizations |
| contact_id | UUID | Foreign key to contacts |
| name | VARCHAR(255) | Opportunity name |
| stage | VARCHAR(50) | Opportunity stage |
| amount | DECIMAL(12, 2) | Opportunity amount |
| close_date | DATE | Expected close date |
| probability | INTEGER | Probability of closing |
| assigned_to | UUID | Foreign key to users |
| notes | TEXT | Notes about the opportunity |
| created_at | TIMESTAMP | When the opportunity was created |
| updated_at | TIMESTAMP | When the opportunity was last updated |

### Invantray Schema

The Invantray schema contains tables related to inventory and asset management.

#### Inventory Items

The `invantray.inventory_items` table stores inventory item information.

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| organization_id | UUID | Foreign key to organizations |
| name | VARCHAR(255) | Item name |
| description | TEXT | Item description |
| sku | VARCHAR(100) | Stock keeping unit |
| category | VARCHAR(100) | Item category |
| quantity | INTEGER | Item quantity |
| unit_price | DECIMAL(12, 2) | Unit price |
| reorder_level | INTEGER | Reorder level |
| created_at | TIMESTAMP | When the item was created |
| updated_at | TIMESTAMP | When the item was last updated |

#### Warehouses

The `invantray.warehouses` table stores warehouse information.

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| organization_id | UUID | Foreign key to organizations |
| name | VARCHAR(255) | Warehouse name |
| address | TEXT | Warehouse address |
| manager_id | UUID | Foreign key to users |
| created_at | TIMESTAMP | When the warehouse was created |
| updated_at | TIMESTAMP | When the warehouse was last updated |

#### Inventory Transactions

The `invantray.inventory_transactions` table stores inventory transactions.

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| organization_id | UUID | Foreign key to organizations |
| item_id | UUID | Foreign key to inventory items |
| warehouse_id | UUID | Foreign key to warehouses |
| transaction_type | VARCHAR(50) | Type of transaction |
| quantity | INTEGER | Transaction quantity |
| reference_id | UUID | Reference ID |
| notes | TEXT | Transaction notes |
| created_by | UUID | Foreign key to users |
| created_at | TIMESTAMP | When the transaction was created |

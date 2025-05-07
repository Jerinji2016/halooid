# Employee API Reference

The Employee API allows you to manage employee records in the Qultrix product.

## Base URL

```
/api/v1/organizations/{org_id}/qultrix/employees
```

## Authentication

All endpoints require authentication using a JWT token. The token should be included in the `Authorization` header as a Bearer token.

```
Authorization: Bearer <token>
```

## Permissions

The following permissions are required to access the Employee API:

- `employee:read` - Required to read employee records
- `employee:write` - Required to create and update employee records
- `employee:delete` - Required to delete employee records

## Endpoints

### Create Employee

Creates a new employee record.

**URL**: `POST /api/v1/organizations/{org_id}/qultrix/employees`

**Permissions**: `employee:write`

**Request Body**:

```json
{
  "user_id": "uuid",
  "employee_id": "string",
  "department": "string",
  "position": "string",
  "hire_date": "date",
  "manager_id": "uuid (optional)",
  "salary": "number"
}
```

**Response**: `201 Created`

```json
{
  "id": "uuid",
  "organization_id": "uuid",
  "employee_id": "string",
  "department": "string",
  "position": "string",
  "hire_date": "date",
  "manager_id": "uuid (optional)",
  "salary": "number",
  "is_active": true,
  "created_at": "datetime",
  "updated_at": "datetime",
  "user": {
    "id": "uuid",
    "email": "string",
    "first_name": "string",
    "last_name": "string"
  },
  "manager": {
    "id": "uuid",
    "employee_id": "string",
    "department": "string",
    "position": "string",
    "user": {
      "id": "uuid",
      "email": "string",
      "first_name": "string",
      "last_name": "string"
    }
  }
}
```

**Error Responses**:

- `400 Bad Request` - Invalid request body
- `404 Not Found` - User or manager not found
- `409 Conflict` - Employee ID already exists or user is already an employee

### Get Employee by ID

Retrieves an employee record by ID.

**URL**: `GET /api/v1/organizations/{org_id}/qultrix/employees/{id}`

**Permissions**: `employee:read`

**Response**: `200 OK`

```json
{
  "id": "uuid",
  "organization_id": "uuid",
  "employee_id": "string",
  "department": "string",
  "position": "string",
  "hire_date": "date",
  "manager_id": "uuid (optional)",
  "salary": "number",
  "is_active": true,
  "created_at": "datetime",
  "updated_at": "datetime",
  "user": {
    "id": "uuid",
    "email": "string",
    "first_name": "string",
    "last_name": "string"
  },
  "manager": {
    "id": "uuid",
    "employee_id": "string",
    "department": "string",
    "position": "string",
    "user": {
      "id": "uuid",
      "email": "string",
      "first_name": "string",
      "last_name": "string"
    }
  }
}
```

**Error Responses**:

- `404 Not Found` - Employee not found

### Get Employee by Employee ID

Retrieves an employee record by employee ID.

**URL**: `GET /api/v1/organizations/{org_id}/qultrix/employees/by-employee-id/{employee_id}`

**Permissions**: `employee:read`

**Response**: `200 OK`

```json
{
  "id": "uuid",
  "organization_id": "uuid",
  "employee_id": "string",
  "department": "string",
  "position": "string",
  "hire_date": "date",
  "manager_id": "uuid (optional)",
  "salary": "number",
  "is_active": true,
  "created_at": "datetime",
  "updated_at": "datetime",
  "user": {
    "id": "uuid",
    "email": "string",
    "first_name": "string",
    "last_name": "string"
  },
  "manager": {
    "id": "uuid",
    "employee_id": "string",
    "department": "string",
    "position": "string",
    "user": {
      "id": "uuid",
      "email": "string",
      "first_name": "string",
      "last_name": "string"
    }
  }
}
```

**Error Responses**:

- `404 Not Found` - Employee not found

### List Employees

Retrieves a list of employees based on filter parameters.

**URL**: `GET /api/v1/organizations/{org_id}/qultrix/employees`

**Permissions**: `employee:read`

**Query Parameters**:

- `department` (optional) - Filter by department
- `position` (optional) - Filter by position
- `is_active` (optional) - Filter by active status (true/false)
- `sort_by` (optional) - Sort by field (employee_id, department, position, hire_date, salary)
- `sort_order` (optional) - Sort order (asc/desc)
- `page` (optional) - Page number (default: 1)
- `page_size` (optional) - Page size (default: 10, max: 100)

**Response**: `200 OK`

```json
{
  "employees": [
    {
      "id": "uuid",
      "organization_id": "uuid",
      "employee_id": "string",
      "department": "string",
      "position": "string",
      "hire_date": "date",
      "manager_id": "uuid (optional)",
      "salary": "number",
      "is_active": true,
      "created_at": "datetime",
      "updated_at": "datetime",
      "user": {
        "id": "uuid",
        "email": "string",
        "first_name": "string",
        "last_name": "string"
      }
    }
  ],
  "pagination": {
    "total": "number",
    "page": "number",
    "page_size": "number",
    "total_pages": "number"
  }
}
```

### Update Employee

Updates an employee record.

**URL**: `PUT /api/v1/organizations/{org_id}/qultrix/employees/{id}`

**Permissions**: `employee:write`

**Request Body**:

```json
{
  "user_id": "uuid",
  "employee_id": "string",
  "department": "string",
  "position": "string",
  "hire_date": "date",
  "manager_id": "uuid (optional)",
  "salary": "number"
}
```

**Response**: `200 OK`

```json
{
  "id": "uuid",
  "organization_id": "uuid",
  "employee_id": "string",
  "department": "string",
  "position": "string",
  "hire_date": "date",
  "manager_id": "uuid (optional)",
  "salary": "number",
  "is_active": true,
  "created_at": "datetime",
  "updated_at": "datetime",
  "user": {
    "id": "uuid",
    "email": "string",
    "first_name": "string",
    "last_name": "string"
  },
  "manager": {
    "id": "uuid",
    "employee_id": "string",
    "department": "string",
    "position": "string",
    "user": {
      "id": "uuid",
      "email": "string",
      "first_name": "string",
      "last_name": "string"
    }
  }
}
```

**Error Responses**:

- `400 Bad Request` - Invalid request body
- `404 Not Found` - Employee or manager not found
- `409 Conflict` - Employee ID already exists

### Delete Employee

Marks an employee record as inactive.

**URL**: `DELETE /api/v1/organizations/{org_id}/qultrix/employees/{id}`

**Permissions**: `employee:delete`

**Response**: `204 No Content`

**Error Responses**:

- `404 Not Found` - Employee not found

## Data Models

### Employee

```json
{
  "id": "uuid",
  "organization_id": "uuid",
  "employee_id": "string",
  "department": "string",
  "position": "string",
  "hire_date": "date",
  "manager_id": "uuid (optional)",
  "salary": "number",
  "is_active": "boolean",
  "created_at": "datetime",
  "updated_at": "datetime",
  "user": {
    "id": "uuid",
    "email": "string",
    "first_name": "string",
    "last_name": "string",
    "is_active": "boolean",
    "created_at": "datetime",
    "updated_at": "datetime"
  },
  "manager": {
    "id": "uuid",
    "employee_id": "string",
    "department": "string",
    "position": "string",
    "user": {
      "id": "uuid",
      "email": "string",
      "first_name": "string",
      "last_name": "string"
    }
  }
}
```

### EmployeeRequest

```json
{
  "user_id": "uuid",
  "employee_id": "string",
  "department": "string",
  "position": "string",
  "hire_date": "date",
  "manager_id": "uuid (optional)",
  "salary": "number"
}
```

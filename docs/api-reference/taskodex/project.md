# Project API Reference

The Project API allows you to manage projects in the Taskodex product.

## Base URL

```
/api/v1/organizations/{org_id}/taskodex/projects
```

## Authentication

All endpoints require authentication using a JWT token. The token should be included in the `Authorization` header as a Bearer token.

```
Authorization: Bearer <token>
```

## Permissions

The following permissions are required to access the Project API:

- `project:read` - Required to read projects
- `project:write` - Required to create, update, and manage projects
- `project:delete` - Required to delete projects

## Endpoints

### Create Project

Creates a new project.

**URL**: `POST /api/v1/organizations/{org_id}/taskodex/projects`

**Permissions**: `project:write`

**Request Body**:

```json
{
  "name": "string",
  "description": "string",
  "status": "planning | active | on_hold | completed | cancelled",
  "start_date": "date (optional)",
  "end_date": "date (optional)"
}
```

**Response**: `201 Created`

```json
{
  "id": "uuid",
  "organization_id": "uuid",
  "name": "string",
  "description": "string",
  "status": "planning | active | on_hold | completed | cancelled",
  "start_date": "date (optional)",
  "end_date": "date (optional)",
  "created_by": "uuid",
  "created_at": "datetime",
  "updated_at": "datetime",
  "creator": {
    "id": "uuid",
    "email": "string",
    "first_name": "string",
    "last_name": "string"
  }
}
```

**Error Responses**:

- `400 Bad Request` - Invalid request body
- `409 Conflict` - Project name already exists in this organization

### Get Project by ID

Retrieves a project by ID.

**URL**: `GET /api/v1/organizations/{org_id}/taskodex/projects/{id}`

**Permissions**: `project:read`

**Response**: `200 OK`

```json
{
  "id": "uuid",
  "organization_id": "uuid",
  "name": "string",
  "description": "string",
  "status": "planning | active | on_hold | completed | cancelled",
  "start_date": "date (optional)",
  "end_date": "date (optional)",
  "created_by": "uuid",
  "created_at": "datetime",
  "updated_at": "datetime",
  "creator": {
    "id": "uuid",
    "email": "string",
    "first_name": "string",
    "last_name": "string"
  }
}
```

**Error Responses**:

- `404 Not Found` - Project not found

### List Projects

Retrieves a list of projects based on filter parameters.

**URL**: `GET /api/v1/organizations/{org_id}/taskodex/projects`

**Permissions**: `project:read`

**Query Parameters**:

- `status` (optional) - Filter by status (planning, active, on_hold, completed, cancelled)
- `created_by` (optional) - Filter by creator ID
- `search` (optional) - Search in name and description
- `sort_by` (optional) - Sort by field (name, status, start_date, end_date, created_at, updated_at)
- `sort_order` (optional) - Sort order (asc/desc)
- `page` (optional) - Page number (default: 1)
- `page_size` (optional) - Page size (default: 20, max: 100)

**Response**: `200 OK`

```json
{
  "projects": [
    {
      "id": "uuid",
      "organization_id": "uuid",
      "name": "string",
      "description": "string",
      "status": "planning | active | on_hold | completed | cancelled",
      "start_date": "date (optional)",
      "end_date": "date (optional)",
      "created_by": "uuid",
      "created_at": "datetime",
      "updated_at": "datetime",
      "creator": {
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

### Update Project

Updates a project.

**URL**: `PUT /api/v1/organizations/{org_id}/taskodex/projects/{id}`

**Permissions**: `project:write`

**Request Body**:

```json
{
  "name": "string",
  "description": "string",
  "status": "planning | active | on_hold | completed | cancelled",
  "start_date": "date (optional)",
  "end_date": "date (optional)"
}
```

**Response**: `200 OK`

```json
{
  "id": "uuid",
  "organization_id": "uuid",
  "name": "string",
  "description": "string",
  "status": "planning | active | on_hold | completed | cancelled",
  "start_date": "date (optional)",
  "end_date": "date (optional)",
  "created_by": "uuid",
  "created_at": "datetime",
  "updated_at": "datetime",
  "creator": {
    "id": "uuid",
    "email": "string",
    "first_name": "string",
    "last_name": "string"
  }
}
```

**Error Responses**:

- `400 Bad Request` - Invalid request body
- `404 Not Found` - Project not found
- `409 Conflict` - Project name already exists in this organization

### Delete Project

Deletes a project.

**URL**: `DELETE /api/v1/organizations/{org_id}/taskodex/projects/{id}`

**Permissions**: `project:delete`

**Response**: `204 No Content`

**Error Responses**:

- `400 Bad Request` - Cannot delete project with tasks
- `404 Not Found` - Project not found

### Get Project Tasks

Retrieves all tasks for a project.

**URL**: `GET /api/v1/organizations/{org_id}/taskodex/projects/{id}/tasks`

**Permissions**: `project:read`

**Query Parameters**:

- `status` (optional) - Filter by status (todo, in_progress, review, done, cancelled)
- `priority` (optional) - Filter by priority (low, medium, high, critical)
- `assigned_to` (optional) - Filter by assignee ID
- `due_before` (optional) - Filter by due date before (ISO 8601 format)
- `due_after` (optional) - Filter by due date after (ISO 8601 format)
- `search` (optional) - Search in title and description
- `sort_by` (optional) - Sort by field (title, status, priority, due_date, created_at, updated_at)
- `sort_order` (optional) - Sort order (asc/desc)
- `page` (optional) - Page number (default: 1)
- `page_size` (optional) - Page size (default: 20, max: 100)

**Response**: `200 OK`

```json
{
  "tasks": [
    {
      "id": "uuid",
      "project_id": "uuid",
      "title": "string",
      "description": "string",
      "status": "todo | in_progress | review | done | cancelled",
      "priority": "low | medium | high | critical",
      "due_date": "date (optional)",
      "created_by": "uuid",
      "assigned_to": "uuid (optional)",
      "estimated_hours": "number (optional)",
      "actual_hours": "number (optional)",
      "tags": ["string"],
      "created_at": "datetime",
      "updated_at": "datetime",
      "creator": {
        "id": "uuid",
        "email": "string",
        "first_name": "string",
        "last_name": "string"
      },
      "assignee": {
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

**Error Responses**:

- `404 Not Found` - Project not found

### Add Task to Project

Adds a task to a project.

**URL**: `POST /api/v1/organizations/{org_id}/taskodex/projects/{id}/tasks`

**Permissions**: `project:write`

**Request Body**:

```json
{
  "task_id": "uuid"
}
```

**Response**: `204 No Content`

**Error Responses**:

- `400 Bad Request` - Invalid request body
- `404 Not Found` - Project or task not found

### Remove Task from Project

Removes a task from a project.

**URL**: `DELETE /api/v1/organizations/{org_id}/taskodex/projects/{id}/tasks/{task_id}`

**Permissions**: `project:write`

**Response**: `204 No Content`

**Error Responses**:

- `400 Bad Request` - Task does not belong to the project
- `404 Not Found` - Project or task not found

## Data Models

### Project

```json
{
  "id": "uuid",
  "organization_id": "uuid",
  "name": "string",
  "description": "string",
  "status": "planning | active | on_hold | completed | cancelled",
  "start_date": "date (optional)",
  "end_date": "date (optional)",
  "created_by": "uuid",
  "created_at": "datetime",
  "updated_at": "datetime",
  "creator": {
    "id": "uuid",
    "email": "string",
    "first_name": "string",
    "last_name": "string"
  }
}
```

### ProjectRequest

```json
{
  "name": "string",
  "description": "string",
  "status": "planning | active | on_hold | completed | cancelled",
  "start_date": "date (optional)",
  "end_date": "date (optional)"
}
```

# Task API Reference

The Task API allows you to manage tasks in the Taskodex product.

## Base URL

```
/api/v1/organizations/{org_id}/taskodex/tasks
```

## Authentication

All endpoints require authentication using a JWT token. The token should be included in the `Authorization` header as a Bearer token.

```
Authorization: Bearer <token>
```

## Permissions

The following permissions are required to access the Task API:

- `task:read` - Required to read tasks
- `task:write` - Required to create, update, and manage tasks
- `task:delete` - Required to delete tasks

## Endpoints

### Create Task

Creates a new task.

**URL**: `POST /api/v1/organizations/{org_id}/taskodex/tasks`

**Permissions**: `task:write`

**Request Body**:

```json
{
  "project_id": "uuid (optional)",
  "title": "string",
  "description": "string",
  "status": "todo | in_progress | review | done | cancelled",
  "priority": "low | medium | high | critical",
  "due_date": "date (optional)",
  "assigned_to": "uuid (optional)",
  "estimated_hours": "number (optional)",
  "tags": ["string"] (optional)
}
```

**Response**: `201 Created`

```json
{
  "id": "uuid",
  "project_id": "uuid (optional)",
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
  "project": {
    "id": "uuid",
    "name": "string",
    "description": "string",
    "status": "string"
  },
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
```

**Error Responses**:

- `400 Bad Request` - Invalid request body
- `404 Not Found` - Project or user not found

### Get Task by ID

Retrieves a task by ID.

**URL**: `GET /api/v1/organizations/{org_id}/taskodex/tasks/{id}`

**Permissions**: `task:read`

**Response**: `200 OK`

```json
{
  "id": "uuid",
  "project_id": "uuid (optional)",
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
  "project": {
    "id": "uuid",
    "name": "string",
    "description": "string",
    "status": "string"
  },
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
```

**Error Responses**:

- `404 Not Found` - Task not found

### List Tasks

Retrieves a list of tasks based on filter parameters.

**URL**: `GET /api/v1/organizations/{org_id}/taskodex/tasks`

**Permissions**: `task:read`

**Query Parameters**:

- `project_id` (optional) - Filter by project ID
- `status` (optional) - Filter by status (todo, in_progress, review, done, cancelled)
- `priority` (optional) - Filter by priority (low, medium, high, critical)
- `created_by` (optional) - Filter by creator ID
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
      "project_id": "uuid (optional)",
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

### Update Task

Updates a task.

**URL**: `PUT /api/v1/organizations/{org_id}/taskodex/tasks/{id}`

**Permissions**: `task:write`

**Request Body**:

```json
{
  "project_id": "uuid (optional)",
  "title": "string",
  "description": "string",
  "status": "todo | in_progress | review | done | cancelled",
  "priority": "low | medium | high | critical",
  "due_date": "date (optional)",
  "assigned_to": "uuid (optional)",
  "estimated_hours": "number (optional)",
  "tags": ["string"] (optional)
}
```

**Response**: `200 OK`

```json
{
  "id": "uuid",
  "project_id": "uuid (optional)",
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
  "project": {
    "id": "uuid",
    "name": "string",
    "description": "string",
    "status": "string"
  },
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
```

**Error Responses**:

- `400 Bad Request` - Invalid request body
- `404 Not Found` - Task, project, or user not found

### Delete Task

Deletes a task.

**URL**: `DELETE /api/v1/organizations/{org_id}/taskodex/tasks/{id}`

**Permissions**: `task:delete`

**Response**: `204 No Content`

**Error Responses**:

- `404 Not Found` - Task not found

### Add Tag to Task

Adds a tag to a task.

**URL**: `POST /api/v1/organizations/{org_id}/taskodex/tasks/{id}/tags`

**Permissions**: `task:write`

**Request Body**:

```json
{
  "tag": "string"
}
```

**Response**: `204 No Content`

**Error Responses**:

- `400 Bad Request` - Invalid request body
- `404 Not Found` - Task not found

### Remove Tag from Task

Removes a tag from a task.

**URL**: `DELETE /api/v1/organizations/{org_id}/taskodex/tasks/{id}/tags/{tag}`

**Permissions**: `task:write`

**Response**: `204 No Content`

**Error Responses**:

- `404 Not Found` - Task not found

## Data Models

### Task

```json
{
  "id": "uuid",
  "project_id": "uuid (optional)",
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
  "project": {
    "id": "uuid",
    "name": "string",
    "description": "string",
    "status": "string"
  },
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
```

### TaskRequest

```json
{
  "project_id": "uuid (optional)",
  "title": "string",
  "description": "string",
  "status": "todo | in_progress | review | done | cancelled",
  "priority": "low | medium | high | critical",
  "due_date": "date (optional)",
  "assigned_to": "uuid (optional)",
  "estimated_hours": "number (optional)",
  "tags": ["string"] (optional)
}
```

# Task Assignment and Tracking API Reference

The Task Assignment and Tracking API allows you to manage task assignments and track task progress in the Taskodex product.

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

The following permissions are required to access the Task Assignment and Tracking API:

- `task:read` - Required to read tasks
- `task:write` - Required to assign, unassign, and update task status

## Endpoints

### Assign Task

Assigns a task to a user.

**URL**: `POST /api/v1/organizations/{org_id}/taskodex/tasks/{id}/assign`

**Permissions**: `task:write`

**Request Body**:

```json
{
  "user_id": "uuid"
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
  "assigned_to": "uuid",
  "estimated_hours": "number (optional)",
  "actual_hours": "number (optional)",
  "tags": ["string"],
  "created_at": "datetime",
  "updated_at": "datetime",
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
- `404 Not Found` - Task or user not found

### Unassign Task

Removes the assignment of a task.

**URL**: `POST /api/v1/organizations/{org_id}/taskodex/tasks/{id}/unassign`

**Permissions**: `task:write`

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
  "assigned_to": null,
  "estimated_hours": "number (optional)",
  "actual_hours": "number (optional)",
  "tags": ["string"],
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

**Error Responses**:

- `400 Bad Request` - Task is not assigned to anyone
- `404 Not Found` - Task not found

### Update Task Status

Updates the status of a task.

**URL**: `PUT /api/v1/organizations/{org_id}/taskodex/tasks/{id}/status`

**Permissions**: `task:write`

**Request Body**:

```json
{
  "status": "todo | in_progress | review | done | cancelled"
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
  "updated_at": "datetime"
}
```

**Error Responses**:

- `400 Bad Request` - Invalid request body or invalid task status
- `404 Not Found` - Task not found

### Get Tasks by Assignee

Retrieves tasks assigned to a user.

**URL**: `GET /api/v1/organizations/{org_id}/taskodex/tasks/assignee/{user_id}`

**Permissions**: `task:read`

**Query Parameters**:

- `status` (optional) - Filter by status (todo, in_progress, review, done, cancelled)
- `priority` (optional) - Filter by priority (low, medium, high, critical)
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
      "assigned_to": "uuid",
      "estimated_hours": "number (optional)",
      "actual_hours": "number (optional)",
      "tags": ["string"],
      "created_at": "datetime",
      "updated_at": "datetime"
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

- `404 Not Found` - User not found

### Get Overdue Tasks

Retrieves overdue tasks.

**URL**: `GET /api/v1/organizations/{org_id}/taskodex/tasks/overdue`

**Permissions**: `task:read`

**Query Parameters**:

- `project_id` (optional) - Filter by project ID
- `assigned_to` (optional) - Filter by assignee ID
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
      "status": "todo | in_progress | review",
      "priority": "low | medium | high | critical",
      "due_date": "date",
      "created_by": "uuid",
      "assigned_to": "uuid (optional)",
      "estimated_hours": "number (optional)",
      "actual_hours": "number (optional)",
      "tags": ["string"],
      "created_at": "datetime",
      "updated_at": "datetime"
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

### Get Tasks Due Soon

Retrieves tasks due within a specified number of days.

**URL**: `GET /api/v1/organizations/{org_id}/taskodex/tasks/due-soon/{days}`

**Permissions**: `task:read`

**Path Parameters**:

- `days` - Number of days to look ahead

**Query Parameters**:

- `project_id` (optional) - Filter by project ID
- `assigned_to` (optional) - Filter by assignee ID
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
      "status": "todo | in_progress | review",
      "priority": "low | medium | high | critical",
      "due_date": "date",
      "created_by": "uuid",
      "assigned_to": "uuid (optional)",
      "estimated_hours": "number (optional)",
      "actual_hours": "number (optional)",
      "tags": ["string"],
      "created_at": "datetime",
      "updated_at": "datetime"
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

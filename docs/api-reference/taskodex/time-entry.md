# Time Entry API Reference

The Time Entry API allows you to manage time entries for tasks in the Taskodex product.

## Base URL

```
/api/v1/organizations/{org_id}/taskodex/time-entries
```

## Authentication

All endpoints require authentication using a JWT token. The token should be included in the `Authorization` header as a Bearer token.

```
Authorization: Bearer <token>
```

## Permissions

The following permissions are required to access the Time Entry API:

- `time_entry:read` - Required to read time entries
- `time_entry:write` - Required to create, update, and manage time entries
- `time_entry:delete` - Required to delete time entries

## Endpoints

### Create Time Entry

Creates a new time entry.

**URL**: `POST /api/v1/organizations/{org_id}/taskodex/time-entries`

**Permissions**: `time_entry:write`

**Request Body**:

```json
{
  "task_id": "uuid",
  "start_time": "datetime",
  "end_time": "datetime (optional)",
  "duration_minutes": "number (optional)",
  "description": "string (optional)"
}
```

**Response**: `201 Created`

```json
{
  "id": "uuid",
  "task_id": "uuid",
  "user_id": "uuid",
  "start_time": "datetime",
  "end_time": "datetime (optional)",
  "duration_minutes": "number (optional)",
  "description": "string",
  "created_at": "datetime",
  "updated_at": "datetime",
  "task": {
    "id": "uuid",
    "title": "string",
    "status": "todo | in_progress | review | done | cancelled"
  },
  "user": {
    "id": "uuid",
    "email": "string",
    "first_name": "string",
    "last_name": "string"
  }
}
```

**Error Responses**:

- `400 Bad Request` - Invalid request body or end time cannot be before start time
- `404 Not Found` - Task not found
- `409 Conflict` - You already have a running timer for this task

### Get Time Entry by ID

Retrieves a time entry by ID.

**URL**: `GET /api/v1/organizations/{org_id}/taskodex/time-entries/{id}`

**Permissions**: `time_entry:read`

**Response**: `200 OK`

```json
{
  "id": "uuid",
  "task_id": "uuid",
  "user_id": "uuid",
  "start_time": "datetime",
  "end_time": "datetime (optional)",
  "duration_minutes": "number (optional)",
  "description": "string",
  "created_at": "datetime",
  "updated_at": "datetime",
  "task": {
    "id": "uuid",
    "title": "string",
    "status": "todo | in_progress | review | done | cancelled"
  },
  "user": {
    "id": "uuid",
    "email": "string",
    "first_name": "string",
    "last_name": "string"
  }
}
```

**Error Responses**:

- `404 Not Found` - Time entry not found

### List Time Entries

Retrieves a list of time entries based on filter parameters.

**URL**: `GET /api/v1/organizations/{org_id}/taskodex/time-entries`

**Permissions**: `time_entry:read`

**Query Parameters**:

- `task_id` (optional) - Filter by task ID
- `user_id` (optional) - Filter by user ID
- `start_after` (optional) - Filter by start time after (ISO 8601 format)
- `start_before` (optional) - Filter by start time before (ISO 8601 format)
- `is_running` (optional) - Filter by running status (true/false)
- `sort_by` (optional) - Sort by field (start_time, end_time, created_at, updated_at)
- `sort_order` (optional) - Sort order (asc/desc)
- `page` (optional) - Page number (default: 1)
- `page_size` (optional) - Page size (default: 20, max: 100)

**Response**: `200 OK`

```json
{
  "time_entries": [
    {
      "id": "uuid",
      "task_id": "uuid",
      "user_id": "uuid",
      "start_time": "datetime",
      "end_time": "datetime (optional)",
      "duration_minutes": "number (optional)",
      "description": "string",
      "created_at": "datetime",
      "updated_at": "datetime",
      "task": {
        "id": "uuid",
        "title": "string",
        "status": "todo | in_progress | review | done | cancelled"
      },
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

### Update Time Entry

Updates a time entry.

**URL**: `PUT /api/v1/organizations/{org_id}/taskodex/time-entries/{id}`

**Permissions**: `time_entry:write`

**Request Body**:

```json
{
  "task_id": "uuid",
  "start_time": "datetime",
  "end_time": "datetime (optional)",
  "duration_minutes": "number (optional)",
  "description": "string (optional)"
}
```

**Response**: `200 OK`

```json
{
  "id": "uuid",
  "task_id": "uuid",
  "user_id": "uuid",
  "start_time": "datetime",
  "end_time": "datetime (optional)",
  "duration_minutes": "number (optional)",
  "description": "string",
  "created_at": "datetime",
  "updated_at": "datetime",
  "task": {
    "id": "uuid",
    "title": "string",
    "status": "todo | in_progress | review | done | cancelled"
  },
  "user": {
    "id": "uuid",
    "email": "string",
    "first_name": "string",
    "last_name": "string"
  }
}
```

**Error Responses**:

- `400 Bad Request` - Invalid request body or end time cannot be before start time
- `404 Not Found` - Time entry or task not found

### Delete Time Entry

Deletes a time entry.

**URL**: `DELETE /api/v1/organizations/{org_id}/taskodex/time-entries/{id}`

**Permissions**: `time_entry:delete`

**Response**: `204 No Content`

**Error Responses**:

- `404 Not Found` - Time entry not found

### Start Timer

Starts a timer for a task.

**URL**: `POST /api/v1/organizations/{org_id}/taskodex/time-entries/start`

**Permissions**: `time_entry:write`

**Request Body**:

```json
{
  "task_id": "uuid",
  "description": "string (optional)"
}
```

**Response**: `200 OK`

```json
{
  "id": "uuid",
  "task_id": "uuid",
  "user_id": "uuid",
  "start_time": "datetime",
  "end_time": null,
  "duration_minutes": null,
  "description": "string",
  "created_at": "datetime",
  "updated_at": "datetime",
  "task": {
    "id": "uuid",
    "title": "string",
    "status": "todo | in_progress | review | done | cancelled"
  },
  "user": {
    "id": "uuid",
    "email": "string",
    "first_name": "string",
    "last_name": "string"
  }
}
```

**Error Responses**:

- `400 Bad Request` - Invalid request body
- `404 Not Found` - Task not found
- `409 Conflict` - You already have a running timer for this task

### Stop Timer

Stops a running timer for a task.

**URL**: `POST /api/v1/organizations/{org_id}/taskodex/time-entries/stop`

**Permissions**: `time_entry:write`

**Request Body**:

```json
{
  "task_id": "uuid"
}
```

**Response**: `200 OK`

```json
{
  "id": "uuid",
  "task_id": "uuid",
  "user_id": "uuid",
  "start_time": "datetime",
  "end_time": "datetime",
  "duration_minutes": "number",
  "description": "string",
  "created_at": "datetime",
  "updated_at": "datetime",
  "task": {
    "id": "uuid",
    "title": "string",
    "status": "todo | in_progress | review | done | cancelled"
  },
  "user": {
    "id": "uuid",
    "email": "string",
    "first_name": "string",
    "last_name": "string"
  }
}
```

**Error Responses**:

- `400 Bad Request` - Invalid request body
- `404 Not Found` - No running timer found for this task

### Get Running Timers

Retrieves all running timers for the authenticated user.

**URL**: `GET /api/v1/organizations/{org_id}/taskodex/time-entries/running`

**Permissions**: `time_entry:read`

**Response**: `200 OK`

```json
[
  {
    "id": "uuid",
    "task_id": "uuid",
    "user_id": "uuid",
    "start_time": "datetime",
    "end_time": null,
    "duration_minutes": null,
    "description": "string",
    "created_at": "datetime",
    "updated_at": "datetime",
    "task": {
      "id": "uuid",
      "title": "string",
      "status": "todo | in_progress | review | done | cancelled"
    },
    "user": {
      "id": "uuid",
      "email": "string",
      "first_name": "string",
      "last_name": "string"
    }
  }
]
```

**Error Responses**:

- `404 Not Found` - User not found

### Aggregate Time Entries

Aggregates time entries based on parameters.

**URL**: `GET /api/v1/organizations/{org_id}/taskodex/time-entries/aggregate`

**Permissions**: `time_entry:read`

**Query Parameters**:

- `task_id` (optional) - Filter by task ID
- `user_id` (optional) - Filter by user ID
- `project_id` (optional) - Filter by project ID
- `start_after` (optional) - Filter by start time after (ISO 8601 format)
- `start_before` (optional) - Filter by start time before (ISO 8601 format)
- `group_by` (optional) - Group by field (day, week, month, year, task, user) (default: day)

**Response**: `200 OK`

```json
[
  {
    "total_duration_minutes": "number",
    "task_id": "uuid (if group_by=task)",
    "user_id": "uuid (if group_by=user)",
    "date": "date (if group_by=day)",
    "week": "number (if group_by=week)",
    "month": "number (if group_by=month)",
    "year": "number (if group_by=week, month, or year)"
  }
]
```

**Error Responses**:

- `400 Bad Request` - Invalid group_by parameter

## Data Models

### TimeEntry

```json
{
  "id": "uuid",
  "task_id": "uuid",
  "user_id": "uuid",
  "start_time": "datetime",
  "end_time": "datetime (optional)",
  "duration_minutes": "number (optional)",
  "description": "string",
  "created_at": "datetime",
  "updated_at": "datetime",
  "task": {
    "id": "uuid",
    "title": "string",
    "status": "todo | in_progress | review | done | cancelled"
  },
  "user": {
    "id": "uuid",
    "email": "string",
    "first_name": "string",
    "last_name": "string"
  }
}
```

### TimeEntryRequest

```json
{
  "task_id": "uuid",
  "start_time": "datetime",
  "end_time": "datetime (optional)",
  "duration_minutes": "number (optional)",
  "description": "string (optional)"
}
```

### TimeEntryAggregation

```json
{
  "total_duration_minutes": "number",
  "task_id": "uuid (if group_by=task)",
  "user_id": "uuid (if group_by=user)",
  "date": "date (if group_by=day)",
  "week": "number (if group_by=week)",
  "month": "number (if group_by=month)",
  "year": "number (if group_by=week, month, or year)"
}
```

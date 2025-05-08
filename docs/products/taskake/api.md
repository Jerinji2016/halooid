# Taskake API

The Taskake API provides a comprehensive set of endpoints for managing tasks, projects, time entries, comments, and file attachments. This documentation outlines the available endpoints, request/response formats, and authentication requirements.

## Authentication

All API requests require authentication using a JWT token. The token should be included in the `Authorization` header as a Bearer token.

```
Authorization: Bearer <token>
```

## Base URL

```
/api/v1/organizations/{org_id}/taskodex
```

## Task Management

### Create Task

Creates a new task.

**URL**: `POST /tasks`

**Permissions**: `task:write`

**Request Body**:
```json
{
  "title": "Implement API documentation",
  "description": "Create comprehensive API documentation for the Taskake API",
  "status": "todo",
  "priority": "high",
  "project_id": "123e4567-e89b-12d3-a456-426614174000",
  "due_date": "2023-12-31T23:59:59Z",
  "assigned_to": "123e4567-e89b-12d3-a456-426614174001",
  "estimated_hours": 8,
  "tags": ["documentation", "api"]
}
```

**Response**: `201 Created`
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174002",
  "title": "Implement API documentation",
  "description": "Create comprehensive API documentation for the Taskake API",
  "status": "todo",
  "priority": "high",
  "project_id": "123e4567-e89b-12d3-a456-426614174000",
  "due_date": "2023-12-31T23:59:59Z",
  "created_by": "123e4567-e89b-12d3-a456-426614174003",
  "assigned_to": "123e4567-e89b-12d3-a456-426614174001",
  "estimated_hours": 8,
  "actual_hours": null,
  "tags": ["documentation", "api"],
  "created_at": "2023-06-01T12:00:00Z",
  "updated_at": "2023-06-01T12:00:00Z"
}
```

### Get Task by ID

Retrieves a task by its ID.

**URL**: `GET /tasks/{id}`

**Permissions**: `task:read`

**Response**: `200 OK`
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174002",
  "title": "Implement API documentation",
  "description": "Create comprehensive API documentation for the Taskake API",
  "status": "todo",
  "priority": "high",
  "project_id": "123e4567-e89b-12d3-a456-426614174000",
  "due_date": "2023-12-31T23:59:59Z",
  "created_by": "123e4567-e89b-12d3-a456-426614174003",
  "assigned_to": "123e4567-e89b-12d3-a456-426614174001",
  "estimated_hours": 8,
  "actual_hours": null,
  "tags": ["documentation", "api"],
  "created_at": "2023-06-01T12:00:00Z",
  "updated_at": "2023-06-01T12:00:00Z",
  "project": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "name": "API Development",
    "description": "Development of the Taskake API",
    "status": "active"
  },
  "creator": {
    "id": "123e4567-e89b-12d3-a456-426614174003",
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@example.com"
  },
  "assignee": {
    "id": "123e4567-e89b-12d3-a456-426614174001",
    "first_name": "Jane",
    "last_name": "Smith",
    "email": "jane.smith@example.com"
  }
}
```

### List Tasks

Retrieves a list of tasks based on filter parameters.

**URL**: `GET /tasks`

**Permissions**: `task:read`

**Query Parameters**:
- `page` (optional): Page number (default: 1)
- `page_size` (optional): Number of items per page (default: 10)
- `status` (optional): Filter by status (e.g., todo, in_progress, done)
- `priority` (optional): Filter by priority (e.g., low, medium, high, critical)
- `project_id` (optional): Filter by project ID
- `assigned_to` (optional): Filter by assignee ID
- `created_by` (optional): Filter by creator ID
- `due_date_from` (optional): Filter by due date (from)
- `due_date_to` (optional): Filter by due date (to)
- `sort_by` (optional): Field to sort by (default: created_at)
- `sort_order` (optional): Sort order (asc or desc, default: desc)

**Response**: `200 OK`
```json
{
  "tasks": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174002",
      "title": "Implement API documentation",
      "description": "Create comprehensive API documentation for the Taskake API",
      "status": "todo",
      "priority": "high",
      "project_id": "123e4567-e89b-12d3-a456-426614174000",
      "due_date": "2023-12-31T23:59:59Z",
      "created_by": "123e4567-e89b-12d3-a456-426614174003",
      "assigned_to": "123e4567-e89b-12d3-a456-426614174001",
      "estimated_hours": 8,
      "actual_hours": null,
      "tags": ["documentation", "api"],
      "created_at": "2023-06-01T12:00:00Z",
      "updated_at": "2023-06-01T12:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 10,
    "total": 1,
    "total_pages": 1
  }
}
```

### Update Task

Updates an existing task.

**URL**: `PUT /tasks/{id}`

**Permissions**: `task:write`

**Request Body**:
```json
{
  "title": "Implement API documentation",
  "description": "Create comprehensive API documentation for the Taskake API",
  "status": "in_progress",
  "priority": "high",
  "project_id": "123e4567-e89b-12d3-a456-426614174000",
  "due_date": "2023-12-31T23:59:59Z",
  "assigned_to": "123e4567-e89b-12d3-a456-426614174001",
  "estimated_hours": 8,
  "tags": ["documentation", "api"]
}
```

**Response**: `200 OK`
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174002",
  "title": "Implement API documentation",
  "description": "Create comprehensive API documentation for the Taskake API",
  "status": "in_progress",
  "priority": "high",
  "project_id": "123e4567-e89b-12d3-a456-426614174000",
  "due_date": "2023-12-31T23:59:59Z",
  "created_by": "123e4567-e89b-12d3-a456-426614174003",
  "assigned_to": "123e4567-e89b-12d3-a456-426614174001",
  "estimated_hours": 8,
  "actual_hours": null,
  "tags": ["documentation", "api"],
  "created_at": "2023-06-01T12:00:00Z",
  "updated_at": "2023-06-01T12:30:00Z"
}
```

### Delete Task

Deletes a task.

**URL**: `DELETE /tasks/{id}`

**Permissions**: `task:delete`

**Response**: `204 No Content`

## Project Management

### Create Project

Creates a new project.

**URL**: `POST /projects`

**Permissions**: `project:write`

**Request Body**:
```json
{
  "name": "API Development",
  "description": "Development of the Taskake API",
  "status": "planning",
  "start_date": "2023-06-01",
  "end_date": "2023-12-31"
}
```

**Response**: `201 Created`
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "organization_id": "123e4567-e89b-12d3-a456-426614174004",
  "name": "API Development",
  "description": "Development of the Taskake API",
  "status": "planning",
  "start_date": "2023-06-01",
  "end_date": "2023-12-31",
  "created_by": "123e4567-e89b-12d3-a456-426614174003",
  "created_at": "2023-06-01T12:00:00Z",
  "updated_at": "2023-06-01T12:00:00Z"
}
```

### Get Project by ID

Retrieves a project by its ID.

**URL**: `GET /projects/{id}`

**Permissions**: `project:read`

**Response**: `200 OK`
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "organization_id": "123e4567-e89b-12d3-a456-426614174004",
  "name": "API Development",
  "description": "Development of the Taskake API",
  "status": "planning",
  "start_date": "2023-06-01",
  "end_date": "2023-12-31",
  "created_by": "123e4567-e89b-12d3-a456-426614174003",
  "created_at": "2023-06-01T12:00:00Z",
  "updated_at": "2023-06-01T12:00:00Z",
  "creator": {
    "id": "123e4567-e89b-12d3-a456-426614174003",
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@example.com"
  }
}
```

## Time Tracking

### Create Time Entry

Creates a new time entry for a task.

**URL**: `POST /time-entries`

**Permissions**: `time_entry:write`

**Request Body**:
```json
{
  "task_id": "123e4567-e89b-12d3-a456-426614174002",
  "start_time": "2023-06-01T13:00:00Z",
  "end_time": "2023-06-01T15:00:00Z",
  "duration_minutes": 120,
  "description": "Working on API documentation"
}
```

**Response**: `201 Created`
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174005",
  "task_id": "123e4567-e89b-12d3-a456-426614174002",
  "user_id": "123e4567-e89b-12d3-a456-426614174003",
  "start_time": "2023-06-01T13:00:00Z",
  "end_time": "2023-06-01T15:00:00Z",
  "duration_minutes": 120,
  "description": "Working on API documentation",
  "created_at": "2023-06-01T15:00:00Z",
  "updated_at": "2023-06-01T15:00:00Z"
}
```

### Start Timer

Starts a timer for a task.

**URL**: `POST /time-entries/start`

**Permissions**: `time_entry:write`

**Request Body**:
```json
{
  "task_id": "123e4567-e89b-12d3-a456-426614174002",
  "description": "Working on API documentation"
}
```

**Response**: `200 OK`
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174006",
  "task_id": "123e4567-e89b-12d3-a456-426614174002",
  "user_id": "123e4567-e89b-12d3-a456-426614174003",
  "start_time": "2023-06-01T16:00:00Z",
  "end_time": null,
  "duration_minutes": null,
  "description": "Working on API documentation",
  "created_at": "2023-06-01T16:00:00Z",
  "updated_at": "2023-06-01T16:00:00Z"
}
```

## Comments

### Create Comment

Creates a new comment on a task.

**URL**: `POST /comments`

**Permissions**: `comment:write`

**Request Body**:
```json
{
  "task_id": "123e4567-e89b-12d3-a456-426614174002",
  "content": "This is a comment on the task."
}
```

**Response**: `201 Created`
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174007",
  "task_id": "123e4567-e89b-12d3-a456-426614174002",
  "user_id": "123e4567-e89b-12d3-a456-426614174003",
  "content": "This is a comment on the task.",
  "created_at": "2023-06-01T17:00:00Z",
  "updated_at": "2023-06-01T17:00:00Z",
  "user": {
    "id": "123e4567-e89b-12d3-a456-426614174003",
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@example.com"
  }
}
```

## File Attachments

### Upload File

Uploads a file and attaches it to a task.

**URL**: `POST /files`

**Permissions**: `file:write`

**Request Body**: Multipart form data
- `task_id`: ID of the task to attach the file to
- `file`: The file to upload

**Response**: `201 Created`
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174008",
  "task_id": "123e4567-e89b-12d3-a456-426614174002",
  "user_id": "123e4567-e89b-12d3-a456-426614174003",
  "file_name": "api-documentation.pdf",
  "file_size": 1024000,
  "content_type": "application/pdf",
  "download_url": "https://api.halooid.com/v1/files/123e4567-e89b-12d3-a456-426614174008/download",
  "created_at": "2023-06-01T18:00:00Z"
}
```

## Notifications

### List Notifications

Retrieves a list of notifications for the authenticated user.

**URL**: `GET /notifications`

**Query Parameters**:
- `page` (optional): Page number (default: 1)
- `page_size` (optional): Number of items per page (default: 10)
- `is_read` (optional): Filter by read status (true or false)
- `type` (optional): Filter by notification type

**Response**: `200 OK`
```json
{
  "notifications": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174009",
      "user_id": "123e4567-e89b-12d3-a456-426614174003",
      "type": "task_assigned",
      "title": "Task Assigned",
      "message": "You have been assigned a new task: Implement API documentation",
      "resource_type": "task",
      "resource_id": "123e4567-e89b-12d3-a456-426614174002",
      "is_read": false,
      "created_at": "2023-06-01T12:00:00Z",
      "read_at": null
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 10,
    "total": 1,
    "total_pages": 1
  }
}
```

### Mark Notification as Read

Marks a notification as read.

**URL**: `PUT /notifications/{id}/read`

**Response**: `204 No Content`

### Mark All Notifications as Read

Marks all notifications for the authenticated user as read.

**URL**: `PUT /notifications/read-all`

**Response**: `204 No Content`

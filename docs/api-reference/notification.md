# Notification API Reference

The Notification API allows you to manage notifications in the system.

## Base URL

```
/api/v1/notifications
```

## Authentication

All endpoints require authentication using a JWT token. The token should be included in the `Authorization` header as a Bearer token.

```
Authorization: Bearer <token>
```

## Endpoints

### List Notifications

Retrieves a list of notifications for the authenticated user.

**URL**: `GET /api/v1/notifications`

**Query Parameters**:

- `type` (optional) - Filter by notification type (task_assigned, task_status_update, task_due_soon, task_overdue, task_comment, task_mention)
- `is_read` (optional) - Filter by read status (true/false)
- `sort_by` (optional) - Sort by field (created_at, type, is_read)
- `sort_order` (optional) - Sort order (asc/desc)
- `page` (optional) - Page number (default: 1)
- `page_size` (optional) - Page size (default: 20, max: 100)

**Response**: `200 OK`

```json
{
  "notifications": [
    {
      "id": "uuid",
      "user_id": "uuid",
      "type": "task_assigned | task_status_update | task_due_soon | task_overdue | task_comment | task_mention",
      "title": "string",
      "message": "string",
      "resource_type": "string",
      "resource_id": "uuid",
      "is_read": "boolean",
      "created_at": "datetime",
      "read_at": "datetime (optional)"
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

### Get Notification by ID

Retrieves a notification by ID.

**URL**: `GET /api/v1/notifications/{id}`

**Response**: `200 OK`

```json
{
  "id": "uuid",
  "user_id": "uuid",
  "type": "task_assigned | task_status_update | task_due_soon | task_overdue | task_comment | task_mention",
  "title": "string",
  "message": "string",
  "resource_type": "string",
  "resource_id": "uuid",
  "is_read": "boolean",
  "created_at": "datetime",
  "read_at": "datetime (optional)"
}
```

**Error Responses**:

- `404 Not Found` - Notification not found
- `403 Forbidden` - You don't have permission to access this notification

### Mark Notification as Read

Marks a notification as read.

**URL**: `PUT /api/v1/notifications/{id}/read`

**Response**: `204 No Content`

**Error Responses**:

- `404 Not Found` - Notification not found
- `403 Forbidden` - You don't have permission to access this notification

### Mark All Notifications as Read

Marks all notifications for the authenticated user as read.

**URL**: `PUT /api/v1/notifications/read-all`

**Response**: `204 No Content`

### Delete Notification

Deletes a notification.

**URL**: `DELETE /api/v1/notifications/{id}`

**Response**: `204 No Content`

**Error Responses**:

- `404 Not Found` - Notification not found
- `403 Forbidden` - You don't have permission to access this notification

### Delete All Notifications

Deletes all notifications for the authenticated user.

**URL**: `DELETE /api/v1/notifications`

**Response**: `204 No Content`

## Data Models

### Notification

```json
{
  "id": "uuid",
  "user_id": "uuid",
  "type": "task_assigned | task_status_update | task_due_soon | task_overdue | task_comment | task_mention",
  "title": "string",
  "message": "string",
  "resource_type": "string",
  "resource_id": "uuid",
  "is_read": "boolean",
  "created_at": "datetime",
  "read_at": "datetime (optional)"
}
```

# Comment API Reference

The Comment API allows you to manage comments on tasks in the Taskodex product.

## Base URL

```
/api/v1/organizations/{org_id}/taskodex/comments
```

## Authentication

All endpoints require authentication using a JWT token. The token should be included in the `Authorization` header as a Bearer token.

```
Authorization: Bearer <token>
```

## Permissions

The following permissions are required to access the Comment API:

- `comment:read` - Required to read comments
- `comment:write` - Required to create and update comments
- `comment:delete` - Required to delete comments

## Endpoints

### Create Comment

Creates a new comment on a task.

**URL**: `POST /api/v1/organizations/{org_id}/taskodex/comments`

**Permissions**: `comment:write`

**Request Body**:

```json
{
  "task_id": "uuid",
  "content": "string"
}
```

**Response**: `201 Created`

```json
{
  "id": "uuid",
  "task_id": "uuid",
  "user_id": "uuid",
  "content": "string",
  "created_at": "datetime",
  "updated_at": "datetime",
  "user": {
    "id": "uuid",
    "email": "string",
    "first_name": "string",
    "last_name": "string"
  },
  "mentions": ["uuid"]
}
```

**Error Responses**:

- `400 Bad Request` - Invalid request body
- `404 Not Found` - Task not found

### Get Comment by ID

Retrieves a comment by ID.

**URL**: `GET /api/v1/organizations/{org_id}/taskodex/comments/{id}`

**Permissions**: `comment:read`

**Response**: `200 OK`

```json
{
  "id": "uuid",
  "task_id": "uuid",
  "user_id": "uuid",
  "content": "string",
  "created_at": "datetime",
  "updated_at": "datetime",
  "user": {
    "id": "uuid",
    "email": "string",
    "first_name": "string",
    "last_name": "string"
  },
  "mentions": ["uuid"]
}
```

**Error Responses**:

- `404 Not Found` - Comment not found

### List Comments

Retrieves a list of comments based on filter parameters.

**URL**: `GET /api/v1/organizations/{org_id}/taskodex/comments`

**Permissions**: `comment:read`

**Query Parameters**:

- `task_id` (optional) - Filter by task ID
- `user_id` (optional) - Filter by user ID
- `sort_by` (optional) - Sort by field (created_at, updated_at)
- `sort_order` (optional) - Sort order (asc/desc)
- `page` (optional) - Page number (default: 1)
- `page_size` (optional) - Page size (default: 20, max: 100)

**Response**: `200 OK`

```json
{
  "comments": [
    {
      "id": "uuid",
      "task_id": "uuid",
      "user_id": "uuid",
      "content": "string",
      "created_at": "datetime",
      "updated_at": "datetime",
      "user": {
        "id": "uuid",
        "email": "string",
        "first_name": "string",
        "last_name": "string"
      },
      "mentions": ["uuid"]
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

### Update Comment

Updates a comment.

**URL**: `PUT /api/v1/organizations/{org_id}/taskodex/comments/{id}`

**Permissions**: `comment:write`

**Request Body**:

```json
{
  "task_id": "uuid",
  "content": "string"
}
```

**Response**: `200 OK`

```json
{
  "id": "uuid",
  "task_id": "uuid",
  "user_id": "uuid",
  "content": "string",
  "created_at": "datetime",
  "updated_at": "datetime",
  "user": {
    "id": "uuid",
    "email": "string",
    "first_name": "string",
    "last_name": "string"
  },
  "mentions": ["uuid"]
}
```

**Error Responses**:

- `400 Bad Request` - Invalid request body
- `403 Forbidden` - You don't have permission to update this comment
- `404 Not Found` - Comment or task not found

### Delete Comment

Deletes a comment.

**URL**: `DELETE /api/v1/organizations/{org_id}/taskodex/comments/{id}`

**Permissions**: `comment:delete`

**Response**: `204 No Content`

**Error Responses**:

- `403 Forbidden` - You don't have permission to delete this comment
- `404 Not Found` - Comment not found

## Mentioning Users

To mention a user in a comment, include their username with an @ symbol in the comment content. For example:

```
@johndoe Please take a look at this task.
```

When a user is mentioned in a comment, they will receive a notification.

## Data Models

### Comment

```json
{
  "id": "uuid",
  "task_id": "uuid",
  "user_id": "uuid",
  "content": "string",
  "created_at": "datetime",
  "updated_at": "datetime",
  "user": {
    "id": "uuid",
    "email": "string",
    "first_name": "string",
    "last_name": "string"
  },
  "mentions": ["uuid"]
}
```

### CommentRequest

```json
{
  "task_id": "uuid",
  "content": "string"
}
```

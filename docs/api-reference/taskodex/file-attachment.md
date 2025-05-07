# File Attachment API Reference

The File Attachment API allows you to manage file attachments for tasks in the Taskodex product.

## Base URL

```
/api/v1/organizations/{org_id}/taskodex/files
```

## Authentication

All endpoints require authentication using a JWT token. The token should be included in the `Authorization` header as a Bearer token.

```
Authorization: Bearer <token>
```

## Permissions

The following permissions are required to access the File Attachment API:

- `file:read` - Required to read file attachments
- `file:write` - Required to upload file attachments
- `file:delete` - Required to delete file attachments

## Endpoints

### Upload File

Uploads a file and creates a file attachment for a task.

**URL**: `POST /api/v1/organizations/{org_id}/taskodex/files`

**Permissions**: `file:write`

**Content-Type**: `multipart/form-data`

**Form Parameters**:

- `task_id` - The ID of the task to attach the file to
- `file` - The file to upload

**Response**: `201 Created`

```json
{
  "id": "uuid",
  "task_id": "uuid",
  "user_id": "uuid",
  "file_name": "string",
  "file_size": "number",
  "content_type": "string",
  "download_url": "string",
  "created_at": "datetime",
  "user": {
    "id": "uuid",
    "email": "string",
    "first_name": "string",
    "last_name": "string"
  }
}
```

**Error Responses**:

- `400 Bad Request` - Invalid request body or file size exceeds the limit (10MB)
- `404 Not Found` - Task not found

### Get File Attachment by ID

Retrieves a file attachment by ID.

**URL**: `GET /api/v1/organizations/{org_id}/taskodex/files/{id}`

**Permissions**: `file:read`

**Response**: `200 OK`

```json
{
  "id": "uuid",
  "task_id": "uuid",
  "user_id": "uuid",
  "file_name": "string",
  "file_size": "number",
  "content_type": "string",
  "download_url": "string",
  "created_at": "datetime",
  "user": {
    "id": "uuid",
    "email": "string",
    "first_name": "string",
    "last_name": "string"
  }
}
```

**Error Responses**:

- `404 Not Found` - File attachment not found

### List File Attachments

Retrieves a list of file attachments based on filter parameters.

**URL**: `GET /api/v1/organizations/{org_id}/taskodex/files`

**Permissions**: `file:read`

**Query Parameters**:

- `task_id` (optional) - Filter by task ID
- `user_id` (optional) - Filter by user ID
- `sort_by` (optional) - Sort by field (created_at, file_name, file_size)
- `sort_order` (optional) - Sort order (asc/desc)
- `page` (optional) - Page number (default: 1)
- `page_size` (optional) - Page size (default: 20, max: 100)

**Response**: `200 OK`

```json
{
  "file_attachments": [
    {
      "id": "uuid",
      "task_id": "uuid",
      "user_id": "uuid",
      "file_name": "string",
      "file_size": "number",
      "content_type": "string",
      "download_url": "string",
      "created_at": "datetime",
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

### Download File

Downloads a file attachment.

**URL**: `GET /api/v1/organizations/{org_id}/taskodex/files/{id}/download`

**Permissions**: `file:read`

**Response**: The file content with appropriate Content-Type and Content-Disposition headers.

**Error Responses**:

- `404 Not Found` - File attachment not found

### Delete File Attachment

Deletes a file attachment.

**URL**: `DELETE /api/v1/organizations/{org_id}/taskodex/files/{id}`

**Permissions**: `file:delete`

**Response**: `204 No Content`

**Error Responses**:

- `403 Forbidden` - You don't have permission to delete this file attachment
- `404 Not Found` - File attachment not found

## File Size Limits

The maximum file size for uploads is 10MB.

## Supported File Types

All file types are supported, but the following common file types are recommended:

- Documents: PDF, DOC, DOCX, XLS, XLSX, PPT, PPTX, TXT
- Images: JPG, JPEG, PNG, GIF, SVG
- Archives: ZIP, RAR, 7Z
- Code: JS, TS, PY, GO, JAVA, HTML, CSS, JSON, XML

## Data Models

### FileAttachment

```json
{
  "id": "uuid",
  "task_id": "uuid",
  "user_id": "uuid",
  "file_name": "string",
  "file_size": "number",
  "content_type": "string",
  "download_url": "string",
  "created_at": "datetime",
  "user": {
    "id": "uuid",
    "email": "string",
    "first_name": "string",
    "last_name": "string"
  }
}
```

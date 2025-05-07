import type { UserResponse } from './user';

export type TaskStatus = 'todo' | 'in_progress' | 'review' | 'done' | 'cancelled';
export type TaskPriority = 'low' | 'medium' | 'high' | 'critical';

export interface Task {
  id: string;
  project_id: string;
  title: string;
  description: string;
  status: TaskStatus;
  priority: TaskPriority;
  due_date?: string;
  created_by: string;
  assigned_to?: string;
  estimated_hours?: number;
  actual_hours?: number;
  tags: string[];
  created_at: string;
  updated_at: string;
}

export interface TaskResponse extends Task {
  creator?: UserResponse;
  assignee?: UserResponse;
}

export interface TaskRequest {
  project_id: string;
  title: string;
  description: string;
  status: TaskStatus;
  priority: TaskPriority;
  due_date?: string;
  assigned_to?: string;
  estimated_hours?: number;
  tags?: string[];
}

export interface TaskListParams {
  status?: TaskStatus;
  priority?: TaskPriority;
  assignedTo?: string;
  projectId?: string;
  search?: string;
  sortBy?: 'title' | 'status' | 'priority' | 'due_date' | 'created_at' | 'updated_at';
  sortOrder?: 'asc' | 'desc';
  page?: number;
  pageSize?: number;
}

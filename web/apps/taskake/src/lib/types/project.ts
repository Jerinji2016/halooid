import type { UserResponse } from './user';

export type ProjectStatus = 'planning' | 'active' | 'on_hold' | 'completed' | 'cancelled';

export interface Project {
  id: string;
  organization_id: string;
  name: string;
  description: string;
  status: ProjectStatus;
  start_date?: string;
  end_date?: string;
  created_by: string;
  created_at: string;
  updated_at: string;
}

export interface ProjectResponse extends Project {
  creator?: UserResponse;
}

export interface ProjectRequest {
  organization_id: string;
  name: string;
  description: string;
  status: ProjectStatus;
  start_date?: string;
  end_date?: string;
}

export interface ProjectListParams {
  status?: ProjectStatus;
  createdBy?: string;
  search?: string;
  sortBy?: 'name' | 'status' | 'start_date' | 'end_date' | 'created_at' | 'updated_at';
  sortOrder?: 'asc' | 'desc';
  page?: number;
  pageSize?: number;
}

export interface ProjectStats {
  tasksByStatus: { status: string, count: number }[];
  tasksByPriority: { priority: string, count: number }[];
  tasksByAssignee: { assignee_id: string, assignee_name: string, count: number }[];
  completionPercentage: number;
  totalTasks: number;
  completedTasks: number;
  overdueTasks: number;
  upcomingDeadlines: any[]; // TaskResponse[]
}

export interface ProjectActivity {
  id: string;
  project_id: string;
  user_id: string;
  user_name: string;
  action: string;
  entity_type: string;
  entity_id: string;
  entity_name: string;
  created_at: string;
}

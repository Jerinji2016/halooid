import type { Project, ProjectListParams, ProjectResponse } from '$lib/types/project';
import type { TaskListParams, TaskResponse } from '$lib/types/task';
import { API_BASE_URL } from '$lib/config';

/**
 * Project API service for interacting with the project endpoints
 */
export class ProjectApi {
  private baseUrl: string;
  private orgId: string;

  constructor(orgId: string) {
    this.baseUrl = `${API_BASE_URL}/organizations/${orgId}/taskodex`;
    this.orgId = orgId;
  }

  /**
   * Get a project by ID
   * @param id Project ID
   * @returns Project data
   */
  async getProject(id: string): Promise<ProjectResponse> {
    const response = await fetch(`${this.baseUrl}/projects/${id}`, {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });

    if (!response.ok) {
      throw new Error(`Failed to fetch project: ${response.statusText}`);
    }

    return await response.json();
  }

  /**
   * Get a list of projects
   * @param params Filter parameters
   * @returns List of projects and pagination info
   */
  async listProjects(params: ProjectListParams = {}): Promise<{ projects: ProjectResponse[], pagination: { total: number, page: number, page_size: number, total_pages: number } }> {
    const queryParams = new URLSearchParams();
    
    if (params.status) queryParams.append('status', params.status);
    if (params.createdBy) queryParams.append('created_by', params.createdBy);
    if (params.search) queryParams.append('search', params.search);
    if (params.sortBy) queryParams.append('sort_by', params.sortBy);
    if (params.sortOrder) queryParams.append('sort_order', params.sortOrder);
    if (params.page) queryParams.append('page', params.page.toString());
    if (params.pageSize) queryParams.append('page_size', params.pageSize.toString());

    const response = await fetch(`${this.baseUrl}/projects?${queryParams.toString()}`, {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });

    if (!response.ok) {
      throw new Error(`Failed to fetch projects: ${response.statusText}`);
    }

    return await response.json();
  }

  /**
   * Get tasks for a project
   * @param projectId Project ID
   * @param params Filter parameters
   * @returns List of tasks and pagination info
   */
  async getProjectTasks(projectId: string, params: TaskListParams = {}): Promise<{ tasks: TaskResponse[], pagination: { total: number, page: number, page_size: number, total_pages: number } }> {
    const queryParams = new URLSearchParams();
    
    if (params.status) queryParams.append('status', params.status);
    if (params.priority) queryParams.append('priority', params.priority);
    if (params.assignedTo) queryParams.append('assigned_to', params.assignedTo);
    if (params.search) queryParams.append('search', params.search);
    if (params.sortBy) queryParams.append('sort_by', params.sortBy);
    if (params.sortOrder) queryParams.append('sort_order', params.sortOrder);
    if (params.page) queryParams.append('page', params.page.toString());
    if (params.pageSize) queryParams.append('page_size', params.pageSize.toString());

    const response = await fetch(`${this.baseUrl}/projects/${projectId}/tasks?${queryParams.toString()}`, {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });

    if (!response.ok) {
      throw new Error(`Failed to fetch project tasks: ${response.statusText}`);
    }

    return await response.json();
  }

  /**
   * Get project statistics
   * @param projectId Project ID
   * @returns Project statistics
   */
  async getProjectStats(projectId: string): Promise<{
    tasksByStatus: { status: string, count: number }[],
    tasksByPriority: { priority: string, count: number }[],
    tasksByAssignee: { assignee_id: string, assignee_name: string, count: number }[],
    completionPercentage: number,
    totalTasks: number,
    completedTasks: number,
    overdueTasks: number,
    upcomingDeadlines: TaskResponse[]
  }> {
    const response = await fetch(`${this.baseUrl}/projects/${projectId}/stats`, {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });

    if (!response.ok) {
      throw new Error(`Failed to fetch project statistics: ${response.statusText}`);
    }

    return await response.json();
  }

  /**
   * Get project activity
   * @param projectId Project ID
   * @param limit Number of activities to return
   * @returns Project activity
   */
  async getProjectActivity(projectId: string, limit: number = 10): Promise<{
    activities: {
      id: string,
      project_id: string,
      user_id: string,
      user_name: string,
      action: string,
      entity_type: string,
      entity_id: string,
      entity_name: string,
      created_at: string
    }[]
  }> {
    const response = await fetch(`${this.baseUrl}/projects/${projectId}/activity?limit=${limit}`, {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });

    if (!response.ok) {
      throw new Error(`Failed to fetch project activity: ${response.statusText}`);
    }

    return await response.json();
  }
}

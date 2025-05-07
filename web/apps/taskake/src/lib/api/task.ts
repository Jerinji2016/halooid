import type { TaskListParams, TaskResponse, TaskRequest } from '$lib/types/task';
import { API_BASE_URL } from '$lib/config';

/**
 * Task API service for interacting with the task endpoints
 */
export class TaskApi {
  private baseUrl: string;

  constructor(orgId: string) {
    this.baseUrl = `${API_BASE_URL}/organizations/${orgId}/taskodex`;
  }

  /**
   * Get a task by ID
   * @param id Task ID
   * @returns Task data
   */
  async getTask(id: string): Promise<TaskResponse> {
    const response = await fetch(`${this.baseUrl}/tasks/${id}`, {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });

    if (!response.ok) {
      throw new Error(`Failed to fetch task: ${response.statusText}`);
    }

    return await response.json();
  }

  /**
   * Get a list of tasks
   * @param params Filter parameters
   * @returns List of tasks and pagination info
   */
  async listTasks(params: TaskListParams = {}): Promise<{ tasks: TaskResponse[], pagination: { total: number, page: number, page_size: number, total_pages: number } }> {
    const queryParams = new URLSearchParams();

    if (params.status) queryParams.append('status', params.status);
    if (params.priority) queryParams.append('priority', params.priority);
    if (params.assignedTo) queryParams.append('assigned_to', params.assignedTo);
    if (params.projectId) queryParams.append('project_id', params.projectId);
    if (params.search) queryParams.append('search', params.search);
    if (params.sortBy) queryParams.append('sort_by', params.sortBy);
    if (params.sortOrder) queryParams.append('sort_order', params.sortOrder);
    if (params.page) queryParams.append('page', params.page.toString());
    if (params.pageSize) queryParams.append('page_size', params.pageSize.toString());

    const response = await fetch(`${this.baseUrl}/tasks?${queryParams.toString()}`, {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });

    if (!response.ok) {
      throw new Error(`Failed to fetch tasks: ${response.statusText}`);
    }

    return await response.json();
  }

  /**
   * Create a new task
   * @param task Task data
   * @returns Created task
   */
  async createTask(task: TaskRequest): Promise<TaskResponse> {
    const response = await fetch(`${this.baseUrl}/tasks`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify(task)
    });

    if (!response.ok) {
      throw new Error(`Failed to create task: ${response.statusText}`);
    }

    return await response.json();
  }

  /**
   * Update a task
   * @param id Task ID
   * @param task Task data
   * @returns Updated task
   */
  async updateTask(id: string, task: Partial<TaskRequest>): Promise<TaskResponse> {
    const response = await fetch(`${this.baseUrl}/tasks/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify(task)
    });

    if (!response.ok) {
      throw new Error(`Failed to update task: ${response.statusText}`);
    }

    return await response.json();
  }

  /**
   * Delete a task
   * @param id Task ID
   */
  async deleteTask(id: string): Promise<void> {
    const response = await fetch(`${this.baseUrl}/tasks/${id}`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });

    if (!response.ok) {
      throw new Error(`Failed to delete task: ${response.statusText}`);
    }
  }

  /**
   * Assign a task to a user
   * @param taskId Task ID
   * @param userId User ID
   * @returns Updated task
   */
  async assignTask(taskId: string, userId: string): Promise<TaskResponse> {
    const response = await fetch(`${this.baseUrl}/tasks/${taskId}/assign`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify({ user_id: userId })
    });

    if (!response.ok) {
      throw new Error(`Failed to assign task: ${response.statusText}`);
    }

    return await response.json();
  }

  /**
   * Update task status
   * @param taskId Task ID
   * @param status New status
   * @returns Updated task
   */
  async updateTaskStatus(taskId: string, status: string): Promise<TaskResponse> {
    const response = await fetch(`${this.baseUrl}/tasks/${taskId}/status`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify({ status })
    });

    if (!response.ok) {
      throw new Error(`Failed to update task status: ${response.statusText}`);
    }

    return await response.json();
  }
}

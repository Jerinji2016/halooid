import { writable, derived } from 'svelte/store';
import { DEFAULT_ORG_ID } from '$lib/config';
import type { Task, TaskResponse, TaskListParams } from '$lib/types/task';

// Define store types
interface TaskStore {
  items: TaskResponse[];
  pagination: {
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
  };
}

// Create writable stores
const tasks = writable<TaskStore>({
  items: [],
  pagination: {
    total: 0,
    page: 1,
    page_size: 20,
    total_pages: 1
  }
});

const currentTask = writable<TaskResponse | null>(null);
const isLoading = writable<boolean>(false);
const error = writable<string | null>(null);

// API base URL
const API_BASE_URL = import.meta.env.VITE_API_URL || 'https://api.halooid.com/v1';

// Task API functions
const fetchTasks = async (params: TaskListParams = {}) => {
  isLoading.set(true);
  error.set(null);
  
  try {
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
    
    const response = await fetch(`${API_BASE_URL}/organizations/${DEFAULT_ORG_ID}/taskodex/tasks?${queryParams.toString()}`, {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });
    
    if (!response.ok) {
      throw new Error(`Failed to fetch tasks: ${response.statusText}`);
    }
    
    const data = await response.json();
    
    tasks.update(store => ({
      items: data.tasks || [],
      pagination: {
        total: data.pagination?.total || 0,
        page: data.pagination?.page || 1,
        page_size: data.pagination?.page_size || 20,
        total_pages: data.pagination?.total_pages || 1
      }
    }));
    
    isLoading.set(false);
    return data;
  } catch (err) {
    console.error('Error fetching tasks:', err);
    error.set(err.message);
    isLoading.set(false);
    throw err;
  }
};

const fetchTask = async (id: string) => {
  isLoading.set(true);
  error.set(null);
  
  try {
    const response = await fetch(`${API_BASE_URL}/organizations/${DEFAULT_ORG_ID}/taskodex/tasks/${id}`, {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });
    
    if (!response.ok) {
      throw new Error(`Failed to fetch task: ${response.statusText}`);
    }
    
    const data = await response.json();
    currentTask.set(data);
    isLoading.set(false);
    return data;
  } catch (err) {
    console.error('Error fetching task:', err);
    error.set(err.message);
    isLoading.set(false);
    throw err;
  }
};

const createTask = async (task: Partial<Task>) => {
  isLoading.set(true);
  error.set(null);
  
  try {
    const response = await fetch(`${API_BASE_URL}/organizations/${DEFAULT_ORG_ID}/taskodex/tasks`, {
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
    
    const data = await response.json();
    isLoading.set(false);
    return data;
  } catch (err) {
    console.error('Error creating task:', err);
    error.set(err.message);
    isLoading.set(false);
    throw err;
  }
};

const updateTask = async (id: string, task: Partial<Task>) => {
  isLoading.set(true);
  error.set(null);
  
  try {
    const response = await fetch(`${API_BASE_URL}/organizations/${DEFAULT_ORG_ID}/taskodex/tasks/${id}`, {
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
    
    const data = await response.json();
    
    // Update current task if it's the one being edited
    currentTask.update(current => {
      if (current && current.id === id) {
        return { ...current, ...data };
      }
      return current;
    });
    
    // Update task in the list if it exists
    tasks.update(store => {
      const updatedItems = store.items.map(item => {
        if (item.id === id) {
          return { ...item, ...data };
        }
        return item;
      });
      
      return {
        ...store,
        items: updatedItems
      };
    });
    
    isLoading.set(false);
    return data;
  } catch (err) {
    console.error('Error updating task:', err);
    error.set(err.message);
    isLoading.set(false);
    throw err;
  }
};

const deleteTask = async (id: string) => {
  isLoading.set(true);
  error.set(null);
  
  try {
    const response = await fetch(`${API_BASE_URL}/organizations/${DEFAULT_ORG_ID}/taskodex/tasks/${id}`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });
    
    if (!response.ok) {
      throw new Error(`Failed to delete task: ${response.statusText}`);
    }
    
    // Remove task from the list
    tasks.update(store => {
      const filteredItems = store.items.filter(item => item.id !== id);
      
      return {
        ...store,
        items: filteredItems,
        pagination: {
          ...store.pagination,
          total: store.pagination.total - 1
        }
      };
    });
    
    // Clear current task if it's the one being deleted
    currentTask.update(current => {
      if (current && current.id === id) {
        return null;
      }
      return current;
    });
    
    isLoading.set(false);
    return true;
  } catch (err) {
    console.error('Error deleting task:', err);
    error.set(err.message);
    isLoading.set(false);
    throw err;
  }
};

// Derived stores
const tasksByStatus = derived(tasks, $tasks => {
  const statusGroups: Record<string, TaskResponse[]> = {};
  
  $tasks.items.forEach(task => {
    if (!statusGroups[task.status]) {
      statusGroups[task.status] = [];
    }
    statusGroups[task.status].push(task);
  });
  
  return statusGroups;
});

const tasksByPriority = derived(tasks, $tasks => {
  const priorityGroups: Record<string, TaskResponse[]> = {};
  
  $tasks.items.forEach(task => {
    if (!priorityGroups[task.priority]) {
      priorityGroups[task.priority] = [];
    }
    priorityGroups[task.priority].push(task);
  });
  
  return priorityGroups;
});

const tasksByAssignee = derived(tasks, $tasks => {
  const assigneeGroups: Record<string, TaskResponse[]> = {};
  
  $tasks.items.forEach(task => {
    const assigneeId = task.assigned_to || 'unassigned';
    if (!assigneeGroups[assigneeId]) {
      assigneeGroups[assigneeId] = [];
    }
    assigneeGroups[assigneeId].push(task);
  });
  
  return assigneeGroups;
});

// Export the store
export const taskStore = {
  subscribe: tasks.subscribe,
  currentTask: {
    subscribe: currentTask.subscribe
  },
  isLoading: {
    subscribe: isLoading.subscribe
  },
  error: {
    subscribe: error.subscribe
  },
  tasksByStatus: {
    subscribe: tasksByStatus.subscribe
  },
  tasksByPriority: {
    subscribe: tasksByPriority.subscribe
  },
  tasksByAssignee: {
    subscribe: tasksByAssignee.subscribe
  },
  fetchTasks,
  fetchTask,
  createTask,
  updateTask,
  deleteTask
};

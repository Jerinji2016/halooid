import { writable, derived } from 'svelte/store';
import type { ProjectResponse, ProjectStats, ProjectActivity } from '$lib/types/project';
import type { TaskResponse } from '$lib/types/task';
import { ProjectApi } from '$lib/api/project';

// Create writable stores
export const currentProject = writable<ProjectResponse | null>(null);
export const projectTasks = writable<TaskResponse[]>([]);
export const projectStats = writable<ProjectStats | null>(null);
export const projectActivity = writable<ProjectActivity[]>([]);
export const isLoading = writable<boolean>(false);
export const error = writable<string | null>(null);

// Derived stores
export const tasksByStatus = derived(projectTasks, ($projectTasks) => {
  const statusGroups: Record<string, TaskResponse[]> = {};
  
  $projectTasks.forEach(task => {
    if (!statusGroups[task.status]) {
      statusGroups[task.status] = [];
    }
    statusGroups[task.status].push(task);
  });
  
  return statusGroups;
});

export const tasksByPriority = derived(projectTasks, ($projectTasks) => {
  const priorityGroups: Record<string, TaskResponse[]> = {};
  
  $projectTasks.forEach(task => {
    if (!priorityGroups[task.priority]) {
      priorityGroups[task.priority] = [];
    }
    priorityGroups[task.priority].push(task);
  });
  
  return priorityGroups;
});

export const tasksByAssignee = derived(projectTasks, ($projectTasks) => {
  const assigneeGroups: Record<string, TaskResponse[]> = {};
  
  $projectTasks.forEach(task => {
    const assigneeId = task.assigned_to || 'unassigned';
    if (!assigneeGroups[assigneeId]) {
      assigneeGroups[assigneeId] = [];
    }
    assigneeGroups[assigneeId].push(task);
  });
  
  return assigneeGroups;
});

export const completionPercentage = derived(projectTasks, ($projectTasks) => {
  if ($projectTasks.length === 0) return 0;
  
  const completedTasks = $projectTasks.filter(task => task.status === 'done').length;
  return Math.round((completedTasks / $projectTasks.length) * 100);
});

export const overdueTasks = derived(projectTasks, ($projectTasks) => {
  const now = new Date();
  return $projectTasks.filter(task => {
    if (!task.due_date) return false;
    const dueDate = new Date(task.due_date);
    return dueDate < now && task.status !== 'done' && task.status !== 'cancelled';
  });
});

export const upcomingDeadlines = derived(projectTasks, ($projectTasks) => {
  const now = new Date();
  const oneWeekFromNow = new Date();
  oneWeekFromNow.setDate(now.getDate() + 7);
  
  return $projectTasks
    .filter(task => {
      if (!task.due_date) return false;
      const dueDate = new Date(task.due_date);
      return dueDate >= now && dueDate <= oneWeekFromNow && task.status !== 'done' && task.status !== 'cancelled';
    })
    .sort((a, b) => {
      if (!a.due_date || !b.due_date) return 0;
      return new Date(a.due_date).getTime() - new Date(b.due_date).getTime();
    });
});

// Actions
export const projectActions = (orgId: string) => {
  const projectApi = new ProjectApi(orgId);
  
  return {
    /**
     * Load project data
     * @param projectId Project ID
     */
    async loadProject(projectId: string) {
      isLoading.set(true);
      error.set(null);
      
      try {
        const project = await projectApi.getProject(projectId);
        currentProject.set(project);
      } catch (err) {
        error.set(err instanceof Error ? err.message : 'Failed to load project');
        console.error(err);
      } finally {
        isLoading.set(false);
      }
    },
    
    /**
     * Load project tasks
     * @param projectId Project ID
     */
    async loadProjectTasks(projectId: string) {
      isLoading.set(true);
      error.set(null);
      
      try {
        const response = await projectApi.getProjectTasks(projectId, { pageSize: 100 });
        projectTasks.set(response.tasks);
      } catch (err) {
        error.set(err instanceof Error ? err.message : 'Failed to load project tasks');
        console.error(err);
      } finally {
        isLoading.set(false);
      }
    },
    
    /**
     * Load project statistics
     * @param projectId Project ID
     */
    async loadProjectStats(projectId: string) {
      isLoading.set(true);
      error.set(null);
      
      try {
        const stats = await projectApi.getProjectStats(projectId);
        projectStats.set(stats);
      } catch (err) {
        error.set(err instanceof Error ? err.message : 'Failed to load project statistics');
        console.error(err);
      } finally {
        isLoading.set(false);
      }
    },
    
    /**
     * Load project activity
     * @param projectId Project ID
     */
    async loadProjectActivity(projectId: string) {
      isLoading.set(true);
      error.set(null);
      
      try {
        const response = await projectApi.getProjectActivity(projectId);
        projectActivity.set(response.activities);
      } catch (err) {
        error.set(err instanceof Error ? err.message : 'Failed to load project activity');
        console.error(err);
      } finally {
        isLoading.set(false);
      }
    },
    
    /**
     * Load all project data
     * @param projectId Project ID
     */
    async loadAllProjectData(projectId: string) {
      isLoading.set(true);
      error.set(null);
      
      try {
        await Promise.all([
          this.loadProject(projectId),
          this.loadProjectTasks(projectId),
          this.loadProjectStats(projectId),
          this.loadProjectActivity(projectId)
        ]);
      } catch (err) {
        error.set(err instanceof Error ? err.message : 'Failed to load project data');
        console.error(err);
      } finally {
        isLoading.set(false);
      }
    }
  };
};

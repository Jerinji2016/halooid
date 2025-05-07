<script lang="ts">
  import { DATE_FORMAT_OPTIONS, PRIORITY_COLORS } from '$lib/config';
  import type { TaskResponse } from '$lib/types/task';
  
  export let tasks: TaskResponse[];
  
  // Format date for display
  const formatDate = (dateString: string | undefined) => {
    if (!dateString) return 'No due date';
    return new Date(dateString).toLocaleDateString(undefined, DATE_FORMAT_OPTIONS);
  };
  
  // Calculate days remaining
  const getDaysRemaining = (dateString: string | undefined) => {
    if (!dateString) return '';
    
    const dueDate = new Date(dateString);
    const today = new Date();
    
    // Reset time part for accurate day calculation
    dueDate.setHours(0, 0, 0, 0);
    today.setHours(0, 0, 0, 0);
    
    const diffTime = dueDate.getTime() - today.getTime();
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
    
    if (diffDays < 0) {
      return `${Math.abs(diffDays)} days overdue`;
    } else if (diffDays === 0) {
      return 'Due today';
    } else if (diffDays === 1) {
      return '1 day remaining';
    } else {
      return `${diffDays} days remaining`;
    }
  };
  
  // Get days remaining class
  const getDaysRemainingClass = (dateString: string | undefined) => {
    if (!dateString) return '';
    
    const dueDate = new Date(dateString);
    const today = new Date();
    
    // Reset time part for accurate day calculation
    dueDate.setHours(0, 0, 0, 0);
    today.setHours(0, 0, 0, 0);
    
    const diffTime = dueDate.getTime() - today.getTime();
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
    
    if (diffDays < 0) {
      return 'text-red-600 font-medium';
    } else if (diffDays === 0) {
      return 'text-orange-600 font-medium';
    } else if (diffDays <= 2) {
      return 'text-yellow-600';
    } else {
      return 'text-green-600';
    }
  };
  
  // Get priority color
  const getPriorityColor = (priority: string) => {
    return PRIORITY_COLORS[priority as keyof typeof PRIORITY_COLORS] || '#9E9E9E';
  };
  
  // Sort tasks by due date
  const sortedTasks = tasks
    .filter(task => task.due_date)
    .sort((a, b) => {
      if (!a.due_date || !b.due_date) return 0;
      return new Date(a.due_date).getTime() - new Date(b.due_date).getTime();
    });
</script>

<div class="upcoming-deadlines bg-white rounded-lg shadow-md p-6">
  <h2 class="text-lg font-semibold mb-4">Upcoming Deadlines</h2>
  
  {#if sortedTasks.length === 0}
    <div class="flex items-center justify-center h-32 text-gray-500">
      No upcoming deadlines
    </div>
  {:else}
    <div class="deadlines-list space-y-4">
      {#each sortedTasks.slice(0, 5) as task}
        <div class="deadline-item p-3 border rounded-lg hover:bg-gray-50">
          <div class="flex justify-between items-start">
            <div>
              <a 
                href="/tasks/{task.id}" 
                class="text-blue-600 hover:text-blue-800 font-medium"
              >
                {task.title}
              </a>
              
              <div class="text-sm text-gray-600 mt-1">
                Assigned to: {task.assignee?.first_name || 'Unassigned'}
              </div>
            </div>
            
            <div 
              class="priority-indicator w-3 h-3 rounded-full mt-1"
              style="background-color: {getPriorityColor(task.priority)};"
              title="Priority: {task.priority.charAt(0).toUpperCase() + task.priority.slice(1)}"
            ></div>
          </div>
          
          <div class="flex justify-between mt-2 text-sm">
            <div>
              <span class="text-gray-500">Due:</span> {formatDate(task.due_date)}
            </div>
            
            <div class={getDaysRemainingClass(task.due_date)}>
              {getDaysRemaining(task.due_date)}
            </div>
          </div>
        </div>
      {/each}
      
      {#if sortedTasks.length > 5}
        <div class="mt-2 text-center">
          <a href="/projects/{sortedTasks[0].project_id}/tasks" class="text-blue-500 hover:text-blue-700 text-sm">
            View all tasks
          </a>
        </div>
      {/if}
    </div>
  {/if}
</div>

<style>
  /* Add any component-specific styles here */
</style>

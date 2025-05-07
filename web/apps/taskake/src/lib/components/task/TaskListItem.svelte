<script lang="ts">
  import { STATUS_COLORS, PRIORITY_COLORS, DATE_FORMAT_OPTIONS } from '$lib/config';
  import type { TaskResponse } from '$lib/types/task';
  import TaskStatusBadge from './TaskStatusBadge.svelte';
  import TaskPriorityBadge from './TaskPriorityBadge.svelte';
  import { Card } from '@halooid/ui-components';
  
  export let task: TaskResponse;
  
  // Format date for display
  const formatDate = (dateString) => {
    if (!dateString) return '';
    return new Date(dateString).toLocaleDateString(undefined, DATE_FORMAT_OPTIONS);
  };
  
  // Calculate days remaining or overdue
  const getDaysRemaining = (dueDate) => {
    if (!dueDate) return '';
    
    const due = new Date(dueDate);
    const today = new Date();
    
    // Reset time part for accurate day calculation
    due.setHours(0, 0, 0, 0);
    today.setHours(0, 0, 0, 0);
    
    const diffTime = due.getTime() - today.getTime();
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
  const getDaysRemainingClass = (dueDate) => {
    if (!dueDate) return '';
    
    const due = new Date(dueDate);
    const today = new Date();
    
    // Reset time part for accurate day calculation
    due.setHours(0, 0, 0, 0);
    today.setHours(0, 0, 0, 0);
    
    const diffTime = due.getTime() - today.getTime();
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
</script>

<Card hoverable>
  <a href="/tasks/{task.id}" class="block p-4 hover:bg-gray-50 rounded-lg transition-colors">
    <div class="flex flex-col md:flex-row md:items-center md:justify-between">
      <!-- Task Title and Status -->
      <div class="flex-1">
        <div class="flex items-start">
          <h3 class="text-lg font-medium text-gray-900 mr-2 {task.status === 'done' ? 'line-through' : ''}">
            {task.title}
          </h3>
          
          <div class="flex flex-wrap gap-2">
            <TaskStatusBadge status={task.status} size="sm" />
            <TaskPriorityBadge priority={task.priority} size="sm" />
          </div>
        </div>
        
        {#if task.description}
          <p class="mt-1 text-sm text-gray-600 line-clamp-2">
            {task.description}
          </p>
        {/if}
        
        {#if task.tags && task.tags.length > 0}
          <div class="mt-2 flex flex-wrap gap-1">
            {#each task.tags as tag}
              <span class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-gray-100 text-gray-800">
                {tag}
              </span>
            {/each}
          </div>
        {/if}
      </div>
      
      <!-- Task Metadata -->
      <div class="mt-4 md:mt-0 md:ml-6 flex flex-col items-end">
        {#if task.due_date}
          <div class="flex items-center">
            <svg class="w-4 h-4 text-gray-400 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
            </svg>
            <span class="text-sm text-gray-500">{formatDate(task.due_date)}</span>
          </div>
          
          <div class="mt-1 text-xs {getDaysRemainingClass(task.due_date)}">
            {getDaysRemaining(task.due_date)}
          </div>
        {/if}
        
        {#if task.assignee}
          <div class="mt-2 flex items-center">
            <svg class="w-4 h-4 text-gray-400 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
            </svg>
            <span class="text-sm text-gray-500">{task.assignee.first_name} {task.assignee.last_name}</span>
          </div>
        {/if}
        
        {#if task.project_id}
          <div class="mt-2 flex items-center">
            <svg class="w-4 h-4 text-gray-400 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"></path>
            </svg>
            <a href="/projects/{task.project_id}" class="text-sm text-blue-600 hover:text-blue-800">
              {task.project_name || 'View Project'}
            </a>
          </div>
        {/if}
      </div>
    </div>
  </a>
</Card>

<style>
  /* Add any component-specific styles here */
  .line-clamp-2 {
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }
</style>

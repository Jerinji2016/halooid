<script lang="ts">
  import { STATUS_COLORS } from '$lib/config';
  import type { TaskStatus } from '$lib/types/task';
  
  export let status: TaskStatus;
  export let size: 'sm' | 'md' | 'lg' = 'md';
  
  // Get status color
  const getStatusColor = (status: TaskStatus) => {
    return STATUS_COLORS[status] || '#9E9E9E';
  };
  
  // Get status text
  const getStatusText = (status: TaskStatus) => {
    switch (status) {
      case 'todo':
        return 'To Do';
      case 'in_progress':
        return 'In Progress';
      case 'review':
        return 'Review';
      case 'done':
        return 'Done';
      case 'cancelled':
        return 'Cancelled';
      default:
        return status.charAt(0).toUpperCase() + status.slice(1).replace('_', ' ');
    }
  };
  
  // Get status icon
  const getStatusIcon = (status: TaskStatus) => {
    switch (status) {
      case 'todo':
        return 'M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2';
      case 'in_progress':
        return 'M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z';
      case 'review':
        return 'M15 12a3 3 0 11-6 0 3 3 0 016 0z M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z';
      case 'done':
        return 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z';
      case 'cancelled':
        return 'M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z';
      default:
        return 'M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z';
    }
  };
  
  // Get size classes
  const getSizeClasses = (size: 'sm' | 'md' | 'lg') => {
    switch (size) {
      case 'sm':
        return 'px-2 py-1 text-xs';
      case 'lg':
        return 'px-4 py-2 text-base';
      case 'md':
      default:
        return 'px-3 py-1 text-sm';
    }
  };
</script>

<span 
  class="inline-flex items-center rounded-full font-medium {getSizeClasses(size)}"
  style="background-color: {getStatusColor(status)}25; color: {getStatusColor(status)};"
>
  <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={getStatusIcon(status)}></path>
  </svg>
  {getStatusText(status)}
</span>

<style>
  /* Add any component-specific styles here */
</style>

<script lang="ts">
  import { PRIORITY_COLORS } from '$lib/config';
  import type { TaskPriority } from '$lib/types/task';
  
  export let priority: TaskPriority;
  export let size: 'sm' | 'md' | 'lg' = 'md';
  
  // Get priority color
  const getPriorityColor = (priority: TaskPriority) => {
    return PRIORITY_COLORS[priority] || '#9E9E9E';
  };
  
  // Get priority text
  const getPriorityText = (priority: TaskPriority) => {
    return priority.charAt(0).toUpperCase() + priority.slice(1);
  };
  
  // Get priority icon
  const getPriorityIcon = (priority: TaskPriority) => {
    switch (priority) {
      case 'low':
        return 'M19 14l-7 7m0 0l-7-7m7 7V3';
      case 'medium':
        return 'M5 12h14';
      case 'high':
        return 'M5 10l7-7m0 0l7 7m-7-7v18';
      case 'critical':
        return 'M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z';
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
  style="background-color: {getPriorityColor(priority)}25; color: {getPriorityColor(priority)};"
>
  <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={getPriorityIcon(priority)}></path>
  </svg>
  {getPriorityText(priority)}
</span>

<style>
  /* Add any component-specific styles here */
</style>

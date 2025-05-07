<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { STATUS_COLORS, PRIORITY_COLORS } from '$lib/config';
  import { userStore } from '$lib/stores/userStore';
  import type { TaskStatus, TaskPriority } from '$lib/types/task';
  import { Card } from '@halooid/ui-components';
  
  export let statusFilter: TaskStatus | null = null;
  export let priorityFilter: TaskPriority | null = null;
  export let assigneeFilter: string | null = null;
  
  const dispatch = createEventDispatcher();
  
  // Handle status filter change
  const handleStatusChange = (status: TaskStatus | null) => {
    statusFilter = status;
    dispatch('filterChange', { name: 'status', value: status });
  };
  
  // Handle priority filter change
  const handlePriorityChange = (priority: TaskPriority | null) => {
    priorityFilter = priority;
    dispatch('filterChange', { name: 'priority', value: priority });
  };
  
  // Handle assignee filter change
  const handleAssigneeChange = (assigneeId: string | null) => {
    assigneeFilter = assigneeId;
    dispatch('filterChange', { name: 'assignee', value: assigneeId });
  };
  
  // Clear all filters
  const clearFilters = () => {
    statusFilter = null;
    priorityFilter = null;
    assigneeFilter = null;
    dispatch('clearFilters');
  };
  
  // Status options
  const statusOptions = [
    { value: 'todo', label: 'To Do', color: STATUS_COLORS.todo },
    { value: 'in_progress', label: 'In Progress', color: STATUS_COLORS.in_progress },
    { value: 'review', label: 'Review', color: STATUS_COLORS.review },
    { value: 'done', label: 'Done', color: STATUS_COLORS.done },
    { value: 'cancelled', label: 'Cancelled', color: STATUS_COLORS.cancelled }
  ];
  
  // Priority options
  const priorityOptions = [
    { value: 'low', label: 'Low', color: PRIORITY_COLORS.low },
    { value: 'medium', label: 'Medium', color: PRIORITY_COLORS.medium },
    { value: 'high', label: 'High', color: PRIORITY_COLORS.high },
    { value: 'critical', label: 'Critical', color: PRIORITY_COLORS.critical }
  ];
</script>

<Card>
  <div class="p-4">
    <div class="flex justify-between items-center mb-4">
      <h2 class="text-lg font-medium">Filters</h2>
      
      <button 
        class="text-sm text-blue-600 hover:text-blue-800"
        on:click={clearFilters}
      >
        Clear All
      </button>
    </div>
    
    <!-- Status Filter -->
    <div class="mb-6">
      <h3 class="text-sm font-medium text-gray-700 mb-2">Status</h3>
      
      <div class="space-y-2">
        <label class="flex items-center">
          <input 
            type="radio" 
            name="status" 
            value="" 
            checked={statusFilter === null}
            on:change={() => handleStatusChange(null)}
            class="mr-2"
          />
          <span>All</span>
        </label>
        
        {#each statusOptions as option}
          <label class="flex items-center">
            <input 
              type="radio" 
              name="status" 
              value={option.value} 
              checked={statusFilter === option.value}
              on:change={() => handleStatusChange(option.value as TaskStatus)}
              class="mr-2"
            />
            <span 
              class="inline-block w-3 h-3 rounded-full mr-2"
              style="background-color: {option.color};"
            ></span>
            <span>{option.label}</span>
          </label>
        {/each}
      </div>
    </div>
    
    <!-- Priority Filter -->
    <div class="mb-6">
      <h3 class="text-sm font-medium text-gray-700 mb-2">Priority</h3>
      
      <div class="space-y-2">
        <label class="flex items-center">
          <input 
            type="radio" 
            name="priority" 
            value="" 
            checked={priorityFilter === null}
            on:change={() => handlePriorityChange(null)}
            class="mr-2"
          />
          <span>All</span>
        </label>
        
        {#each priorityOptions as option}
          <label class="flex items-center">
            <input 
              type="radio" 
              name="priority" 
              value={option.value} 
              checked={priorityFilter === option.value}
              on:change={() => handlePriorityChange(option.value as TaskPriority)}
              class="mr-2"
            />
            <span 
              class="inline-block w-3 h-3 rounded-full mr-2"
              style="background-color: {option.color};"
            ></span>
            <span>{option.label}</span>
          </label>
        {/each}
      </div>
    </div>
    
    <!-- Assignee Filter -->
    <div>
      <h3 class="text-sm font-medium text-gray-700 mb-2">Assignee</h3>
      
      <div class="space-y-2">
        <label class="flex items-center">
          <input 
            type="radio" 
            name="assignee" 
            value="" 
            checked={assigneeFilter === null}
            on:change={() => handleAssigneeChange(null)}
            class="mr-2"
          />
          <span>All</span>
        </label>
        
        <label class="flex items-center">
          <input 
            type="radio" 
            name="assignee" 
            value="me" 
            checked={assigneeFilter === 'me'}
            on:change={() => handleAssigneeChange('me')}
            class="mr-2"
          />
          <span>Assigned to me</span>
        </label>
        
        <label class="flex items-center">
          <input 
            type="radio" 
            name="assignee" 
            value="unassigned" 
            checked={assigneeFilter === 'unassigned'}
            on:change={() => handleAssigneeChange('unassigned')}
            class="mr-2"
          />
          <span>Unassigned</span>
        </label>
      </div>
    </div>
  </div>
</Card>

<style>
  /* Add any component-specific styles here */
</style>

<script lang="ts">
  import { onMount } from 'svelte';
  import { DEFAULT_ORG_ID } from '$lib/config';
  import { taskStore } from '$lib/stores/taskStore';
  import { userStore } from '$lib/stores/userStore';
  import { STATUS_COLORS, PRIORITY_COLORS } from '$lib/config';
  import TaskListItem from '$lib/components/task/TaskListItem.svelte';
  import TaskFilterPanel from '$lib/components/task/TaskFilterPanel.svelte';
  import Pagination from '$lib/components/common/Pagination.svelte';
  import Spinner from '$lib/components/common/Spinner.svelte';
  import Alert from '$lib/components/common/Alert.svelte';
  import { Button, Input, Card } from '@halooid/ui-components';
  
  // Local state
  let searchQuery = '';
  let isLoading = true;
  let error = null;
  let currentPage = 1;
  let pageSize = 20;
  let totalTasks = 0;
  let totalPages = 1;
  
  // Filter state
  let statusFilter = null;
  let priorityFilter = null;
  let assigneeFilter = null;
  let sortBy = 'due_date';
  let sortOrder = 'asc';
  
  // Initialize task store
  const { fetchTasks, tasks, isLoading: storeLoading, error: storeError } = taskStore;
  
  // Load tasks on mount
  onMount(async () => {
    try {
      isLoading = true;
      await fetchTasks({
        page: currentPage,
        pageSize,
        status: statusFilter,
        priority: priorityFilter,
        assignedTo: assigneeFilter,
        sortBy,
        sortOrder,
        search: searchQuery
      });
      
      // Get total from response
      totalTasks = $tasks.pagination?.total || 0;
      totalPages = $tasks.pagination?.total_pages || 1;
      
      isLoading = false;
    } catch (err) {
      error = err.message;
      isLoading = false;
    }
  });
  
  // Handle search
  const handleSearch = async () => {
    currentPage = 1; // Reset to first page on new search
    await fetchTasks({
      page: currentPage,
      pageSize,
      status: statusFilter,
      priority: priorityFilter,
      assignedTo: assigneeFilter,
      sortBy,
      sortOrder,
      search: searchQuery
    });
  };
  
  // Handle filter changes
  const handleFilterChange = async (event) => {
    const { name, value } = event.detail;
    
    if (name === 'status') {
      statusFilter = value;
    } else if (name === 'priority') {
      priorityFilter = value;
    } else if (name === 'assignee') {
      assigneeFilter = value;
    }
    
    currentPage = 1; // Reset to first page on filter change
    await fetchTasks({
      page: currentPage,
      pageSize,
      status: statusFilter,
      priority: priorityFilter,
      assignedTo: assigneeFilter,
      sortBy,
      sortOrder,
      search: searchQuery
    });
  };
  
  // Handle sort changes
  const handleSortChange = async (event) => {
    const { field, order } = event.detail;
    sortBy = field;
    sortOrder = order;
    
    await fetchTasks({
      page: currentPage,
      pageSize,
      status: statusFilter,
      priority: priorityFilter,
      assignedTo: assigneeFilter,
      sortBy,
      sortOrder,
      search: searchQuery
    });
  };
  
  // Handle page change
  const handlePageChange = async (event) => {
    currentPage = event.detail.page;
    
    await fetchTasks({
      page: currentPage,
      pageSize,
      status: statusFilter,
      priority: priorityFilter,
      assignedTo: assigneeFilter,
      sortBy,
      sortOrder,
      search: searchQuery
    });
  };
  
  // Clear all filters
  const clearFilters = async () => {
    statusFilter = null;
    priorityFilter = null;
    assigneeFilter = null;
    searchQuery = '';
    sortBy = 'due_date';
    sortOrder = 'asc';
    currentPage = 1;
    
    await fetchTasks({
      page: currentPage,
      pageSize
    });
  };
</script>

<svelte:head>
  <title>Tasks | Taskake</title>
</svelte:head>

<div class="task-list-page">
  <div class="flex justify-between items-center mb-6">
    <h1 class="text-2xl font-bold">Tasks</h1>
    
    <div class="flex space-x-2">
      <a href="/tasks/create" class="btn btn-primary">
        <svg class="w-5 h-5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
        </svg>
        New Task
      </a>
    </div>
  </div>
  
  <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
    <!-- Filter Panel -->
    <div class="md:col-span-1">
      <TaskFilterPanel 
        {statusFilter}
        {priorityFilter}
        {assigneeFilter}
        on:filterChange={handleFilterChange}
        on:clearFilters={clearFilters}
      />
    </div>
    
    <!-- Task List -->
    <div class="md:col-span-3">
      <!-- Search Bar -->
      <div class="mb-4 flex">
        <div class="flex-1 mr-2">
          <Input
            type="text"
            placeholder="Search tasks..."
            value={searchQuery}
            on:input={(e) => searchQuery = e.target.value}
            on:keydown={(e) => e.key === 'Enter' && handleSearch()}
            fullWidth
          />
        </div>
        <Button variant="primary" on:click={handleSearch}>
          Search
        </Button>
      </div>
      
      <!-- Sort Controls -->
      <div class="mb-4 flex justify-between items-center">
        <div class="text-sm text-gray-500">
          {#if totalTasks > 0}
            Showing {(currentPage - 1) * pageSize + 1} to {Math.min(currentPage * pageSize, totalTasks)} of {totalTasks} tasks
          {:else}
            No tasks found
          {/if}
        </div>
        
        <div class="flex items-center space-x-2">
          <span class="text-sm text-gray-500">Sort by:</span>
          <select 
            class="select select-sm select-bordered"
            value={sortBy}
            on:change={(e) => handleSortChange({ detail: { field: e.target.value, order: sortOrder } })}
          >
            <option value="due_date">Due Date</option>
            <option value="title">Title</option>
            <option value="status">Status</option>
            <option value="priority">Priority</option>
            <option value="created_at">Created Date</option>
          </select>
          
          <button 
            class="btn btn-sm btn-ghost"
            on:click={() => handleSortChange({ detail: { field: sortBy, order: sortOrder === 'asc' ? 'desc' : 'asc' } })}
          >
            {#if sortOrder === 'asc'}
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4h13M3 8h9m-9 4h6m4 0l4-4m0 0l4 4m-4-4v12"></path>
              </svg>
            {:else}
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4h13M3 8h9m-9 4h9m5-4v12m0 0l-4-4m4 4l4-4"></path>
              </svg>
            {/if}
          </button>
        </div>
      </div>
      
      <!-- Task List -->
      {#if isLoading || $storeLoading}
        <div class="flex justify-center items-center h-64">
          <Spinner size="lg" />
        </div>
      {:else if error || $storeError}
        <Alert type="error" message={error || $storeError} />
      {:else if $tasks.items && $tasks.items.length > 0}
        <div class="space-y-4">
          {#each $tasks.items as task (task.id)}
            <TaskListItem {task} />
          {/each}
        </div>
        
        <!-- Pagination -->
        {#if totalPages > 1}
          <div class="mt-6">
            <Pagination 
              currentPage={currentPage}
              totalPages={totalPages}
              on:pageChange={handlePageChange}
            />
          </div>
        {/if}
      {:else}
        <Card>
          <div class="flex flex-col items-center justify-center py-12">
            <svg class="w-16 h-16 text-gray-400 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01"></path>
            </svg>
            <h3 class="text-lg font-medium text-gray-900">No tasks found</h3>
            <p class="mt-1 text-sm text-gray-500">
              {#if statusFilter || priorityFilter || assigneeFilter || searchQuery}
                Try changing your filters or search query
              {:else}
                Get started by creating your first task
              {/if}
            </p>
            {#if statusFilter || priorityFilter || assigneeFilter || searchQuery}
              <button class="mt-4 btn btn-outline btn-sm" on:click={clearFilters}>
                Clear Filters
              </button>
            {:else}
              <a href="/tasks/create" class="mt-4 btn btn-primary btn-sm">
                Create Task
              </a>
            {/if}
          </div>
        </Card>
      {/if}
    </div>
  </div>
</div>

<style>
  /* Add any component-specific styles here */
</style>

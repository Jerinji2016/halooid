<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { DEFAULT_ORG_ID } from '$lib/config';
  import { taskStore } from '$lib/stores/taskStore';
  import { userStore } from '$lib/stores/userStore';
  import { STATUS_COLORS, PRIORITY_COLORS, DATE_FORMAT_OPTIONS } from '$lib/config';
  import TaskStatusBadge from '$lib/components/task/TaskStatusBadge.svelte';
  import TaskPriorityBadge from '$lib/components/task/TaskPriorityBadge.svelte';
  import TaskComments from '$lib/components/task/TaskComments.svelte';
  import TaskTimer from '$lib/components/task/TaskTimer.svelte';
  import TaskAssigneeSelector from '$lib/components/task/TaskAssigneeSelector.svelte';
  import Spinner from '$lib/components/common/Spinner.svelte';
  import Alert from '$lib/components/common/Alert.svelte';
  import ConfirmDialog from '$lib/components/common/ConfirmDialog.svelte';
  import { Button, Card } from '@halooid/ui-components';
  import { goto } from '$app/navigation';
  
  // Get task ID from URL
  const taskId = $page.params.id;
  
  // Local state
  let isLoading = true;
  let error = null;
  let task = null;
  let showDeleteConfirm = false;
  let isDeleting = false;
  let deleteError = null;
  
  // Initialize task store
  const { fetchTask, updateTask, deleteTask, isLoading: storeLoading, error: storeError } = taskStore;
  
  // Load task on mount
  onMount(async () => {
    try {
      isLoading = true;
      task = await fetchTask(taskId);
      isLoading = false;
    } catch (err) {
      error = err.message;
      isLoading = false;
    }
  });
  
  // Format date for display
  const formatDate = (dateString) => {
    if (!dateString) return 'Not set';
    return new Date(dateString).toLocaleDateString(undefined, DATE_FORMAT_OPTIONS);
  };
  
  // Format date with time for display
  const formatDateTime = (dateString) => {
    if (!dateString) return '';
    return new Date(dateString).toLocaleString();
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
  
  // Handle status change
  const handleStatusChange = async (newStatus) => {
    try {
      await updateTask(taskId, { status: newStatus });
      task = await fetchTask(taskId); // Refresh task data
    } catch (err) {
      error = err.message;
    }
  };
  
  // Handle assignee change
  const handleAssigneeChange = async (event) => {
    const assigneeId = event.detail.assigneeId;
    
    try {
      await updateTask(taskId, { assigned_to: assigneeId });
      task = await fetchTask(taskId); // Refresh task data
    } catch (err) {
      error = err.message;
    }
  };
  
  // Handle delete task
  const handleDeleteTask = async () => {
    try {
      isDeleting = true;
      await deleteTask(taskId);
      isDeleting = false;
      showDeleteConfirm = false;
      goto('/tasks');
    } catch (err) {
      deleteError = err.message;
      isDeleting = false;
    }
  };
</script>

<svelte:head>
  <title>{task ? task.title : 'Task Details'} | Taskake</title>
</svelte:head>

<div class="task-detail-page">
  {#if isLoading || $storeLoading}
    <div class="flex justify-center items-center h-64">
      <Spinner size="lg" />
    </div>
  {:else if error || $storeError}
    <Alert type="error" message={error || $storeError} />
  {:else if task}
    <div class="mb-6">
      <div class="flex flex-col md:flex-row justify-between items-start md:items-center">
        <div class="flex items-center space-x-2">
          <a href="/tasks" class="text-blue-600 hover:text-blue-800">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path>
            </svg>
          </a>
          <h1 class="text-2xl font-bold">{task.title}</h1>
        </div>
        
        <div class="flex space-x-2 mt-4 md:mt-0">
          <a href="/tasks/{taskId}/edit" class="btn btn-outline">
            <svg class="w-5 h-5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
            </svg>
            Edit
          </a>
          
          <button class="btn btn-outline btn-error" on:click={() => showDeleteConfirm = true}>
            <svg class="w-5 h-5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
            </svg>
            Delete
          </button>
        </div>
      </div>
    </div>
    
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <!-- Main Content -->
      <div class="md:col-span-2">
        <Card>
          <div class="p-6">
            <!-- Status and Priority -->
            <div class="flex flex-wrap gap-2 mb-6">
              <TaskStatusBadge status={task.status} />
              <TaskPriorityBadge priority={task.priority} />
              
              {#if task.tags && task.tags.length > 0}
                {#each task.tags as tag}
                  <span class="badge badge-outline">{tag}</span>
                {/each}
              {/if}
            </div>
            
            <!-- Description -->
            <div class="mb-6">
              <h2 class="text-lg font-semibold mb-2">Description</h2>
              <div class="prose max-w-none">
                {#if task.description}
                  <p>{task.description}</p>
                {:else}
                  <p class="text-gray-500 italic">No description provided</p>
                {/if}
              </div>
            </div>
            
            <!-- Time Tracking -->
            {#if task.status === 'in_progress'}
              <div class="mb-6">
                <TaskTimer taskId={taskId} />
              </div>
            {/if}
            
            <!-- Comments -->
            <div>
              <h2 class="text-lg font-semibold mb-2">Comments</h2>
              <TaskComments taskId={taskId} />
            </div>
          </div>
        </Card>
      </div>
      
      <!-- Sidebar -->
      <div class="md:col-span-1">
        <Card>
          <div class="p-6">
            <!-- Assignee -->
            <div class="mb-6">
              <h3 class="text-sm font-medium text-gray-500 mb-2">Assignee</h3>
              <TaskAssigneeSelector 
                taskId={taskId}
                currentAssigneeId={task.assigned_to}
                currentAssigneeName={task.assignee?.first_name ? `${task.assignee.first_name} ${task.assignee.last_name}` : 'Unassigned'}
                on:assigneeChange={handleAssigneeChange}
              />
            </div>
            
            <!-- Due Date -->
            <div class="mb-6">
              <h3 class="text-sm font-medium text-gray-500 mb-2">Due Date</h3>
              <p class="text-lg font-semibold">{formatDate(task.due_date)}</p>
              {#if task.due_date}
                <p class={getDaysRemainingClass(task.due_date)}>{getDaysRemaining(task.due_date)}</p>
              {/if}
            </div>
            
            <!-- Project -->
            {#if task.project_id}
              <div class="mb-6">
                <h3 class="text-sm font-medium text-gray-500 mb-2">Project</h3>
                <a href="/projects/{task.project_id}" class="text-blue-600 hover:text-blue-800 font-semibold">
                  {task.project_name || 'View Project'}
                </a>
              </div>
            {/if}
            
            <!-- Estimated Time -->
            {#if task.estimated_hours}
              <div class="mb-6">
                <h3 class="text-sm font-medium text-gray-500 mb-2">Estimated Time</h3>
                <p class="text-lg font-semibold">{task.estimated_hours} hours</p>
                
                {#if task.actual_hours}
                  <div class="mt-2">
                    <div class="w-full bg-gray-200 rounded-full h-2.5">
                      <div 
                        class="bg-blue-600 h-2.5 rounded-full" 
                        style="width: {Math.min((task.actual_hours / task.estimated_hours) * 100, 100)}%"
                      ></div>
                    </div>
                    <p class="text-sm mt-1">
                      {task.actual_hours} / {task.estimated_hours} hours ({Math.round((task.actual_hours / task.estimated_hours) * 100)}%)
                    </p>
                  </div>
                {/if}
              </div>
            {/if}
            
            <!-- Status Actions -->
            <div class="mb-6">
              <h3 class="text-sm font-medium text-gray-500 mb-2">Change Status</h3>
              <div class="flex flex-col space-y-2">
                {#if task.status !== 'todo'}
                  <button 
                    class="btn btn-sm btn-outline w-full justify-start"
                    style="color: {STATUS_COLORS.todo}; border-color: {STATUS_COLORS.todo};"
                    on:click={() => handleStatusChange('todo')}
                  >
                    To Do
                  </button>
                {/if}
                
                {#if task.status !== 'in_progress'}
                  <button 
                    class="btn btn-sm btn-outline w-full justify-start"
                    style="color: {STATUS_COLORS.in_progress}; border-color: {STATUS_COLORS.in_progress};"
                    on:click={() => handleStatusChange('in_progress')}
                  >
                    In Progress
                  </button>
                {/if}
                
                {#if task.status !== 'review'}
                  <button 
                    class="btn btn-sm btn-outline w-full justify-start"
                    style="color: {STATUS_COLORS.review}; border-color: {STATUS_COLORS.review};"
                    on:click={() => handleStatusChange('review')}
                  >
                    Review
                  </button>
                {/if}
                
                {#if task.status !== 'done'}
                  <button 
                    class="btn btn-sm btn-outline w-full justify-start"
                    style="color: {STATUS_COLORS.done}; border-color: {STATUS_COLORS.done};"
                    on:click={() => handleStatusChange('done')}
                  >
                    Done
                  </button>
                {/if}
                
                {#if task.status !== 'cancelled'}
                  <button 
                    class="btn btn-sm btn-outline w-full justify-start"
                    style="color: {STATUS_COLORS.cancelled}; border-color: {STATUS_COLORS.cancelled};"
                    on:click={() => handleStatusChange('cancelled')}
                  >
                    Cancelled
                  </button>
                {/if}
              </div>
            </div>
            
            <!-- Metadata -->
            <div>
              <h3 class="text-sm font-medium text-gray-500 mb-2">Details</h3>
              <div class="text-sm text-gray-500">
                <p>Created: {formatDateTime(task.created_at)}</p>
                {#if task.updated_at}
                  <p>Updated: {formatDateTime(task.updated_at)}</p>
                {/if}
                {#if task.creator}
                  <p>Created by: {task.creator.first_name} {task.creator.last_name}</p>
                {/if}
              </div>
            </div>
          </div>
        </Card>
      </div>
    </div>
    
    <!-- Delete Confirmation Dialog -->
    {#if showDeleteConfirm}
      <ConfirmDialog
        title="Delete Task"
        message="Are you sure you want to delete this task? This action cannot be undone."
        confirmText="Delete"
        cancelText="Cancel"
        confirmVariant="danger"
        isLoading={isDeleting}
        error={deleteError}
        on:confirm={handleDeleteTask}
        on:cancel={() => showDeleteConfirm = false}
      />
    {/if}
  {:else}
    <Alert type="warning" message="Task not found" />
  {/if}
</div>

<style>
  /* Add any component-specific styles here */
</style>

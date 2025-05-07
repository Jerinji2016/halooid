<script lang="ts">
  import { onMount, createEventDispatcher } from 'svelte';
  import { DEFAULT_ORG_ID } from '$lib/config';
  import type { UserResponse } from '$lib/types/user';
  import { Button } from '@halooid/ui-components';
  import Spinner from '../common/Spinner.svelte';
  import Alert from '../common/Alert.svelte';
  
  export let taskId: string;
  export let currentAssigneeId: string | null = null;
  export let currentAssigneeName: string = 'Unassigned';
  
  // Local state
  let users: UserResponse[] = [];
  let isLoading = true;
  let isSubmitting = false;
  let error = null;
  let showDropdown = false;
  
  const dispatch = createEventDispatcher();
  
  // API base URL
  const API_BASE_URL = import.meta.env.VITE_API_URL || 'https://api.halooid.com/v1';
  
  // Load users on mount
  onMount(async () => {
    await fetchUsers();
    
    // Close dropdown when clicking outside
    const handleClickOutside = (event) => {
      if (showDropdown && !event.target.closest('.assignee-selector')) {
        showDropdown = false;
      }
    };
    
    document.addEventListener('click', handleClickOutside);
    
    return () => {
      document.removeEventListener('click', handleClickOutside);
    };
  });
  
  // Fetch users
  const fetchUsers = async () => {
    isLoading = true;
    error = null;
    
    try {
      const response = await fetch(`${API_BASE_URL}/organizations/${DEFAULT_ORG_ID}/users`, {
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      });
      
      if (!response.ok) {
        throw new Error(`Failed to fetch users: ${response.statusText}`);
      }
      
      const data = await response.json();
      users = data.users || [];
      isLoading = false;
    } catch (err) {
      console.error('Error fetching users:', err);
      error = err.message;
      isLoading = false;
    }
  };
  
  // Assign task to user
  const assignTask = async (userId: string | null) => {
    isSubmitting = true;
    error = null;
    
    try {
      const response = await fetch(`${API_BASE_URL}/organizations/${DEFAULT_ORG_ID}/taskodex/tasks/${taskId}/assign`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify({
          user_id: userId
        })
      });
      
      if (!response.ok) {
        throw new Error(`Failed to assign task: ${response.statusText}`);
      }
      
      const data = await response.json();
      
      // Update current assignee
      currentAssigneeId = userId;
      
      if (userId) {
        const assignedUser = users.find(user => user.id === userId);
        if (assignedUser) {
          currentAssigneeName = `${assignedUser.first_name} ${assignedUser.last_name}`;
        }
      } else {
        currentAssigneeName = 'Unassigned';
      }
      
      // Notify parent component
      dispatch('assigneeChange', { assigneeId: userId });
      
      isSubmitting = false;
      showDropdown = false;
    } catch (err) {
      console.error('Error assigning task:', err);
      error = err.message;
      isSubmitting = false;
    }
  };
  
  // Toggle dropdown
  const toggleDropdown = () => {
    if (!isSubmitting) {
      showDropdown = !showDropdown;
    }
  };
</script>

<div class="assignee-selector relative">
  {#if isLoading}
    <div class="flex items-center">
      <Spinner size="sm" />
      <span class="ml-2">Loading users...</span>
    </div>
  {:else if error}
    <Alert type="error" message={error} />
  {:else}
    <div 
      class="flex items-center justify-between p-2 border rounded-lg cursor-pointer hover:bg-gray-50"
      on:click={toggleDropdown}
    >
      <div class="flex items-center">
        <div class="w-8 h-8 rounded-full bg-blue-100 flex items-center justify-center text-blue-600 font-semibold mr-2">
          {currentAssigneeName === 'Unassigned' ? 'UN' : currentAssigneeName.split(' ').map(n => n[0]).join('')}
        </div>
        <span>{currentAssigneeName}</span>
      </div>
      
      <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
      </svg>
    </div>
    
    {#if showDropdown}
      <div class="absolute z-10 mt-1 w-full bg-white border rounded-lg shadow-lg">
        <div class="p-2">
          <div class="mb-2">
            <input 
              type="text" 
              placeholder="Search users..." 
              class="w-full px-3 py-2 border rounded-lg focus:outline-none focus:border-blue-500"
            />
          </div>
          
          <div class="max-h-60 overflow-y-auto">
            <div 
              class="flex items-center p-2 hover:bg-gray-100 rounded-lg cursor-pointer {currentAssigneeId === null ? 'bg-blue-50' : ''}"
              on:click={() => assignTask(null)}
            >
              <div class="w-8 h-8 rounded-full bg-gray-200 flex items-center justify-center text-gray-600 font-semibold mr-2">
                UN
              </div>
              <span>Unassigned</span>
            </div>
            
            {#each users as user (user.id)}
              <div 
                class="flex items-center p-2 hover:bg-gray-100 rounded-lg cursor-pointer {currentAssigneeId === user.id ? 'bg-blue-50' : ''}"
                on:click={() => assignTask(user.id)}
              >
                <div class="w-8 h-8 rounded-full bg-blue-100 flex items-center justify-center text-blue-600 font-semibold mr-2">
                  {user.first_name[0]}{user.last_name[0]}
                </div>
                <span>{user.first_name} {user.last_name}</span>
              </div>
            {/each}
          </div>
        </div>
      </div>
    {/if}
    
    {#if isSubmitting}
      <div class="absolute inset-0 bg-white bg-opacity-75 flex items-center justify-center rounded-lg">
        <Spinner size="sm" />
      </div>
    {/if}
  {/if}
</div>

<style>
  /* Add any component-specific styles here */
</style>

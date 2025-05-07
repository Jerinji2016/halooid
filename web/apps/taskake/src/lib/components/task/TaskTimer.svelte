<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { DEFAULT_ORG_ID } from '$lib/config';
  import { Button } from '@halooid/ui-components';
  import Spinner from '../common/Spinner.svelte';
  import Alert from '../common/Alert.svelte';
  
  export let taskId: string;
  
  // Local state
  let isRunning = false;
  let startTime = null;
  let elapsedTime = 0;
  let timerInterval = null;
  let isLoading = true;
  let isSubmitting = false;
  let error = null;
  let timeEntries = [];
  let totalTime = 0;
  
  // API base URL
  const API_BASE_URL = import.meta.env.VITE_API_URL || 'https://api.halooid.com/v1';
  
  // Load time entries on mount
  onMount(async () => {
    await fetchTimeEntries();
    
    // Check if there's an active timer in localStorage
    const storedTimer = localStorage.getItem(`timer_${taskId}`);
    if (storedTimer) {
      const timerData = JSON.parse(storedTimer);
      startTime = new Date(timerData.startTime);
      isRunning = true;
      startTimer();
    }
  });
  
  // Clean up on destroy
  onDestroy(() => {
    if (timerInterval) {
      clearInterval(timerInterval);
    }
  });
  
  // Fetch time entries
  const fetchTimeEntries = async () => {
    isLoading = true;
    error = null;
    
    try {
      const response = await fetch(`${API_BASE_URL}/organizations/${DEFAULT_ORG_ID}/taskodex/tasks/${taskId}/time`, {
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      });
      
      if (!response.ok) {
        throw new Error(`Failed to fetch time entries: ${response.statusText}`);
      }
      
      const data = await response.json();
      timeEntries = data.entries || [];
      totalTime = data.total_hours || 0;
      isLoading = false;
    } catch (err) {
      console.error('Error fetching time entries:', err);
      error = err.message;
      isLoading = false;
    }
  };
  
  // Start timer
  const startTimer = () => {
    if (!startTime) {
      startTime = new Date();
      
      // Store timer start in localStorage
      localStorage.setItem(`timer_${taskId}`, JSON.stringify({
        startTime: startTime.toISOString()
      }));
    }
    
    isRunning = true;
    
    timerInterval = setInterval(() => {
      const now = new Date();
      elapsedTime = Math.floor((now.getTime() - startTime.getTime()) / 1000);
    }, 1000);
  };
  
  // Stop timer
  const stopTimer = async () => {
    if (timerInterval) {
      clearInterval(timerInterval);
      timerInterval = null;
    }
    
    isRunning = false;
    
    // Calculate hours
    const now = new Date();
    const durationMs = now.getTime() - startTime.getTime();
    const durationHours = durationMs / (1000 * 60 * 60);
    
    // Remove timer from localStorage
    localStorage.removeItem(`timer_${taskId}`);
    
    // Submit time entry
    isSubmitting = true;
    error = null;
    
    try {
      const response = await fetch(`${API_BASE_URL}/organizations/${DEFAULT_ORG_ID}/taskodex/tasks/${taskId}/time`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify({
          hours: durationHours,
          start_time: startTime.toISOString(),
          end_time: now.toISOString()
        })
      });
      
      if (!response.ok) {
        throw new Error(`Failed to submit time entry: ${response.statusText}`);
      }
      
      // Reset timer
      startTime = null;
      elapsedTime = 0;
      
      // Refresh time entries
      await fetchTimeEntries();
      
      isSubmitting = false;
    } catch (err) {
      console.error('Error submitting time entry:', err);
      error = err.message;
      isSubmitting = false;
    }
  };
  
  // Format time as HH:MM:SS
  const formatTime = (seconds) => {
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    const secs = seconds % 60;
    
    return [
      hours.toString().padStart(2, '0'),
      minutes.toString().padStart(2, '0'),
      secs.toString().padStart(2, '0')
    ].join(':');
  };
  
  // Format date for display
  const formatDate = (dateString) => {
    if (!dateString) return '';
    return new Date(dateString).toLocaleDateString();
  };
</script>

<div class="task-timer">
  <h3 class="text-lg font-semibold mb-2">Time Tracking</h3>
  
  {#if isLoading}
    <div class="flex justify-center items-center h-24">
      <Spinner />
    </div>
  {:else if error}
    <Alert type="error" message={error} />
  {:else}
    <!-- Timer Display -->
    <div class="bg-gray-50 p-4 rounded-lg mb-4">
      <div class="flex justify-between items-center">
        <div>
          <div class="text-2xl font-mono font-semibold">
            {#if isRunning}
              {formatTime(elapsedTime)}
            {:else}
              00:00:00
            {/if}
          </div>
          
          <div class="text-sm text-gray-500 mt-1">
            {#if isRunning}
              Started at {startTime.toLocaleTimeString()}
            {:else}
              Timer stopped
            {/if}
          </div>
        </div>
        
        <div>
          {#if isRunning}
            <Button
              variant="danger"
              on:click={stopTimer}
              disabled={isSubmitting}
            >
              {#if isSubmitting}
                <Spinner size="sm" />
                <span class="ml-2">Stopping...</span>
              {:else}
                Stop Timer
              {/if}
            </Button>
          {:else}
            <Button
              variant="primary"
              on:click={startTimer}
            >
              Start Timer
            </Button>
          {/if}
        </div>
      </div>
    </div>
    
    <!-- Total Time -->
    <div class="mb-4">
      <div class="flex justify-between items-center">
        <div class="text-sm font-medium text-gray-700">Total Time:</div>
        <div class="text-lg font-semibold">{totalTime.toFixed(2)} hours</div>
      </div>
    </div>
    
    <!-- Time Entries -->
    {#if timeEntries.length > 0}
      <div class="border rounded-lg overflow-hidden">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Date
              </th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Duration
              </th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                User
              </th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            {#each timeEntries as entry (entry.id)}
              <tr>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {formatDate(entry.start_time)}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {entry.hours.toFixed(2)} hours
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {entry.user?.first_name || 'Unknown'} {entry.user?.last_name || 'User'}
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    {:else}
      <p class="text-gray-500 italic">No time entries yet</p>
    {/if}
  {/if}
</div>

<style>
  /* Add any component-specific styles here */
</style>

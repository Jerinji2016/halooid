<script lang="ts">
  import { onMount } from 'svelte';
  import { DEFAULT_ORG_ID } from '$lib/config';
  import type { TaskComment } from '$lib/types/task';
  import { Button, Input } from '@halooid/ui-components';
  import Spinner from '../common/Spinner.svelte';
  import Alert from '../common/Alert.svelte';
  
  export let taskId: string;
  
  // Local state
  let comments: TaskComment[] = [];
  let newComment: string = '';
  let isLoading = true;
  let isSubmitting = false;
  let error = null;
  
  // API base URL
  const API_BASE_URL = import.meta.env.VITE_API_URL || 'https://api.halooid.com/v1';
  
  // Load comments on mount
  onMount(async () => {
    await fetchComments();
  });
  
  // Fetch comments
  const fetchComments = async () => {
    isLoading = true;
    error = null;
    
    try {
      const response = await fetch(`${API_BASE_URL}/organizations/${DEFAULT_ORG_ID}/taskodex/tasks/${taskId}/comments`, {
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      });
      
      if (!response.ok) {
        throw new Error(`Failed to fetch comments: ${response.statusText}`);
      }
      
      const data = await response.json();
      comments = data.comments || [];
      isLoading = false;
    } catch (err) {
      console.error('Error fetching comments:', err);
      error = err.message;
      isLoading = false;
    }
  };
  
  // Add comment
  const addComment = async () => {
    if (!newComment.trim()) {
      return;
    }
    
    isSubmitting = true;
    error = null;
    
    try {
      const response = await fetch(`${API_BASE_URL}/organizations/${DEFAULT_ORG_ID}/taskodex/tasks/${taskId}/comments`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify({
          content: newComment
        })
      });
      
      if (!response.ok) {
        throw new Error(`Failed to add comment: ${response.statusText}`);
      }
      
      const data = await response.json();
      comments = [...comments, data];
      newComment = '';
      isSubmitting = false;
    } catch (err) {
      console.error('Error adding comment:', err);
      error = err.message;
      isSubmitting = false;
    }
  };
  
  // Format date for display
  const formatDate = (dateString) => {
    if (!dateString) return '';
    
    const date = new Date(dateString);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffSec = Math.round(diffMs / 1000);
    const diffMin = Math.round(diffSec / 60);
    const diffHour = Math.round(diffMin / 60);
    const diffDay = Math.round(diffHour / 24);
    
    if (diffSec < 60) {
      return `${diffSec} second${diffSec !== 1 ? 's' : ''} ago`;
    } else if (diffMin < 60) {
      return `${diffMin} minute${diffMin !== 1 ? 's' : ''} ago`;
    } else if (diffHour < 24) {
      return `${diffHour} hour${diffHour !== 1 ? 's' : ''} ago`;
    } else if (diffDay < 30) {
      return `${diffDay} day${diffDay !== 1 ? 's' : ''} ago`;
    } else {
      return date.toLocaleDateString();
    }
  };
</script>

<div class="task-comments">
  {#if isLoading}
    <div class="flex justify-center items-center h-24">
      <Spinner />
    </div>
  {:else if error}
    <Alert type="error" message={error} />
  {:else}
    <!-- Comment List -->
    <div class="space-y-4 mb-6">
      {#if comments.length === 0}
        <p class="text-gray-500 italic">No comments yet</p>
      {:else}
        {#each comments as comment (comment.id)}
          <div class="bg-gray-50 p-4 rounded-lg">
            <div class="flex justify-between items-start">
              <div class="flex items-start">
                <div class="flex-shrink-0">
                  <div class="w-10 h-10 rounded-full bg-blue-100 flex items-center justify-center text-blue-600 font-semibold">
                    {comment.user?.first_name?.charAt(0) || '?'}{comment.user?.last_name?.charAt(0) || ''}
                  </div>
                </div>
                
                <div class="ml-3">
                  <p class="text-sm font-medium text-gray-900">
                    {comment.user?.first_name || 'Unknown'} {comment.user?.last_name || 'User'}
                  </p>
                  <p class="text-sm text-gray-500">
                    {formatDate(comment.created_at)}
                  </p>
                </div>
              </div>
            </div>
            
            <div class="mt-2 text-sm text-gray-700 whitespace-pre-line">
              {comment.content}
            </div>
          </div>
        {/each}
      {/if}
    </div>
    
    <!-- Add Comment Form -->
    <div>
      <div class="mb-2">
        <textarea
          class="w-full px-3 py-2 text-gray-700 border rounded-lg focus:outline-none focus:border-blue-500"
          rows="3"
          placeholder="Add a comment..."
          bind:value={newComment}
        ></textarea>
      </div>
      
      <div class="flex justify-end">
        <Button
          variant="primary"
          on:click={addComment}
          disabled={isSubmitting || !newComment.trim()}
        >
          {#if isSubmitting}
            <Spinner size="sm" />
            <span class="ml-2">Posting...</span>
          {:else}
            Post Comment
          {/if}
        </Button>
      </div>
    </div>
  {/if}
</div>

<style>
  /* Add any component-specific styles here */
</style>

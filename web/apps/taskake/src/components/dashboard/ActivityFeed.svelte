<script lang="ts">
  import { DATE_TIME_FORMAT_OPTIONS } from '$lib/config';
  import type { ProjectActivity } from '$lib/types/project';
  
  export let activities: ProjectActivity[];
  
  // Format date for display
  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString(undefined, DATE_TIME_FORMAT_OPTIONS);
  };
  
  // Format time ago
  const formatTimeAgo = (dateString: string) => {
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
      return formatDate(dateString);
    }
  };
  
  // Get icon for activity
  const getActivityIcon = (action: string, entityType: string) => {
    if (action === 'created') {
      if (entityType === 'task') return 'plus-circle';
      if (entityType === 'comment') return 'message-circle';
      return 'plus';
    } else if (action === 'updated') {
      if (entityType === 'task') return 'edit';
      return 'edit-2';
    } else if (action === 'deleted') {
      return 'trash-2';
    } else if (action === 'completed') {
      return 'check-circle';
    } else if (action === 'assigned') {
      return 'user-plus';
    } else if (action === 'unassigned') {
      return 'user-minus';
    } else if (action === 'commented') {
      return 'message-square';
    } else {
      return 'activity';
    }
  };
  
  // Get activity description
  const getActivityDescription = (activity: ProjectActivity) => {
    const { action, entity_type, entity_name } = activity;
    
    if (action === 'created') {
      return `created ${entity_type} "${entity_name}"`;
    } else if (action === 'updated') {
      return `updated ${entity_type} "${entity_name}"`;
    } else if (action === 'deleted') {
      return `deleted ${entity_type} "${entity_name}"`;
    } else if (action === 'completed') {
      return `completed ${entity_type} "${entity_name}"`;
    } else if (action === 'assigned') {
      return `assigned ${entity_type} "${entity_name}"`;
    } else if (action === 'unassigned') {
      return `unassigned ${entity_type} "${entity_name}"`;
    } else if (action === 'commented') {
      return `commented on ${entity_type} "${entity_name}"`;
    } else {
      return `${action} ${entity_type} "${entity_name}"`;
    }
  };
</script>

<div class="activity-feed bg-white rounded-lg shadow-md p-6">
  <h2 class="text-lg font-semibold mb-4">Recent Activity</h2>
  
  {#if activities.length === 0}
    <div class="flex items-center justify-center h-32 text-gray-500">
      No recent activity
    </div>
  {:else}
    <div class="activities">
      {#each activities as activity}
        <div class="activity flex items-start py-3 border-b">
          <div class="activity-icon mr-3 mt-1">
            <svg class="w-5 h-5 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
            </svg>
          </div>
          
          <div class="activity-content flex-1">
            <div class="flex justify-between">
              <div class="font-medium">{activity.user_name}</div>
              <div class="text-sm text-gray-500" title={formatDate(activity.created_at)}>
                {formatTimeAgo(activity.created_at)}
              </div>
            </div>
            
            <div class="text-sm text-gray-600 mt-1">
              {getActivityDescription(activity)}
            </div>
          </div>
        </div>
      {/each}
    </div>
    
    {#if activities.length >= 10}
      <div class="mt-4 text-center">
        <a href="/projects/{activities[0].project_id}/activity" class="text-blue-500 hover:text-blue-700 text-sm">
          View all activity
        </a>
      </div>
    {/if}
  {/if}
</div>

<style>
  /* Add any component-specific styles here */
</style>

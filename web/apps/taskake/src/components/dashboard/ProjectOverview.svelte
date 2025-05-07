<script lang="ts">
  import { currentProject, completionPercentage } from '$lib/stores/projectStore';
  import { PROJECT_STATUS_COLORS, DATE_FORMAT_OPTIONS } from '$lib/config';
  import type { ProjectResponse } from '$lib/types/project';
  
  export let project: ProjectResponse;
  
  // Format dates
  const formatDate = (dateString: string | undefined) => {
    if (!dateString) return 'Not set';
    return new Date(dateString).toLocaleDateString(undefined, DATE_FORMAT_OPTIONS);
  };
  
  // Get status color
  const getStatusColor = (status: string) => {
    return PROJECT_STATUS_COLORS[status as keyof typeof PROJECT_STATUS_COLORS] || '#9E9E9E';
  };
  
  // Calculate days remaining or overdue
  const getDaysRemaining = () => {
    if (!project.end_date) return 'No deadline';
    
    const endDate = new Date(project.end_date);
    const today = new Date();
    
    // Reset time part for accurate day calculation
    endDate.setHours(0, 0, 0, 0);
    today.setHours(0, 0, 0, 0);
    
    const diffTime = endDate.getTime() - today.getTime();
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
</script>

<div class="project-overview bg-white rounded-lg shadow-md p-6 mb-6">
  <div class="flex flex-col md:flex-row justify-between items-start md:items-center mb-4">
    <div>
      <h1 class="text-2xl font-bold">{project.name}</h1>
      <p class="text-gray-600 mt-1">{project.description || 'No description'}</p>
    </div>
    
    <div class="mt-4 md:mt-0">
      <span 
        class="px-3 py-1 rounded-full text-sm font-medium"
        style="background-color: {getStatusColor(project.status)}25; color: {getStatusColor(project.status)};"
      >
        {project.status.charAt(0).toUpperCase() + project.status.slice(1).replace('_', ' ')}
      </span>
    </div>
  </div>
  
  <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mt-6">
    <div class="bg-gray-50 p-4 rounded-lg">
      <h3 class="text-sm font-medium text-gray-500">Start Date</h3>
      <p class="mt-1 text-lg font-semibold">{formatDate(project.start_date)}</p>
    </div>
    
    <div class="bg-gray-50 p-4 rounded-lg">
      <h3 class="text-sm font-medium text-gray-500">End Date</h3>
      <p class="mt-1 text-lg font-semibold">{formatDate(project.end_date)}</p>
    </div>
    
    <div class="bg-gray-50 p-4 rounded-lg">
      <h3 class="text-sm font-medium text-gray-500">Time Remaining</h3>
      <p class="mt-1 text-lg font-semibold">{getDaysRemaining()}</p>
    </div>
    
    <div class="bg-gray-50 p-4 rounded-lg">
      <h3 class="text-sm font-medium text-gray-500">Completion</h3>
      <div class="mt-1">
        <div class="w-full bg-gray-200 rounded-full h-2.5">
          <div 
            class="bg-blue-600 h-2.5 rounded-full" 
            style="width: {$completionPercentage}%"
          ></div>
        </div>
        <p class="mt-1 text-lg font-semibold">{$completionPercentage}%</p>
      </div>
    </div>
  </div>
  
  <div class="mt-6 flex flex-wrap gap-2">
    <a 
      href="/projects/{project.id}/tasks" 
      class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
    >
      View Tasks
    </a>
    
    <a 
      href="/projects/{project.id}/edit" 
      class="inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md shadow-sm text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
    >
      Edit Project
    </a>
    
    <button 
      class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500"
    >
      Add Task
    </button>
  </div>
</div>

<style>
  /* Add any component-specific styles here */
</style>

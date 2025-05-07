<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { DEFAULT_ORG_ID } from '$lib/config';
  import { projectActions, currentProject, projectTasks, projectStats, projectActivity, isLoading, error } from '$lib/stores/projectStore';
  
  // Import dashboard components
  import ProjectOverview from '../../../../components/dashboard/ProjectOverview.svelte';
  import TaskStatusChart from '../../../../components/dashboard/TaskStatusChart.svelte';
  import TaskAssignmentChart from '../../../../components/dashboard/TaskAssignmentChart.svelte';
  import ProjectTimeline from '../../../../components/dashboard/ProjectTimeline.svelte';
  import ActivityFeed from '../../../../components/dashboard/ActivityFeed.svelte';
  import MetricsPanel from '../../../../components/dashboard/MetricsPanel.svelte';
  import UpcomingDeadlines from '../../../../components/dashboard/UpcomingDeadlines.svelte';
  
  // Get project ID from URL
  const projectId = $page.params.id;
  
  // Initialize project actions
  const actions = projectActions(DEFAULT_ORG_ID);
  
  // Load project data
  onMount(async () => {
    await actions.loadAllProjectData(projectId);
  });
  
  // Refresh data
  const refreshData = async () => {
    await actions.loadAllProjectData(projectId);
  };
</script>

<svelte:head>
  <title>{$currentProject?.name || 'Project'} Dashboard | Taskake</title>
</svelte:head>

<div class="project-dashboard p-4 md:p-6">
  {#if $isLoading && !$currentProject}
    <div class="flex justify-center items-center h-64">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
    </div>
  {:else if $error}
    <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative mb-6" role="alert">
      <strong class="font-bold">Error!</strong>
      <span class="block sm:inline"> {$error}</span>
    </div>
  {:else if $currentProject}
    <!-- Project Overview -->
    <ProjectOverview project={$currentProject} />
    
    <!-- Refresh Button -->
    <div class="flex justify-end mb-6">
      <button 
        on:click={refreshData}
        class="inline-flex items-center px-3 py-2 border border-gray-300 text-sm font-medium rounded-md shadow-sm text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
        disabled={$isLoading}
      >
        {#if $isLoading}
          <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-gray-700" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          Refreshing...
        {:else}
          <svg class="-ml-1 mr-2 h-4 w-4 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
          </svg>
          Refresh
        {/if}
      </button>
    </div>
    
    <!-- Metrics and Charts -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
      {#if $projectStats}
        <MetricsPanel stats={$projectStats} />
      {/if}
      
      {#if $projectStats?.tasksByStatus?.length > 0}
        <TaskStatusChart tasksByStatus={$projectStats.tasksByStatus} />
      {/if}
    </div>
    
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
      {#if $projectStats?.tasksByAssignee?.length > 0}
        <TaskAssignmentChart tasksByAssignee={$projectStats.tasksByAssignee} />
      {/if}
      
      {#if $projectTasks?.length > 0}
        <UpcomingDeadlines tasks={$projectTasks} />
      {/if}
    </div>
    
    <!-- Timeline -->
    {#if $currentProject && $projectTasks?.length > 0}
      <div class="mb-6">
        <ProjectTimeline project={$currentProject} tasks={$projectTasks} />
      </div>
    {/if}
    
    <!-- Activity Feed -->
    {#if $projectActivity?.length > 0}
      <div>
        <ActivityFeed activities={$projectActivity} />
      </div>
    {/if}
  {:else}
    <div class="bg-yellow-100 border border-yellow-400 text-yellow-700 px-4 py-3 rounded relative mb-6" role="alert">
      <strong class="font-bold">Project not found!</strong>
      <span class="block sm:inline"> The requested project could not be found.</span>
    </div>
  {/if}
</div>

<style>
  /* Add any page-specific styles here */
</style>

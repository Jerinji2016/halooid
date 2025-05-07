<script lang="ts">
  import type { ProjectStats } from '$lib/types/project';
  
  export let stats: ProjectStats;
  
  // Calculate completion rate
  const completionRate = stats.totalTasks > 0 
    ? Math.round((stats.completedTasks / stats.totalTasks) * 100) 
    : 0;
  
  // Calculate overdue rate
  const overdueRate = stats.totalTasks > 0 
    ? Math.round((stats.overdueTasks / stats.totalTasks) * 100) 
    : 0;
</script>

<div class="metrics-panel bg-white rounded-lg shadow-md p-6">
  <h2 class="text-lg font-semibold mb-4">Project Metrics</h2>
  
  <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
    <div class="metric-card bg-blue-50 p-4 rounded-lg">
      <h3 class="text-sm font-medium text-blue-700">Total Tasks</h3>
      <p class="mt-2 text-2xl font-bold text-blue-900">{stats.totalTasks}</p>
    </div>
    
    <div class="metric-card bg-green-50 p-4 rounded-lg">
      <h3 class="text-sm font-medium text-green-700">Completed Tasks</h3>
      <p class="mt-2 text-2xl font-bold text-green-900">{stats.completedTasks}</p>
      <p class="text-xs text-green-700 mt-1">{completionRate}% completion rate</p>
    </div>
    
    <div class="metric-card bg-red-50 p-4 rounded-lg">
      <h3 class="text-sm font-medium text-red-700">Overdue Tasks</h3>
      <p class="mt-2 text-2xl font-bold text-red-900">{stats.overdueTasks}</p>
      <p class="text-xs text-red-700 mt-1">{overdueRate}% overdue rate</p>
    </div>
    
    <div class="metric-card bg-purple-50 p-4 rounded-lg">
      <h3 class="text-sm font-medium text-purple-700">Completion</h3>
      <div class="mt-2">
        <div class="w-full bg-gray-200 rounded-full h-2.5">
          <div 
            class="bg-purple-600 h-2.5 rounded-full" 
            style="width: {stats.completionPercentage}%"
          ></div>
        </div>
        <p class="mt-1 text-2xl font-bold text-purple-900">{stats.completionPercentage}%</p>
      </div>
    </div>
  </div>
  
  <div class="mt-6">
    <h3 class="text-sm font-medium text-gray-700 mb-2">Task Distribution</h3>
    
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div>
        <h4 class="text-xs font-medium text-gray-500 mb-1">By Status</h4>
        <div class="space-y-2">
          {#each stats.tasksByStatus as { status, count }}
            <div>
              <div class="flex justify-between text-sm">
                <span>{status.charAt(0).toUpperCase() + status.slice(1).replace('_', ' ')}</span>
                <span>{count}</span>
              </div>
              <div class="w-full bg-gray-200 rounded-full h-1.5 mt-1">
                <div 
                  class="h-1.5 rounded-full" 
                  style="width: {(count / stats.totalTasks) * 100}%; background-color: var(--status-{status}-color, #9E9E9E);"
                ></div>
              </div>
            </div>
          {/each}
        </div>
      </div>
      
      <div>
        <h4 class="text-xs font-medium text-gray-500 mb-1">By Priority</h4>
        <div class="space-y-2">
          {#each stats.tasksByPriority as { priority, count }}
            <div>
              <div class="flex justify-between text-sm">
                <span>{priority.charAt(0).toUpperCase() + priority.slice(1)}</span>
                <span>{count}</span>
              </div>
              <div class="w-full bg-gray-200 rounded-full h-1.5 mt-1">
                <div 
                  class="h-1.5 rounded-full" 
                  style="width: {(count / stats.totalTasks) * 100}%; background-color: var(--priority-{priority}-color, #9E9E9E);"
                ></div>
              </div>
            </div>
          {/each}
        </div>
      </div>
    </div>
  </div>
</div>

<style>
  /* Status colors */
  :root {
    --status-todo-color: #9E9E9E;
    --status-in_progress-color: #2196F3;
    --status-review-color: #FF9800;
    --status-done-color: #4CAF50;
    --status-cancelled-color: #F44336;
    
    --priority-low-color: #4CAF50;
    --priority-medium-color: #FFC107;
    --priority-high-color: #FF9800;
    --priority-critical-color: #F44336;
  }
</style>

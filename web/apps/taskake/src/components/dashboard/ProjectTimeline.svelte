<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { STATUS_COLORS, PRIORITY_COLORS } from '$lib/config';
  import type { TaskResponse } from '$lib/types/task';
  import type { ProjectResponse } from '$lib/types/project';
  
  export let project: ProjectResponse;
  export let tasks: TaskResponse[];
  
  let ganttElement: HTMLDivElement;
  
  // Format date for display
  const formatDate = (dateString: string | undefined) => {
    if (!dateString) return '';
    return new Date(dateString).toLocaleDateString();
  };
  
  // Get task color based on status
  const getTaskColor = (task: TaskResponse) => {
    return STATUS_COLORS[task.status as keyof typeof STATUS_COLORS] || '#9E9E9E';
  };
  
  // Get task border color based on priority
  const getTaskBorderColor = (task: TaskResponse) => {
    return PRIORITY_COLORS[task.priority as keyof typeof PRIORITY_COLORS] || '#9E9E9E';
  };
  
  // Sort tasks by due date
  const sortedTasks = () => {
    return [...tasks]
      .filter(task => task.due_date)
      .sort((a, b) => {
        if (!a.due_date || !b.due_date) return 0;
        return new Date(a.due_date).getTime() - new Date(b.due_date).getTime();
      });
  };
  
  // Calculate timeline dates
  const getTimelineDates = () => {
    const sorted = sortedTasks();
    if (sorted.length === 0) {
      // Default to project dates or current month
      const start = project.start_date ? new Date(project.start_date) : new Date();
      const end = project.end_date ? new Date(project.end_date) : new Date();
      end.setMonth(end.getMonth() + 1);
      return { start, end };
    }
    
    const start = new Date(sorted[0].due_date!);
    const end = new Date(sorted[sorted.length - 1].due_date!);
    
    // Add padding
    start.setDate(start.getDate() - 7);
    end.setDate(end.getDate() + 7);
    
    return { start, end };
  };
  
  // Generate timeline dates
  const generateTimelineDates = () => {
    const { start, end } = getTimelineDates();
    const dates = [];
    const current = new Date(start);
    
    while (current <= end) {
      dates.push(new Date(current));
      current.setDate(current.getDate() + 1);
    }
    
    return dates;
  };
  
  // Calculate task position and width
  const getTaskStyle = (task: TaskResponse) => {
    if (!task.due_date) return '';
    
    const { start, end } = getTimelineDates();
    const taskDate = new Date(task.due_date);
    
    const totalDays = (end.getTime() - start.getTime()) / (1000 * 60 * 60 * 24);
    const taskDays = (taskDate.getTime() - start.getTime()) / (1000 * 60 * 60 * 24);
    
    const leftPosition = (taskDays / totalDays) * 100;
    
    return `left: ${leftPosition}%;`;
  };
</script>

<div class="project-timeline bg-white rounded-lg shadow-md p-6">
  <h2 class="text-lg font-semibold mb-4">Project Timeline</h2>
  
  {#if tasks.filter(t => t.due_date).length === 0}
    <div class="flex items-center justify-center h-32 text-gray-500">
      No tasks with due dates
    </div>
  {:else}
    <div class="timeline-container overflow-x-auto" bind:this={ganttElement}>
      <div class="timeline-header flex border-b pb-2">
        <div class="w-1/4 font-medium">Task</div>
        <div class="w-1/4 font-medium">Assignee</div>
        <div class="w-1/4 font-medium">Status</div>
        <div class="w-1/4 font-medium">Due Date</div>
      </div>
      
      <div class="timeline-tasks">
        {#each sortedTasks() as task}
          <div class="timeline-task flex py-2 border-b">
            <div class="w-1/4 truncate" title={task.title}>{task.title}</div>
            <div class="w-1/4">{task.assignee?.first_name || 'Unassigned'}</div>
            <div class="w-1/4">
              <span 
                class="px-2 py-1 rounded-full text-xs font-medium"
                style="background-color: {getTaskColor(task)}25; color: {getTaskColor(task)};"
              >
                {task.status.charAt(0).toUpperCase() + task.status.slice(1).replace('_', ' ')}
              </span>
            </div>
            <div class="w-1/4">{formatDate(task.due_date)}</div>
          </div>
        {/each}
      </div>
      
      <div class="timeline-gantt mt-8 relative h-20 border-t border-b">
        {#each generateTimelineDates() as date, i}
          {#if date.getDate() === 1 || i === 0}
            <div 
              class="timeline-month absolute top-0 bottom-0 border-l"
              style="left: {(i / generateTimelineDates().length) * 100}%;"
            >
              <div class="text-xs text-gray-500 ml-1">
                {date.toLocaleDateString(undefined, { month: 'short', day: 'numeric' })}
              </div>
            </div>
          {/if}
        {/each}
        
        {#each sortedTasks() as task}
          <div 
            class="timeline-task-marker absolute w-4 h-4 rounded-full -mt-2 -ml-2"
            style="{getTaskStyle(task)} background-color: {getTaskColor(task)}; border: 2px solid {getTaskBorderColor(task)}; top: 50%;"
            title={`${task.title} - Due: ${formatDate(task.due_date)}`}
          ></div>
        {/each}
        
        <!-- Today marker -->
        <div 
          class="timeline-today absolute top-0 bottom-0 border-l-2 border-red-500"
          style="left: {((new Date().getTime() - getTimelineDates().start.getTime()) / (getTimelineDates().end.getTime() - getTimelineDates().start.getTime())) * 100}%;"
        >
          <div class="text-xs text-red-500 ml-1">Today</div>
        </div>
      </div>
    </div>
  {/if}
</div>

<style>
  .timeline-container {
    min-height: 200px;
  }
  
  .timeline-gantt {
    min-width: 500px;
  }
</style>

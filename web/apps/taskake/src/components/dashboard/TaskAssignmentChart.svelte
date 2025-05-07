<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import Chart from 'chart.js/auto';
  
  export let tasksByAssignee: { assignee_id: string, assignee_name: string, count: number }[];
  
  let chartElement: HTMLCanvasElement;
  let chart: Chart;
  
  // Generate random colors for assignees
  const generateColors = (count: number) => {
    const colors = [];
    for (let i = 0; i < count; i++) {
      const hue = (i * 137) % 360; // Use golden angle approximation for better distribution
      colors.push(`hsl(${hue}, 70%, 60%)`);
    }
    return colors;
  };
  
  // Create chart
  const createChart = () => {
    if (!chartElement) return;
    
    const labels = tasksByAssignee.map(item => item.assignee_name || 'Unassigned');
    const data = tasksByAssignee.map(item => item.count);
    const backgroundColors = generateColors(tasksByAssignee.length);
    
    chart = new Chart(chartElement, {
      type: 'bar',
      data: {
        labels,
        datasets: [{
          label: 'Tasks',
          data,
          backgroundColor: backgroundColors,
          borderWidth: 1
        }]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        indexAxis: 'y',
        plugins: {
          legend: {
            display: false
          },
          tooltip: {
            callbacks: {
              label: (context) => {
                const value = context.raw as number;
                return `Tasks: ${value}`;
              }
            }
          }
        },
        scales: {
          x: {
            beginAtZero: true,
            ticks: {
              precision: 0
            }
          }
        }
      }
    });
  };
  
  // Update chart
  const updateChart = () => {
    if (!chart) return;
    
    const labels = tasksByAssignee.map(item => item.assignee_name || 'Unassigned');
    const data = tasksByAssignee.map(item => item.count);
    const backgroundColors = generateColors(tasksByAssignee.length);
    
    chart.data.labels = labels;
    chart.data.datasets[0].data = data;
    chart.data.datasets[0].backgroundColor = backgroundColors;
    chart.update();
  };
  
  onMount(() => {
    createChart();
  });
  
  onDestroy(() => {
    if (chart) {
      chart.destroy();
    }
  });
  
  // Watch for changes in tasksByAssignee
  $: if (chart && tasksByAssignee) {
    updateChart();
  }
</script>

<div class="task-assignment-chart bg-white rounded-lg shadow-md p-6">
  <h2 class="text-lg font-semibold mb-4">Tasks by Assignee</h2>
  
  {#if tasksByAssignee.length === 0}
    <div class="flex items-center justify-center h-64 text-gray-500">
      No tasks assigned
    </div>
  {:else}
    <div class="chart-container" style="position: relative; height: 300px;">
      <canvas bind:this={chartElement}></canvas>
    </div>
  {/if}
</div>

<style>
  /* Add any component-specific styles here */
</style>

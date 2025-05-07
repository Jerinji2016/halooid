<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import Chart from 'chart.js/auto';
  import { STATUS_COLORS } from '$lib/config';
  
  export let tasksByStatus: { status: string, count: number }[];
  
  let chartElement: HTMLCanvasElement;
  let chart: Chart;
  
  // Format status labels
  const formatStatusLabel = (status: string) => {
    return status.charAt(0).toUpperCase() + status.slice(1).replace('_', ' ');
  };
  
  // Create chart
  const createChart = () => {
    if (!chartElement) return;
    
    const labels = tasksByStatus.map(item => formatStatusLabel(item.status));
    const data = tasksByStatus.map(item => item.count);
    const backgroundColors = tasksByStatus.map(item => 
      STATUS_COLORS[item.status as keyof typeof STATUS_COLORS] || '#9E9E9E'
    );
    
    chart = new Chart(chartElement, {
      type: 'doughnut',
      data: {
        labels,
        datasets: [{
          data,
          backgroundColor: backgroundColors,
          borderWidth: 1
        }]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          legend: {
            position: 'right',
            labels: {
              usePointStyle: true,
              padding: 20
            }
          },
          tooltip: {
            callbacks: {
              label: (context) => {
                const label = context.label || '';
                const value = context.raw as number;
                const total = data.reduce((a, b) => a + b, 0);
                const percentage = Math.round((value / total) * 100);
                return `${label}: ${value} (${percentage}%)`;
              }
            }
          }
        }
      }
    });
  };
  
  // Update chart
  const updateChart = () => {
    if (!chart) return;
    
    const labels = tasksByStatus.map(item => formatStatusLabel(item.status));
    const data = tasksByStatus.map(item => item.count);
    const backgroundColors = tasksByStatus.map(item => 
      STATUS_COLORS[item.status as keyof typeof STATUS_COLORS] || '#9E9E9E'
    );
    
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
  
  // Watch for changes in tasksByStatus
  $: if (chart && tasksByStatus) {
    updateChart();
  }
</script>

<div class="task-status-chart bg-white rounded-lg shadow-md p-6">
  <h2 class="text-lg font-semibold mb-4">Tasks by Status</h2>
  
  {#if tasksByStatus.length === 0}
    <div class="flex items-center justify-center h-64 text-gray-500">
      No tasks available
    </div>
  {:else}
    <div class="chart-container" style="position: relative; height: 300px;">
      <canvas bind:this={chartElement}></canvas>
    </div>
    
    <div class="mt-4 grid grid-cols-2 md:grid-cols-5 gap-2">
      {#each tasksByStatus as { status, count }}
        <div class="flex items-center">
          <div 
            class="w-3 h-3 rounded-full mr-2" 
            style="background-color: {STATUS_COLORS[status as keyof typeof STATUS_COLORS] || '#9E9E9E'};"
          ></div>
          <span class="text-sm">
            {formatStatusLabel(status)}: {count}
          </span>
        </div>
      {/each}
    </div>
  {/if}
</div>

<style>
  /* Add any component-specific styles here */
</style>

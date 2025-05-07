<script lang="ts">
  export let type: 'info' | 'success' | 'warning' | 'error' = 'info';
  export let message: string;
  export let dismissible: boolean = false;
  
  let visible = true;
  
  // Get alert classes based on type
  const getAlertClasses = (type: 'info' | 'success' | 'warning' | 'error') => {
    switch (type) {
      case 'success':
        return 'bg-green-100 border-green-400 text-green-700';
      case 'warning':
        return 'bg-yellow-100 border-yellow-400 text-yellow-700';
      case 'error':
        return 'bg-red-100 border-red-400 text-red-700';
      case 'info':
      default:
        return 'bg-blue-100 border-blue-400 text-blue-700';
    }
  };
  
  // Get alert icon based on type
  const getAlertIcon = (type: 'info' | 'success' | 'warning' | 'error') => {
    switch (type) {
      case 'success':
        return 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z';
      case 'warning':
        return 'M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z';
      case 'error':
        return 'M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z';
      case 'info':
      default:
        return 'M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z';
    }
  };
  
  // Dismiss the alert
  const dismiss = () => {
    visible = false;
  };
</script>

{#if visible}
  <div class="border px-4 py-3 rounded relative mb-4 {getAlertClasses(type)}" role="alert">
    <div class="flex">
      <div class="flex-shrink-0">
        <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={getAlertIcon(type)}></path>
        </svg>
      </div>
      <div>
        <p>{message}</p>
      </div>
    </div>
    
    {#if dismissible}
      <button 
        class="absolute top-0 bottom-0 right-0 px-4 py-3"
        on:click={dismiss}
      >
        <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
        </svg>
      </button>
    {/if}
  </div>
{/if}

<style>
  /* Add any component-specific styles here */
</style>

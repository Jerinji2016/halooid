<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  
  export let currentPage: number = 1;
  export let totalPages: number = 1;
  export let showPageNumbers: boolean = true;
  export let maxPageNumbers: number = 5;
  
  const dispatch = createEventDispatcher();
  
  // Handle page change
  const changePage = (page: number) => {
    if (page < 1 || page > totalPages || page === currentPage) {
      return;
    }
    
    dispatch('pageChange', { page });
  };
  
  // Generate page numbers to display
  $: pageNumbers = getPageNumbers(currentPage, totalPages, maxPageNumbers);
  
  function getPageNumbers(current: number, total: number, max: number) {
    if (total <= max) {
      return Array.from({ length: total }, (_, i) => i + 1);
    }
    
    const half = Math.floor(max / 2);
    let start = current - half;
    let end = current + half;
    
    if (start < 1) {
      end += (1 - start);
      start = 1;
    }
    
    if (end > total) {
      start -= (end - total);
      end = total;
    }
    
    start = Math.max(start, 1);
    
    return Array.from({ length: end - start + 1 }, (_, i) => start + i);
  }
</script>

<div class="flex justify-center">
  <nav class="inline-flex rounded-md shadow-sm -space-x-px" aria-label="Pagination">
    <!-- Previous Page -->
    <button
      class="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 {currentPage === 1 ? 'cursor-not-allowed opacity-50' : ''}"
      on:click={() => changePage(currentPage - 1)}
      disabled={currentPage === 1}
      aria-label="Previous page"
    >
      <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
        <path fill-rule="evenodd" d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z" clip-rule="evenodd" />
      </svg>
    </button>
    
    <!-- Page Numbers -->
    {#if showPageNumbers}
      {#if pageNumbers[0] > 1}
        <button
          class="relative inline-flex items-center px-4 py-2 border border-gray-300 bg-white text-sm font-medium text-gray-700 hover:bg-gray-50"
          on:click={() => changePage(1)}
        >
          1
        </button>
        
        {#if pageNumbers[0] > 2}
          <span class="relative inline-flex items-center px-4 py-2 border border-gray-300 bg-white text-sm font-medium text-gray-700">
            ...
          </span>
        {/if}
      {/if}
      
      {#each pageNumbers as page}
        <button
          class="relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium {page === currentPage ? 'z-10 bg-blue-50 border-blue-500 text-blue-600' : 'bg-white text-gray-700 hover:bg-gray-50'}"
          on:click={() => changePage(page)}
          aria-current={page === currentPage ? 'page' : undefined}
        >
          {page}
        </button>
      {/each}
      
      {#if pageNumbers[pageNumbers.length - 1] < totalPages}
        {#if pageNumbers[pageNumbers.length - 1] < totalPages - 1}
          <span class="relative inline-flex items-center px-4 py-2 border border-gray-300 bg-white text-sm font-medium text-gray-700">
            ...
          </span>
        {/if}
        
        <button
          class="relative inline-flex items-center px-4 py-2 border border-gray-300 bg-white text-sm font-medium text-gray-700 hover:bg-gray-50"
          on:click={() => changePage(totalPages)}
        >
          {totalPages}
        </button>
      {/if}
    {:else}
      <span class="relative inline-flex items-center px-4 py-2 border border-gray-300 bg-white text-sm font-medium text-gray-700">
        {currentPage} of {totalPages}
      </span>
    {/if}
    
    <!-- Next Page -->
    <button
      class="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 {currentPage === totalPages ? 'cursor-not-allowed opacity-50' : ''}"
      on:click={() => changePage(currentPage + 1)}
      disabled={currentPage === totalPages}
      aria-label="Next page"
    >
      <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
        <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
      </svg>
    </button>
  </nav>
</div>

<style>
  /* Add any component-specific styles here */
</style>

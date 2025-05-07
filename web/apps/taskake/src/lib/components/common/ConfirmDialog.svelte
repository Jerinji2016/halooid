<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { Button } from '$lib/components/ui';
  import Alert from './Alert.svelte';
  import Spinner from './Spinner.svelte';

  export let title: string;
  export let message: string;
  export let confirmText: string = 'Confirm';
  export let cancelText: string = 'Cancel';
  export let confirmVariant: 'primary' | 'danger' = 'primary';
  export let isLoading: boolean = false;
  export let error: string | null = null;

  const dispatch = createEventDispatcher();

  // Handle confirm action
  const handleConfirm = () => {
    dispatch('confirm');
  };

  // Handle cancel action
  const handleCancel = () => {
    dispatch('cancel');
  };
</script>

<div class="fixed inset-0 z-50 overflow-y-auto" aria-labelledby="modal-title" role="dialog" aria-modal="true">
  <div class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
    <!-- Background overlay -->
    <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" aria-hidden="true"></div>

    <!-- Center modal -->
    <span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>

    <!-- Modal panel -->
    <div class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
      <div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
        <div class="sm:flex sm:items-start">
          {#if confirmVariant === 'danger'}
            <div class="mx-auto flex-shrink-0 flex items-center justify-center h-12 w-12 rounded-full bg-red-100 sm:mx-0 sm:h-10 sm:w-10">
              <svg class="h-6 w-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path>
              </svg>
            </div>
          {:else}
            <div class="mx-auto flex-shrink-0 flex items-center justify-center h-12 w-12 rounded-full bg-blue-100 sm:mx-0 sm:h-10 sm:w-10">
              <svg class="h-6 w-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
              </svg>
            </div>
          {/if}

          <div class="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left">
            <h3 class="text-lg leading-6 font-medium text-gray-900" id="modal-title">
              {title}
            </h3>
            <div class="mt-2">
              <p class="text-sm text-gray-500">
                {message}
              </p>
            </div>

            {#if error}
              <div class="mt-4">
                <Alert type="error" message={error} />
              </div>
            {/if}
          </div>
        </div>
      </div>

      <div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
        <Button
          variant={confirmVariant === 'danger' ? 'danger' : 'primary'}
          on:click={handleConfirm}
          disabled={isLoading}
          className="w-full sm:w-auto sm:ml-3"
        >
          {#if isLoading}
            <Spinner size="sm" />
            <span class="ml-2">Loading...</span>
          {:else}
            {confirmText}
          {/if}
        </Button>

        <Button
          variant="secondary"
          on:click={handleCancel}
          disabled={isLoading}
          className="mt-3 w-full sm:mt-0 sm:w-auto"
        >
          {cancelText}
        </Button>
      </div>
    </div>
  </div>
</div>

<style>
  /* Add any component-specific styles here */
</style>

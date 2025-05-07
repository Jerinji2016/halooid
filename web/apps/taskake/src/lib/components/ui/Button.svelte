<script lang="ts">
  export let variant: 'primary' | 'secondary' | 'danger' = 'primary';
  export let disabled: boolean = false;
  export let className: string = '';
</script>

<button
  class="btn btn-{variant} {className}"
  disabled={disabled}
  on:click
>
  <slot />
</button>

<style>
  .btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    padding: 0.5rem 1rem;
    font-weight: 500;
    border-radius: 0.375rem;
    cursor: pointer;
    transition: all 0.2s;
  }
  
  .btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
  
  .btn-primary {
    background-color: #3b82f6;
    color: white;
  }
  
  .btn-primary:hover:not(:disabled) {
    background-color: #2563eb;
  }
  
  .btn-secondary {
    background-color: #e5e7eb;
    color: #1f2937;
  }
  
  .btn-secondary:hover:not(:disabled) {
    background-color: #d1d5db;
  }
  
  .btn-danger {
    background-color: #ef4444;
    color: white;
  }
  
  .btn-danger:hover:not(:disabled) {
    background-color: #dc2626;
  }
</style>

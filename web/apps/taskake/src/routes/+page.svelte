<script>
  let count = 0;
  
  function increment() {
    count += 1;
  }
</script>

<main>
  <h1>Taskake Web App</h1>
  
  <div class="card">
    <button on:click={increment}>
      Count is {count}
    </button>
  </div>
  
  <p>
    This is a simple Svelte app for the Taskake web application.
  </p>
</main>

<style>
  main {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: 100vh;
    padding: 2rem;
    text-align: center;
  }
  
  h1 {
    font-size: 2.5rem;
    margin-bottom: 2rem;
    color: #3b82f6;
  }
  
  .card {
    padding: 2rem;
    border-radius: 0.5rem;
    background-color: #f8fafc;
    box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
    margin-bottom: 2rem;
  }
  
  button {
    padding: 0.75rem 1.5rem;
    font-size: 1rem;
    font-weight: 500;
    color: white;
    background-color: #3b82f6;
    border: none;
    border-radius: 0.375rem;
    cursor: pointer;
    transition: background-color 0.2s;
  }
  
  button:hover {
    background-color: #2563eb;
  }
  
  p {
    max-width: 600px;
    line-height: 1.6;
    color: #4b5563;
  }
</style>

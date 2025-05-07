<script>
  import '../app.css';
</script>

<div class="app">
  <header>
    <nav>
      <a href="/" class="logo">Taskake</a>
      <div class="nav-links">
        <a href="/">Dashboard</a>
        <a href="/tasks">Tasks</a>
        <a href="/projects">Projects</a>
      </div>
    </nav>
  </header>

  <main>
    <slot />
  </main>

  <footer>
    <p>© 2025 Halooid - Taskake</p>
  </footer>
</div>

<style>
  .app {
    display: flex;
    flex-direction: column;
    min-height: 100vh;
  }

  header {
    background-color: #ffffff;
    border-bottom: 1px solid #e5e7eb;
    padding: 1rem 2rem;
  }

  nav {
    display: flex;
    justify-content: space-between;
    align-items: center;
    max-width: 1200px;
    margin: 0 auto;
    width: 100%;
  }

  .logo {
    font-size: 1.5rem;
    font-weight: 700;
    color: #3b82f6;
    text-decoration: none;
  }

  .nav-links {
    display: flex;
    gap: 1.5rem;
  }

  .nav-links a {
    color: #4b5563;
    text-decoration: none;
    font-weight: 500;
    transition: color 0.2s;
  }

  .nav-links a:hover {
    color: #3b82f6;
  }

  main {
    flex: 1;
    width: 100%;
  }

  footer {
    background-color: #f8fafc;
    border-top: 1px solid #e5e7eb;
    padding: 1.5rem;
    text-align: center;
    color: #6b7280;
  }
</style>

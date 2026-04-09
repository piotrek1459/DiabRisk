<script lang="ts">
  export let onLoginSuccess: (user: any) => void;
  export let onSwitchToRegister: () => void;

  let email = "";
  let password = "";
  let loading = false;
  let error: string | null = null;

  async function handleLogin() {
    loading = true;
    error = null;

    try {
      const res = await fetch("/auth/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        credentials: "include",
        body: JSON.stringify({ email, password })
      });

      if (!res.ok) {
        const data = await res.json();
        throw new Error(data.error || `Login failed: ${res.status}`);
      }

      const data = await res.json();
      onLoginSuccess(data.user);
    } catch (e: any) {
      error = e.message || "Login failed";
    } finally {
      loading = false;
    }
  }

  function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    handleLogin();
  }
</script>

<div class="login-card card">
  <div class="login-header">
    <h2>🔐 Sign In</h2>
    <p>Welcome back to DiabRisk</p>
  </div>

  <form on:submit={handleSubmit}>
    <div class="form-group">
      <label for="email">Email Address</label>
      <input
        id="email"
        type="email"
        placeholder="your.email@example.com"
        bind:value={email}
        required
        disabled={loading}
      />
    </div>

    <div class="form-group">
      <label for="password">Password</label>
      <input
        id="password"
        type="password"
        placeholder="••••••••"
        bind:value={password}
        required
        disabled={loading}
      />
    </div>

    {#if error}
      <div class="error-message">
        <span>⚠️ {error}</span>
      </div>
    {/if}

    <button type="submit" class="submit-btn" disabled={loading}>
      {#if loading}
        <span class="spinner"></span>
        Signing in...
      {:else}
        Sign In
      {/if}
    </button>
  </form>

  <div class="divider">or</div>

  <button class="register-link" on:click={onSwitchToRegister} disabled={loading}>
    Don't have an account? Create one
  </button>
</div>

<style>
  .login-card {
    max-width: 400px;
    margin: 2rem auto;
    padding: 2rem;
  }

  .login-header {
    text-align: center;
    margin-bottom: 2rem;
  }

  .login-header h2 {
    margin: 0;
    font-size: 1.8rem;
    color: #333;
  }

  .login-header p {
    margin: 0.5rem 0 0 0;
    color: #666;
    font-size: 0.95rem;
  }

  form {
    display: flex;
    flex-direction: column;
    gap: 1.2rem;
  }

  .form-group {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  label {
    font-weight: 600;
    color: #333;
    font-size: 0.9rem;
  }

  input {
    padding: 0.75rem;
    border: 2px solid #e0e0e0;
    border-radius: 6px;
    font-size: 1rem;
    transition: all 0.3s ease;
  }

  input:focus {
    outline: none;
    border-color: #667eea;
    box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
  }

  input:disabled {
    background-color: #f5f5f5;
    cursor: not-allowed;
  }

  .submit-btn {
    padding: 0.75rem 1rem;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    border: none;
    border-radius: 6px;
    font-size: 1rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.3s ease;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
  }

  .submit-btn:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 8px 16px rgba(102, 126, 234, 0.3);
  }

  .submit-btn:disabled {
    opacity: 0.7;
    cursor: not-allowed;
  }

  .spinner {
    display: inline-block;
    width: 14px;
    height: 14px;
    border: 2px solid rgba(255, 255, 255, 0.3);
    border-top-color: white;
    border-radius: 50%;
    animation: spin 0.6s linear infinite;
  }

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }

  .error-message {
    background-color: #fee;
    color: #c33;
    padding: 0.75rem;
    border-radius: 6px;
    font-size: 0.9rem;
  }

  .divider {
    text-align: center;
    color: #999;
    font-size: 0.85rem;
    margin: 1.5rem 0;
    position: relative;
  }

  .divider::before,
  .divider::after {
    content: "";
    position: absolute;
    top: 50%;
    width: 45%;
    height: 1px;
    background-color: #e0e0e0;
  }

  .divider::before {
    left: 0;
  }

  .divider::after {
    right: 0;
  }

  .register-link {
    background-color: #f5f5f5;
    color: #667eea;
    border: 2px solid #e0e0e0;
    padding: 0.75rem;
    border-radius: 6px;
    font-size: 0.9rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.3s ease;
  }

  .register-link:hover:not(:disabled) {
    background-color: #efefef;
    border-color: #667eea;
  }

  .register-link:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  :global(.card) {
    background: white;
    border-radius: 12px;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.07), 0 1px 2px rgba(0, 0, 0, 0.05);
    padding: 1.5rem;
  }
</style>

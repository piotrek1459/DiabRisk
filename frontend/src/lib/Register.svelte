<script lang="ts">
  export let onRegisterSuccess: (user: any) => void;
  export let onSwitchToLogin: () => void;

  let email = "";
  let password = "";
  let confirmPassword = "";
  let fullName = "";
  let loading = false;
  let error: string | null = null;

  async function handleRegister() {
    if (password !== confirmPassword) {
      error = "Passwords do not match";
      return;
    }

    if (password.length < 6) {
      error = "Password must be at least 6 characters";
      return;
    }

    loading = true;
    error = null;

    try {
      const res = await fetch("/auth/register", {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        credentials: "include",
        body: JSON.stringify({ email, password, full_name: fullName })
      });

      if (!res.ok) {
        const data = await res.json();
        throw new Error(data.error || `Registration failed: ${res.status}`);
      }

      const data = await res.json();
      onRegisterSuccess(data.user);
    } catch (e: any) {
      error = e.message || "Registration failed";
    } finally {
      loading = false;
    }
  }

  function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    handleRegister();
  }
</script>

<div class="register-card card">
  <div class="register-header">
    <h2>📝 Create Account</h2>
    <p>Join DiabRisk to assess your diabetes risk</p>
  </div>

  <form on:submit={handleSubmit}>
    <div class="form-group">
      <label for="fullname">Full Name</label>
      <input
        id="fullname"
        type="text"
        placeholder="John Doe"
        bind:value={fullName}
        required
        disabled={loading}
      />
    </div>

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
        placeholder="At least 6 characters"
        bind:value={password}
        required
        disabled={loading}
      />
    </div>

    <div class="form-group">
      <label for="confirmPassword">Confirm Password</label>
      <input
        id="confirmPassword"
        type="password"
        placeholder="••••••••"
        bind:value={confirmPassword}
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
        Creating account...
      {:else}
        Create Account
      {/if}
    </button>
  </form>

  <div class="divider">or</div>

  <button class="login-link" on:click={onSwitchToLogin} disabled={loading}>
    Already have an account? Sign in
  </button>
</div>

<style>
  .register-card {
    max-width: 400px;
    margin: 2rem auto;
    padding: 2rem;
  }

  .register-header {
    text-align: center;
    margin-bottom: 2rem;
  }

  .register-header h2 {
    margin: 0;
    font-size: 1.8rem;
    color: #333;
  }

  .register-header p {
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

  .login-link {
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

  .login-link:hover:not(:disabled) {
    background-color: #efefef;
    border-color: #667eea;
  }

  .login-link:disabled {
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

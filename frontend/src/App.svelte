<script lang="ts">
  import { onMount } from "svelte";
  import Login from "./lib/Login.svelte";
  import Register from "./lib/Register.svelte";

  const API_BASE = "";

  let features = {
    HighBP: 0,
    HighChol: 0,
    CholCheck: 0,
    BMI: 25.5,
    Smoker: 0,
    Stroke: 0,
    HeartDiseaseorAttack: 0,
    PhysActivity: 0,
    Fruits: 0,
    Veggies: 0,
    HvyAlcoholConsump: 0,
    AnyHealthcare: 0,
    NoDocbcCost: 0,
    GenHlth: 3,
    MentHlth: 2,
    PhysHlth: 1,
    DiffWalk: 0,
    Sex: 1,
    Age: 45,
    Education: 4,
    Income: 5
  };

  const fieldLabels = {
    HighBP: "High Blood Pressure",
    HighChol: "High Cholesterol",
    CholCheck: "Cholesterol Check",
    BMI: "Body Mass Index (BMI)",
    Smoker: "Smoker",
    Stroke: "History of Stroke",
    HeartDiseaseorAttack: "Heart Disease or Attack",
    PhysActivity: "Physical Activity",
    Fruits: "Consumes Fruits",
    Veggies: "Consumes Vegetables",
    HvyAlcoholConsump: "Heavy Alcohol Consumption",
    AnyHealthcare: "Has Healthcare Coverage",
    NoDocbcCost: "Could Not See Doctor Due to Cost",
    GenHlth: "General Health (1-5)",
    MentHlth: "Mental Health (Days)",
    PhysHlth: "Physical Health (Days)",
    DiffWalk: "Difficulty Walking",
    Sex: "Sex (0=Female, 1=Male)",
    Age: "Age",
    Education: "Education Level (1-6)",
    Income: "Income Level (1-8)"
  };

  let loading = false;
  let error: string | null = null;
  let result: any = null;
  let user: any = null;
  let checkingAuth = true;
  let authMode: "login" | "register" = "login";

  onMount(async () => {
    await checkSession();
  });

  async function checkSession() {
    try {
      const res = await fetch(`/auth/session`, {
        credentials: "include"
      });
      if (res.ok) {
        user = await res.json();
      }
    } catch (e) {
      // Not logged in, that's ok
    } finally {
      checkingAuth = false;
    }
  }

  function handleLoginSuccess(userData: any) {
    user = userData;
  }

  function handleRegisterSuccess(userData: any) {
    user = userData;
  }

  function switchToLogin() {
    authMode = "login";
  }

  function switchToRegister() {
    authMode = "register";
  }

  async function logout() {
    try {
      await fetch(`/auth/logout`, {
        method: "POST",
        credentials: "include"
      });
      user = null;
      result = null;
      authMode = "login";
    } catch (e) {
      console.error("Logout failed:", e);
    }
  }

  async function submitRisk() {
    loading = true;
    error = null;
    result = null;

    try {
      const res = await fetch(`/api/risk`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        credentials: "include",
        body: JSON.stringify({ features })
      });

      if (!res.ok) {
        const text = await res.text();
        throw new Error(text || `HTTP ${res.status}`);
      }

      result = await res.json();
    } catch (e: any) {
      error = e.message ?? "Unknown error";
    } finally {
      loading = false;
    }
  }

  function getInputType(key: string) {
    return key === "BMI" ? "number" : "number";
  }

  function getInputStep(key: string) {
    return key === "BMI" ? "0.1" : "1";
  }
</script>

<main class="app">
  <div class="header">
    <div class="header-left">
      <h1>🩺 DiabRisk</h1>
      <p class="subtitle">Diabetes Risk Assessment Tool</p>
    </div>
    <div class="header-right">
      {#if !checkingAuth && user}
        <div class="user-info">
          <span class="user-name">{user.full_name || user.email}</span>
          <button on:click={logout} class="logout-btn">Sign Out</button>
        </div>
      {/if}
    </div>
  </div>

  {#if checkingAuth}
    <div class="loading-container">
      <div class="spinner"></div>
      <p>Loading...</p>
    </div>
  {:else if !user}
    <div class="auth-container">
      {#if authMode === "login"}
        <Login onLoginSuccess={handleLoginSuccess} onSwitchToRegister={switchToRegister} />
      {:else}
        <Register onRegisterSuccess={handleRegisterSuccess} onSwitchToLogin={switchToLogin} />
      {/if}
    </div>
  {:else}
    <div class="intro-card card">
      <p>Fill in your health information below to estimate your diabetes risk. All fields are required for an accurate assessment.</p>
    </div>

    <form on:submit|preventDefault={submitRisk} class="form-card card">
      <h2>Health Information</h2>
      <div class="form-grid">
        {#each Object.keys(features) as key}
          <label class="form-field">
            <span class="field-label">{fieldLabels[key]}</span>
            <input
              type={getInputType(key)}
              step={getInputStep(key)}
              bind:value={features[key]}
              required
            />
          </label>
        {/each}
      </div>

      <button type="submit" class="submit-btn" disabled={loading}>
        {#if loading}
          <span class="spinner"></span>
          Calculating...
        {:else}
          Estimate Risk
        {/if}
      </button>
    </form>

    {#if error}
      <div class="error card">
        <strong>⚠️ Error:</strong> {error}
      </div>
    {/if}

    {#if result}
      <section class="result-card card">
        <h2>📊 Risk Assessment Results</h2>
        <div class="risk-score" class:high={result.Category === 'high'} class:medium={result.Category === 'medium'} class:low={result.Category === 'low'}>
          <div class="score-value">{(result.RiskPercent * 100).toFixed(1)}%</div>
          <div class="score-label">Risk Level: <strong>{result.Category.toUpperCase()}</strong></div>
        </div>
        <div class="message">
          <p>{result.Message}</p>
        </div>
      </section>
    {/if}
  {/if}

  <p class="disclaimer">
    ⚕️ <strong>Medical Disclaimer:</strong> This is an educational demonstration only. Not a medical device and not medical advice. Please consult with healthcare professionals for actual medical assessments.
  </p>
</main>

<style>
  :global(body) {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    min-height: 100vh;
    margin: 0;
  }

  .app {
    max-width: 800px;
    margin: 0 auto;
    padding: 2rem 1rem 3rem;
    font-family: system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
  }

  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2rem;
    color: white;
  }

  .header-left {
    text-align: left;
  }

  .header-left h1 {
    margin: 0;
    font-size: 2rem;
  }

  .subtitle {
    margin: 0.5rem 0 0 0;
    font-size: 0.95rem;
    opacity: 0.9;
  }

  .header-right {
    display: flex;
    align-items: center;
  }

  .user-info {
    display: flex;
    align-items: center;
    gap: 1rem;
    background: rgba(255, 255, 255, 0.2);
    padding: 0.75rem 1.5rem;
    border-radius: 50px;
    backdrop-filter: blur(10px);
  }

  .user-name {
    color: white;
    font-weight: 500;
  }

  .logout-btn {
    padding: 0.5rem 1rem;
    border: none;
    border-radius: 6px;
    background: rgba(255, 255, 255, 0.3);
    color: white;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
  }

  .logout-btn:hover {
    background: rgba(255, 255, 255, 0.5);
  }

  .loading-container {
    text-align: center;
    padding: 4rem 2rem;
    color: white;
  }

  .spinner {
    display: inline-block;
    width: 40px;
    height: 40px;
    border: 4px solid rgba(255, 255, 255, 0.3);
    border-top-color: white;
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
    margin-bottom: 1rem;
  }

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }

  .auth-container {
    width: 100%;
  }

  :global(.card) {
    background: white;
    border-radius: 12px;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.07), 0 1px 2px rgba(0, 0, 0, 0.05);
    padding: 1.5rem;
  }

  .intro-card {
    margin-bottom: 2rem;
    text-align: center;
    color: #666;
  }

  .form-card {
    margin-bottom: 2rem;
  }

  .form-card h2 {
    margin-top: 0;
    color: #1f2937;
  }

  .form-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 1.5rem;
    margin-bottom: 1.5rem;
  }

  .form-field {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .field-label {
    font-weight: 600;
    color: #333;
    font-size: 0.9rem;
  }

  .form-field input {
    padding: 0.75rem;
    border: 2px solid #e0e0e0;
    border-radius: 6px;
    font-size: 1rem;
    transition: all 0.3s ease;
  }

  .form-field input:focus {
    outline: none;
    border-color: #667eea;
    box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
  }

  .submit-btn {
    padding: 0.75rem 1.5rem;
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
    width: 100%;
  }

  .submit-btn:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 8px 16px rgba(102, 126, 234, 0.3);
  }

  .submit-btn:disabled {
    opacity: 0.7;
    cursor: not-allowed;
  }

  .error {
    margin-bottom: 2rem;
    background-color: #fee;
    color: #c33;
    border-left: 4px solid #c33;
  }

  .result-card {
    text-align: center;
    margin-bottom: 2rem;
  }

  .result-card h2 {
    margin-top: 0;
    color: #1f2937;
  }

  .risk-score {
    padding: 2rem;
    border-radius: 8px;
    margin-bottom: 1.5rem;
  }

  .risk-score.high {
    background-color: #fee;
    border: 2px solid #dc2626;
  }

  .risk-score.medium {
    background-color: #fef3c7;
    border: 2px solid #f59e0b;
  }

  .risk-score.low {
    background-color: #f0fdf4;
    border: 2px solid #16a34a;
  }

  .score-value {
    font-size: 2.5rem;
    font-weight: bold;
    margin-bottom: 0.5rem;
  }

  .risk-score.high .score-value {
    color: #dc2626;
  }

  .risk-score.medium .score-value {
    color: #f59e0b;
  }

  .risk-score.low .score-value {
    color: #16a34a;
  }

  .score-label {
    font-size: 1.1rem;
    color: #666;
  }

  .message {
    text-align: left;
    padding: 1rem;
    background-color: white;
    border-radius: 6px;
    line-height: 1.6;
  }

  .disclaimer {
    text-align: center;
    color: rgba(255, 255, 255, 0.9);
    font-size: 0.85rem;
    margin-top: 2rem;
    line-height: 1.5;
  }
</style>

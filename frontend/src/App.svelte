<script lang="ts">
  import { onMount } from "svelte";

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
    <h1>🩺 DiabRisk</h1>
    <p class="subtitle">Diabetes Risk Assessment Tool</p>
  </div>

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
    text-align: center;
    margin-bottom: 2rem;
    color: white;
  }

  h1 {
    font-size: 3rem;
    margin: 0;
    text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.3);
  }

  .subtitle {
    font-size: 1.2rem;
    margin: 0.5rem 0 0;
    opacity: 0.95;
  }

  .card {
    border-radius: 20px;
    padding: 2rem;
    margin: 1.5rem 0;
    box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);
    background: white;
    color: #1f2937;
  }

  .intro-card {
    background: #f0f9ff;
    border-left: 4px solid #3b82f6;
  }

  .intro-card p {
    margin: 0;
    line-height: 1.6;
  }

  .form-card h2 {
    margin-top: 0;
    color: #1f2937;
    font-size: 1.5rem;
  }

  .form-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
    gap: 1.25rem;
    margin-bottom: 2rem;
  }

  .form-field {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .field-label {
    font-size: 0.9rem;
    font-weight: 600;
    color: #374151;
  }

  input[type="number"] {
    padding: 0.75rem;
    border-radius: 10px;
    border: 2px solid #e5e7eb;
    background: #f9fafb;
    color: #1f2937;
    font-size: 1rem;
    transition: all 0.2s;
  }

  input[type="number"]:focus {
    outline: none;
    border-color: #3b82f6;
    background: white;
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
  }

  .submit-btn {
    width: 100%;
    border: none;
    padding: 1rem 2rem;
    border-radius: 12px;
    cursor: pointer;
    font-weight: 700;
    font-size: 1.1rem;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    transition: all 0.3s;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
  }

  .submit-btn:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 10px 25px rgba(102, 126, 234, 0.4);
  }

  .submit-btn:disabled {
    opacity: 0.7;
    cursor: not-allowed;
  }

  .spinner {
    width: 16px;
    height: 16px;
    border: 2px solid rgba(255, 255, 255, 0.3);
    border-top-color: white;
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  .error {
    background: #fef2f2;
    border-left: 4px solid #ef4444;
    color: #991b1b;
  }

  .result-card {
    background: linear-gradient(135deg, #f0f9ff 0%, #e0f2fe 100%);
    border: 2px solid #3b82f6;
  }

  .result-card h2 {
    margin-top: 0;
    color: #1e40af;
  }

  .risk-score {
    text-align: center;
    padding: 2rem;
    border-radius: 16px;
    margin: 1.5rem 0;
    background: white;
  }

  .risk-score.high {
    background: linear-gradient(135deg, #fee2e2 0%, #fecaca 100%);
    border: 2px solid #ef4444;
  }

  .risk-score.medium {
    background: linear-gradient(135deg, #fef3c7 0%, #fde68a 100%);
    border: 2px solid #f59e0b;
  }

  .risk-score.low {
    background: linear-gradient(135deg, #d1fae5 0%, #a7f3d0 100%);
    border: 2px solid #10b981;
  }

  .score-value {
    font-size: 3.5rem;
    font-weight: 800;
    line-height: 1;
    margin-bottom: 0.5rem;
  }

  .score-label {
    font-size: 1.1rem;
    color: #374151;
  }

  .message {
    background: white;
    padding: 1.5rem;
    border-radius: 12px;
    border-left: 4px solid #3b82f6;
  }

  .message p {
    margin: 0;
    font-size: 1.05rem;
    line-height: 1.6;
  }

  .disclaimer {
    font-size: 0.85rem;
    color: white;
    text-align: center;
    margin-top: 2rem;
    padding: 1rem;
    background: rgba(0, 0, 0, 0.2);
    border-radius: 12px;
    line-height: 1.6;
  }

  @media (max-width: 640px) {
    h1 {
      font-size: 2rem;
    }

    .form-grid {
      grid-template-columns: 1fr;
    }
  }
</style>

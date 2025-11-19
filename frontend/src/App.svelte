<script lang="ts">
  import { onMount } from "svelte";

const API_BASE = "";

  let age: number = 40;
  let bmi: number = 25;
  let systolic_bp: number = 120;
  let glucose: number = 90;
  let smoker: boolean = false;

  let loading = false;
  let error: string | null = null;

  type RiskResponse = {
    risk_percent: number;
    category: string;
    message: string;
  };

  let result: RiskResponse | null = null;

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
        body: JSON.stringify({
          age,
          bmi,
          systolic_bp,
          glucose,
          smoker
        })
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
</script>

<main class="app">
  <h1>DiabRisk – demo</h1>
  <p>Fill in basic health data to get a rough, non-medical risk estimate.</p>

  <form on:submit|preventDefault={submitRisk} class="card">
    <label>
      Age
      <input type="number" bind:value={age} min="1" max="120" />
    </label>

    <label>
      BMI
      <input type="number" step="0.1" bind:value={bmi} min="10" max="60" />
    </label>

    <label>
      Systolic blood pressure (mmHg)
      <input type="number" bind:value={systolic_bp} min="80" max="250" />
    </label>

    <label>
      Fasting glucose (mg/dL)
      <input type="number" bind:value={glucose} min="50" max="400" />
    </label>

    <label class="checkbox">
      <input type="checkbox" bind:checked={smoker} />
      Current smoker
    </label>

    <button type="submit" disabled={loading}>
      {#if loading}Calculating...{/if}
      {#if !loading}Estimate risk{/if}
    </button>
  </form>

  {#if error}
    <div class="error">
      <strong>Error:</strong> {error}
    </div>
  {/if}

  {#if result}
    <section class="card">
      <h2>Result</h2>
      <p>
        Estimated risk: <strong>{result.risk_percent.toFixed(1)}%</strong> ({result.category})
      </p>
      <p>{result.message}</p>
    </section>
  {/if}

  <p class="disclaimer">
    This is an educational demo only – not a medical device and not medical advice.
  </p>
</main>

<style>
  .app {
    max-width: 640px;
    margin: 0 auto;
    padding: 2rem 1rem 3rem;
    font-family: system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
  }
  h1 {
    margin-bottom: 0.5rem;
  }
  .card {
    border-radius: 16px;
    padding: 1.5rem;
    margin: 1rem 0;
    box-shadow: 0 4px 18px rgba(0, 0, 0, 0.06);
    background: #111827;
  }
  form.card {
    display: grid;
    gap: 1rem;
  }
  label {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
    font-size: 0.95rem;
  }
  input[type="number"] {
    padding: 0.4rem 0.6rem;
    border-radius: 0.5rem;
    border: 1px solid #374151;
    background: #020617;
    color: white;
  }
  .checkbox {
    flex-direction: row;
    align-items: center;
    gap: 0.5rem;
  }
  button {
    border: none;
    padding: 0.6rem 1rem;
    border-radius: 999px;
    cursor: pointer;
    font-weight: 600;
  }
  .error {
    background: #7f1d1d;
    padding: 1rem;
    border-radius: 0.75rem;
  }
  .disclaimer {
    font-size: 0.8rem;
    opacity: 0.7;
    margin-top: 1.5rem;
  }
</style>

<script lang="ts">
  export let items = [];
  export let loading = false;
  export let error: string | null = null;

  function formatDate(value: string) {
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) {
      return value;
    }

    return date.toLocaleString();
  }
</script>

<section class="history-card card">
  <div class="history-header">
    <h2>Assessment History</h2>
    <p>Your recent diabetes risk assessments appear here.</p>
  </div>

  {#if loading}
    <p class="history-state">Loading history...</p>
  {:else if error}
    <p class="history-state error">{error}</p>
  {:else if items.length === 0}
    <p class="history-state">No assessments yet. Your first result will appear here.</p>
  {:else}
    <div class="history-list">
      {#each items as item}
        <article class="history-item">
          <div class="history-item-top">
            <span class="history-date">{formatDate(item.created_at)}</span>
            <span class={`history-badge ${item.category}`}>{item.category.toUpperCase()}</span>
          </div>
          <div class="history-score">{(item.risk_percent * 100).toFixed(1)}%</div>
          <p class="history-message">{item.message}</p>
        </article>
      {/each}
    </div>
  {/if}
</section>

<style>
  .history-card {
    margin-bottom: 2rem;
  }

  .history-header h2 {
    margin: 0;
    color: #1f2937;
  }

  .history-header p {
    margin: 0.5rem 0 1.5rem;
    color: #6b7280;
  }

  .history-list {
    display: grid;
    gap: 1rem;
  }

  .history-item {
    border: 1px solid #e5e7eb;
    border-radius: 12px;
    padding: 1rem;
    background: #f9fafb;
  }

  .history-item-top {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 1rem;
    margin-bottom: 0.75rem;
  }

  .history-date {
    color: #6b7280;
    font-size: 0.9rem;
  }

  .history-badge {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    min-width: 88px;
    padding: 0.35rem 0.75rem;
    border-radius: 999px;
    font-size: 0.8rem;
    font-weight: 700;
    letter-spacing: 0.03em;
  }

  .history-badge.high {
    background: #fee2e2;
    color: #b91c1c;
  }

  .history-badge.medium {
    background: #fef3c7;
    color: #b45309;
  }

  .history-badge.low {
    background: #dcfce7;
    color: #15803d;
  }

  .history-score {
    font-size: 2rem;
    font-weight: 700;
    color: #111827;
    margin-bottom: 0.5rem;
  }

  .history-message {
    margin: 0;
    color: #4b5563;
    line-height: 1.5;
  }

  .history-state {
    margin: 0;
    color: #6b7280;
  }

  .history-state.error {
    color: #b91c1c;
  }
</style>

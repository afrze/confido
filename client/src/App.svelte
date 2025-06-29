<script lang="ts">
  let count = $state(0);
</script>


<main>
  <div>
      Welcome to Confido
    </div>
    <button onclick={() => count++}>
      count: {count}
    </button>
</main>

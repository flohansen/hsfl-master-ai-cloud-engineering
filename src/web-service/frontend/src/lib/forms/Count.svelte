<script lang="ts">
    export let count: number = 1
</script>

<section>
    <p class="text-gray-dark text-sm font-medium mb-2 lg:mb-3 lg:text-base">
        Anzahl auswählen:
    </p>

    <div class="flex items-center gap-x-2 lg:gap-x-3">
        <button
            aria-label="Anzahl des Produktes um Einen herabsetzen"
            on:click={() => (count > 1) ? count-- : null}
            disabled="{count < 2}"
            class="rounded-full bg-green-light/25 w-8 h-8 flex items-center justify-center transition-all ease-in-out duration-300 hover:bg-green-light/50 group disabled:bg-gray-light">
            <span class="text-green-dark font-semibold text-xl group-disabled:text-gray-dark">-</span>
        </button>

        <span aria-live="polite">
            {count}
        </span>

        <button
            aria-label="Anzahl des Produktes um Einen erhöhen"
            on:click={() => count++}
            class="rounded-full bg-green-light/25 w-8 h-8 flex items-center justify-center transition-all ease-in-out duration-300 hover:bg-green-light/50">
            <span class="text-green-dark font-semibold text-xl">+</span>
        </button>
    </div>
</section>
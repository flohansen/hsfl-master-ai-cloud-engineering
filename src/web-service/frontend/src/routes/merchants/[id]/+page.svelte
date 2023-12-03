<script lang="ts">
    import BackLink from "$lib/general/BackLink.svelte";
    import ProductListItem from "$lib/merchant/ProductListItem.svelte";

    interface Data {
        merchant: { name: string },
        prices: { productId: number, price: number }[],
    }

    export let data: Data;
</script>

<header>
    <BackLink url="/merchants" label="Alle Supermärkte" />
</header>

<main>
    <h1 class="px-5 text-lg mb-4 font-semibold lg:text-xl xl:text-2xl">
        {data.merchant.name}
    </h1>
    <ul class="px-5 mt-4 grid grid-cols-1 gap-y-4 lg:gap-y-6 lg:mt-6">
        {#if data.prices}
            {#each data.prices as price}
                <ProductListItem price={price} />
            {/each}
        {:else}
            <p>Es konnten keine Daten geladen werden.</p>
        {/if}
    </ul>
</main>

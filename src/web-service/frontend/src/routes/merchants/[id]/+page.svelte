<script lang="ts">
    import BackLink from "$lib/general/BackLink.svelte";
    import ProductListItem from "$lib/merchant/ProductListItem.svelte";
    import UpdateOrCreateModal from "$lib/merchant/UpdateOrCreateModal.svelte";

    interface Price {
        userId: number,
        productId: number,
        price: number
    }

    interface Data {
        merchant: { id: number, name: string },
        prices: Price[],
        products: { id: number, description: string, ean: number }[],
    }

    function findPriceByProductId(productId: number): Price | undefined {
        return data.prices.find(price => price.productId === productId);
    }

    function isCurrentUserCurrentMerchant(): boolean {
        return sessionStorage.getItem('user_id') ?
            data.merchant.id.toString() == sessionStorage.getItem('user_id')
            : false;
    }

    export let data: Data;
    let isOpen: boolean = false;
</script>

<header>
    <BackLink
        url="/merchants"
        label="Alle Supermärkte" />
</header>

<main>
    <div class="mb-4 px-5 flex flex-wrap justify-between items-center gap-x-4">
        <h1 class="text-lg font-semibold lg:text-xl xl:text-2xl">
            {data.merchant.name}
        </h1>

        {#if isCurrentUserCurrentMerchant()}
            <button
                on:click={() => isOpen = ! isOpen}
                aria-label="Neues Produkt erstellen"
                class="rounded-full bg-green-light w-8 h-8 flex items-center justify-center transition-all ease-in-out duration-300 cursor-pointer hover:bg-green-light/75">
                <span class="text-white font-semibold text-xl">+</span>
            </button>
        {/if}
    </div>

    <ul class="px-5 mt-4 grid grid-cols-1 gap-y-4 lg:gap-y-6 lg:mt-6">
        {#if data.products.length > 0}
            {#each data.products as product}
                <ProductListItem
                    product={product}
                    price={findPriceByProductId(product.id)} />
            {/each}
        {:else}
            <p>Es sind noch keine Produkte vorhanden.</p>
        {/if}
    </ul>

    <UpdateOrCreateModal
        bind:isOpen={isOpen} />
</main>
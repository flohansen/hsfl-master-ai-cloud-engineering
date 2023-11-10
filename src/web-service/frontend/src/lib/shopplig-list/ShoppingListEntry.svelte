<script lang="ts">
    import { onMount } from 'svelte';
    import Checkbox from "$lib/forms/Checkbox.svelte";
    export let productId: number;
    export let productCount: number;
    export let view: ViewState;

    type ViewState = "detailed" | "compressed";

    interface Product {
        id: number,
        description: string,
        count: number,
    }

    interface Price {
        price: number,
    }

    let productData: Product = { id: 0, description: '', count: 0 };
    let priceData: Price = { price: 0 };
    let merchant: string = 'Aldi';

    const apiUrlProduct = `/api/v1/product/${productId}`;
    const apiUrlPrice = `/api/v1/price/${productId}/1`;

    onMount(async () => {
        try {
            const response = await fetch(apiUrlProduct);
            response.ok
                ? productData = await response.json()
                : console.error('Failed to fetch data');
        } catch (error) {
            console.error(error);
        }

        try {
            const response = await fetch(apiUrlPrice);
            response.ok
                ? priceData = await response.json()
                : console.error('Failed to fetch data');
        } catch (error) {
            console.error(error);
        }
    });
</script>

<li class="border-t-2 border-t-gray-light py-3 lg:py-6">
    <div class="flex gap-x-4 items-start justify-between">
        <Checkbox label="{productData.description}" id="{productData.id}"/>
        <span class="mt-0.5 block text-gray-dark text-sm whitespace-nowrap lg:text-base">
            {productCount} Stk.
        </span>
    </div>
    {#if view === 'detailed'}
        <p class="text-gray-dark mt-1 ml-[2.1rem] text-sm flex flex-wrap items-center gap-2 lg:text-sm">
            Am günstigsten bei
            <strong class="text-green-dark font-semibold">{merchant}</strong>für
            <strong class="text-green-dark font-semibold">{priceData.price} €</strong>
        </p>
    {/if}
</li>



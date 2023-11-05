<script lang="ts">
    import { onMount } from 'svelte';
    import Badge from "$lib/shopplig-list/Badge.svelte";
    export let productId: number;

    interface Product {
        id: number,
        description: string,
    }

    interface Price {
        price: number,
    }

    let productData: Product = { id: 0, description: '' };
    let priceData: Price = { price: 0 };

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

<li class="border-t-2 border-t-gray-light py-4 flex justify-between gap-x-4 lg:py-6">
    <div class="flex gap-x-4 items-start">
        <input type="checkbox">
        <div class="-mt-1">
            <h2 class="text-sm font-medium lg:text-base">{productData.description}</h2>
            <p class="text-gray-dark mt-2 text-xs flex flex-wrap items-center gap-2 lg:text-sm">
                Am günstigsten bei: <Badge />
            </p>
        </div>
    </div>
    <span class="block text-gray-dark text-sm lg:text-base whitespace-nowrap">{priceData.price} €</span>
</li>
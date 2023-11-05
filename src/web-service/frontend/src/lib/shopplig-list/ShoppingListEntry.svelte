<script lang="ts">
    import { onMount } from 'svelte';
    export let productId: number;

    interface Product {
        id: number,
        description: string,
    }

    let productData: Product = { id: 0, description: '' };
    const apiUrl = `/api/v1/product/${productId}`;

    onMount(async () => {
        try {
            const response = await fetch(apiUrl);
            response.ok
                ? productData = await response.json()
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
            <span class="block text-gray-dark mt-2 text-xs lg:text-sm">Am günstigsten bei:</span>
        </div>
    </div>
    <span class="block text-gray-dark text-sm lg:text-base">Preis €</span>
</li>
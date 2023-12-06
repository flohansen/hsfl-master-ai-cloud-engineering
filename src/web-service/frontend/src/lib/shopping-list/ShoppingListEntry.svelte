<script lang="ts">
    import {createEventDispatcher, onMount} from 'svelte';
    import Checkbox from "$lib/forms/Checkbox.svelte";
    import {handleErrors} from "../../assets/helper/handleErrors";

    type ViewState = "detailed" | "compressed";

    interface Product {
        id: number,
        description: string,
        count: number,
    }

    interface ShoppingListEntry {
        productId: number,
        count: number,
        checked: boolean
    }

    interface Price {
        price: number,
    }

    export let listId: number;
    export let entry: ShoppingListEntry;
    export let view: ViewState = "detailed";

    let productData: Product = { id: 0, description: '', count: 0 };
    let priceData: Price = { price: 0 };
    let merchant: string = 'Aldi';

    const apiUrlProduct = `/api/v1/product/${entry.productId}`;
    const apiUrlPrice = `/api/v1/price/${entry.productId}/2`;

    const dispatch = createEventDispatcher();

    onMount(async () => {
        fetch(apiUrlProduct)
            .then(handleErrors)
            .then(data => productData = data)
            .catch(error => console.error("Failed to fetch data:", error.message));

        fetch(apiUrlPrice)
            .then(handleErrors)
            .then(data => priceData = data)
            .catch(error => console.error("Failed to fetch data:", error.message));
    });

    function updateShoppingListEntry(): void {
        if (! listId || !productData.id) return;

        const apiUrl: string = `/api/v1/shoppinglistentries/${listId}/${productData.id}`;
        const requestOptions = {
            method: "PUT",
            headers: { 'Content-Type': 'application/json' },
            body: `{ "count": ${entry.count}, "checked": ${entry.checked} }`,
        };

        fetch(apiUrl, requestOptions)
            .then(handleErrors)
            .then(()=> location.reload())
            .catch(error => console.error("Failed to fetch data:", error.message));
    }
</script>

<li class="border-t-2 border-t-gray-light py-3 lg:py-6">
    <div class="flex gap-x-4 items-start justify-between">
        <Checkbox
            label="{productData.description}"
            id="{productData.id}"
            bind:checked={entry.checked}
            on:updateShoppingListEntry={updateShoppingListEntry} />
        <span class="mt-0.5 block text-gray-dark text-sm whitespace-nowrap lg:text-base { entry.checked ? 'opacity-50' : '' }">
            {entry.count} Stk.
        </span>
    </div>
    {#if view === 'detailed' && priceData}
        <p class="text-gray-dark mt-1 ml-[2.1rem] text-sm flex flex-wrap items-center gap-2 lg:text-sm { entry.checked ? 'opacity-50' : '' }">
            Am günstigsten bei
            <strong class="text-green-dark font-semibold">{merchant}</strong>für
            <strong class="text-green-dark font-semibold">{priceData.price} €</strong>
        </p>
    {/if}
</li>



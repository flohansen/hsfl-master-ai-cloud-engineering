<script lang="ts">
    import {createEventDispatcher, onMount} from 'svelte';
    import Checkbox from "$lib/forms/Checkbox.svelte";
    import {handleErrors} from "../../assets/helper/handleErrors";
    import Trash from "../../assets/svg/Trash.svelte";
    import {fetchHelper} from "../../assets/helper/fetchHelper";

    type ViewState = "detailed" | "compressed";

    export let listId: number;
    export let entry: { productId: number, count: number, checked?: boolean };
    export let product: { id: number, description: string, ean: number };
    export let view: ViewState = "detailed";

    let price: any;
    let merchant: string = 'Aldi';
    const token: string | null = sessionStorage.getItem('access_token');

    const dispatch = createEventDispatcher();

    onMount(async () => {
        const apiUrl: string = `/api/v1/price/${product.id}/${sessionStorage.getItem('user_id')}`;
        price = await fetchHelper(apiUrl);
    });

    function updateShoppingListEntry(): void {
        if (! listId || ! product.id || ! token) return;

        const apiUrl: string = `/api/v1/shoppinglistentries/${listId}/${product.id}`;
        const requestOptions = {
            method: "PUT",
            headers: { 'Authorization': `Bearer ${token}` },
            body: `{ "count": ${entry.count}, "checked": ${entry.checked} }`,
        };

        fetch(apiUrl, requestOptions)
            .then(handleErrors)
            .then(()=> dispatch('updateCheckedEntriesCount', { state: entry.checked }))
            .catch(error => console.error("Failed to fetch data:", error.message));
    }

    function deleteShoppingListEntry(): void {
        if (! listId || ! product.id || ! token) return;

        const apiUrl: string = `/api/v1/shoppinglistentries/${listId}/${product.id}`;
        const requestOptions = {
            method: "DELETE",
            headers: { 'Authorization': `Bearer ${token}` },
        };

        fetch(apiUrl, requestOptions)
            .then(handleErrors)
            .then(()=> { location.reload(); dispatch('updateCheckedEntriesCount', { state: true }) })
            .catch(error => console.error("Failed to fetch data:", error.message));
    }
</script>

<li class="border-t-2 border-t-gray-light py-3 lg:py-6">
    <div class="flex gap-x-2 items-center justify-between">
        <div class="flex gap-x-2 items-center">
            <Checkbox
                label={product.description}
                id={product.id}
                count={entry.count}
                bind:checked={entry.checked}
                on:updateShoppingListEntry={updateShoppingListEntry} />
        </div>
        <button
            aria-label="Eintrag löschen"
            on:click={deleteShoppingListEntry}
            class="bg-gray-light rounded-full p-2 text-gray-dark transition-all ease-in-out duration-300 hover:bg-gray-dark/25">
            <Trash classes="w-4 h-4 md:w-5 md:h-5" />
        </button>
    </div>
    {#if view === 'detailed' && price}
        <p class="text-gray-dark mt-1 ml-[2.1rem] text-sm flex flex-wrap items-center gap-2 lg:text-sm { entry.checked ? 'opacity-50' : '' }">
            Am günstigsten bei
            <strong class="text-green-dark font-semibold">{merchant}</strong>für
            <strong class="text-green-dark font-semibold">{price.price ?? 0} €</strong>
        </p>
    {/if}
</li>



<script lang="ts">
    import { page } from "$app/stores";
    import Header from "$lib/general/Header.svelte";
    import ShoppingList from "$lib/shopplig-list/ShoppingList.svelte";

    interface Data {
        lists: { id: number, description: string }[]
        headline: string;
    }

    export let data: Data;
</script>

<Header headline="{$page.data.headline}"/>

<main>
    <h2 class="px-5 text-gray-dark text-sm font-medium mt-6 lg:mt-10 lg:text-base">
        Offene Einkaufslisten
    </h2>
    <ul class="px-5 mt-4 grid grid-cols-1 gap-y-4 lg:gap-y-6 lg:mt-6">
        {#if data.lists}
            {#each data.lists as list}
                <ShoppingList description={list.description} id="{list.id}"/>
            {/each}
        {:else}
            <p>Es konnten keine Daten geladen werden.</p>
        {/if}
    </ul>
</main>
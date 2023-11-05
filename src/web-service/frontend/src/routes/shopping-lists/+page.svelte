<script lang="ts">
    import ShoppingListItem from "$lib/ShoppingListItem.svelte";
    import { page } from "$app/stores";
    import Header from "$lib/Header.svelte";
    export let data;
</script>

<Header headline="{$page.data.headline}"/>

<main class="mt-8 sm:ml-20 md:ml-24 lg:max-w-4xl lg:mx-auto lg:mt-10 xl:max-w-5xl">
    <h2 class="px-5 text-gray-dark text-sm font-medium mt-6 lg:mt-10 lg:text-base">
        Offene Einkaufslisten
    </h2>
    <ul class="px-5 mt-4 grid grid-cols-1 gap-y-4 lg:gap-y-6 lg:mt-6">
        {#if data.lists}
            {#each data.lists as item}
                <ShoppingListItem description={item.description} id="{item.id}"/>
            {/each}
        {:else}
            <p>Es konnten keine Daten geladen werden.</p>
        {/if}
    </ul>
</main>
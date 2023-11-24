<script lang="ts">
    import {page} from "$app/stores";
    import ShoppingList from "$lib/shopplig-list/ShoppingList.svelte";

    interface Data {
        lists: { id: number, description: string }[]
        headline: string;
    }

    export let data: Data;
</script>

<header>
    <h1 class="font-bold text-xl md:text-2xl xl:text-3xl">
        {$page.data.headline}
    </h1>

    <a
        href="/shopping-lists/add"
        aria-label="Neue Einkaufsliste erstellen"
        class="rounded-full bg-green-light w-8 h-8 flex items-center justify-center transition-all ease-in-out duration-300 cursor-pointer hover:bg-green-light/75">
        <span class="text-white font-semibold text-xl">+</span>
    </a>
</header>

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
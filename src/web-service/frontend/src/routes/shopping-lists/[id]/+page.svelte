<script lang="ts">
    import ShoppingListEntry from "$lib/shopplig-list/ShoppingListEntry.svelte";

    interface Data {
        list: { description: string },
        entries: { productId: number, checked: boolean }[],
    }

    export let data: Data;
</script>

<header class="px-5 mt-16 flex items-center justify-between sm:ml-20 md:ml-24 lg:max-w-4xl lg:mx-auto xl:max-w-5xl">
    <a href="/shopping-lists" class="flex gap-x-2 items-center text-gray-dark transition-all duration-300 ease-in-out hover:text-green-dark lg:gap-x-4">
        <svg class="w-6 h-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M6.75 15.75L3 12m0 0l3.75-3.75M3 12h18" />
        </svg>
        <p class="text-sm lg:text-base">Alle Einkaufslisten</p>
    </a>
</header>

<main class="mt-8 sm:ml-20 md:ml-24 lg:max-w-4xl lg:mx-auto lg:mt-10 xl:max-w-5xl">
    <div class="mx-5 bg-white rounded-xl p-4 lg:p-6">
        <header class="flex items-center gap-x-4 lg:gap-x-6">
            <figure class="bg-green-light/25 rounded-full w-14 h-14 flex items-center justify-center lg:w-16 lg:h-16">
                <span class="text-2xl lg:text-3xl">🥗</span>
            </figure>
            <div>
                <h1 class="text-lg font-semibold lg:text-xl xl:text-2xl">
                    {data.list.description}
                </h1>
                <p class="text-gray-dark text-sm mt-1">
                    Anzahl der Einträge: {data.entries.length}
                </p>
            </div>
        </header>

        <p class="text-gray-dark text-sm mt-10">Deine Einkaufsliste</p>
        <ul class="mt-4">
            {#each data.entries as entry}
                <ShoppingListEntry productId="{entry.productId}"/>
            {/each}
        </ul>
    </div>
</main>

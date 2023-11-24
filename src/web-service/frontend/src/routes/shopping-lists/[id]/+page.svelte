<script lang="ts">
    import ShoppingListEntry from "$lib/shopplig-list/ShoppingListEntry.svelte";
    import AddEntryModal from "$lib/shopplig-list/AddEntryModal.svelte";
    import ViewSelect from "$lib/shopplig-list/ViewSelect.svelte";

    interface Data {
        list: { id: number, description: string },
        entries: { productId: number, count: number, checked: boolean }[],
    }

    type ViewState = "detailed" | "compressed";

    let view: ViewState;
    export let data: Data;
</script>

<header>
    <a href="/shopping-lists" class="flex gap-x-2 items-center text-gray-dark transition-all duration-300 ease-in-out hover:text-green-dark lg:gap-x-4">
        <svg class="w-6 h-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M6.75 15.75L3 12m0 0l3.75-3.75M3 12h18" />
        </svg>
        <p class="text-sm lg:text-base">Alle Einkaufslisten</p>
    </a>
</header>

<main>
    <div class="mx-5 bg-white rounded-xl p-4 lg:p-6">
        <section class="flex items-center gap-x-4 lg:gap-x-6">
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
        </section>

        <ViewSelect bind:view={view}/>

        <p class="text-gray-dark text-sm">Deine Einkaufsliste</p>
        <ul class="mt-4">
            {#each data.entries as entry}
                <ShoppingListEntry
                    view="{view}"
                    productId="{entry.productId}"
                    productCount="{entry.count}"/>
            {/each}
        </ul>

        <AddEntryModal listId="{data.list.id}" currentEntries="{data.entries}"/>
    </div>
</main>

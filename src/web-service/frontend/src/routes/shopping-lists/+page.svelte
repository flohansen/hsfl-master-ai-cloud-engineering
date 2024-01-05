<script lang="ts">
    import ShoppingListSection from "$lib/shopping-list/ShoppingListSection.svelte";

    interface List {
        id: number,
        description: string,
        complete?: boolean,
    }

    interface Data {
        completedLists: List[],
        incompleteLists: List[],
        headline: string;
    }

    export let data: Data;
</script>

<header>
    <h1 class="font-bold text-xl md:text-2xl xl:text-3xl">
        {data.headline}
    </h1>

    <a
        href="/shopping-lists/add"
        aria-label="Neue Einkaufsliste erstellen"
        class="rounded-full bg-green-light w-8 h-8 flex items-center justify-center transition-all ease-in-out duration-300 cursor-pointer hover:bg-green-light/75">
        <span class="text-white font-semibold text-xl">+</span>
    </a>
</header>

<main>
    {#if data.completedLists.length === 0 && data.incompleteLists.length === 0}
        <p>Es konnten keine Daten geladen werden.</p>
    {:else}
        <ShoppingListSection
            label="Offene Einkaufslisten"
            lists={data.incompleteLists} />

        <ShoppingListSection
            label="Abgeschlossene Einkaufslisten"
            lists={data.completedLists} />
    {/if}
</main>
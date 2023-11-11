<script lang="ts">
    import List from "../../assets/svg/List.svelte";
    import Tiles from "../../assets/svg/Tiles.svelte";

    export let view: ViewState = 'detailed';
    type ViewState = "detailed" | "compressed";
</script>

<div class="flex items-center mt-10 justify-end">
    <button
        aria-label="Einkaufstliste mit Details anzeigen"
        on:click={() => view  === 'detailed'}
        class="rounded-l-xl py-1 pl-2.5 pr-1 transition-all ease-in-out duration-300 cursor-pointer hover:bg-gray-dark/25 lg:py-2 lg:pl-4 lg:pr-2
            {view === 'detailed' ? 'bg-green-light/50 text-green-dark' : 'bg-gray-light text-black'}">
        <Tiles classes="w-5 h-5"/>
    </button>

    <button
        aria-label="Einkaufstliste ohne Details anzeigen"
        on:click={() => view  === 'compressed'}
        class="rounded-r-xl py-1 pr-2.5 pl-1 transition-all ease-in-out duration-300 cursor-pointer hover:bg-gray-dark/25 lg:py-2 lg:pr-4 lg:pl-2
            {view === 'compressed' ? 'bg-green-light/50 text-green-dark' : 'bg-gray-light text-black'}">
        <List classes="w-5 h-5"/>
    </button>
</div>
<script lang="ts">
    import Add from "../../assets/svg/Add.svelte";
    import Select from "$lib/forms/Select.svelte";
    import Count from "$lib/forms/Count.svelte";
    import Checkmark from "../../assets/svg/Checkmark.svelte";
    import {clickOutside} from "../../assets/helper/clickOutside";
    import {handleErrors} from "../../assets/helper/handleErrors";

    interface NewEntry {
        id: number,
        count: number,
        checked: boolean
    }

    interface ShoppingListEntry {
        productId: number,
        count: number,
    }

    let isOpen = false;
    let entry: NewEntry = { id: 0, count: 0, checked: false };

    export let listId: number;
    export let currentEntries: ShoppingListEntry[];

    function handleClickOutside(): void {
        if (! isOpen) return;
        isOpen = false;
    }

    // Checks if the selected item is already listed in the current shopping list.
    function findExistingListEntry(): ShoppingListEntry | undefined {
        return currentEntries.find((listEntry) => listEntry.productId === entry.id);
    }

    function submit(): void {
        if (! entry.id || ! entry.count) return;

        let existingEntry: ShoppingListEntry | undefined = findExistingListEntry();

        existingEntry
            ? fetchContent("PUT", entry.count + existingEntry.count)
            : fetchContent("POST", entry.count);
    }

    function fetchContent(method: string, count: number): void {
        const apiUrl: string = `/api/v1/shoppinglistentries/${listId}/${entry.id}`
        const requestOptions = {
            method: method,
            headers: { 'Content-Type': 'application/json' },
            body: `{"count": ${count}, "note": "", "checked": false}`,
        };

        fetch(apiUrl, requestOptions)
            .then(handleErrors)
            .then(()=> location.reload())
            .catch(error => console.error("Failed to fetch data:", error.message));
    }
</script>

<button on:click={() => isOpen = ! isOpen} class="-ml-[2px] mt-6 text-green-dark hover:text-green-light flex items-center justify-center gap-x-4">
    <Add classes="w-6 h-6 transition-all ease-in-out duration-300"/>
    <span class="block transition-all ease-in-out duration-300 text-sm lg:text-base">
        Einträge hinzufügen
    </span>
</button>

<div class:hidden={! isOpen} class="bg-black/80 fixed inset-0 w-screen h-screen"></div>

<section
    use:clickOutside
    on:click_outside={handleClickOutside}
    class:hidden={! isOpen}
    class="fixed inset-x-4 h-min top-1/2 -translate-y-1/2 bg-white rounded-xl px-4 py-6 sm:left-28 sm:right-8 lg:inset-x-0 lg:max-w-[52rem] lg:mx-auto lg:px-6 lg:px-10 xl:max-w-[62rem]">
    <h3 class="text-lg font-semibold lg:text-xl">Eintrag hinzufügen:</h3>

    <Select bind:entryId={entry.id}/>
    <Count bind:count={entry.count}/>

    <button
        on:click={submit}
        disabled="{! entry.id || ! entry.count}"
        class="bg-green-light mt-8 mx-auto text-white rounded-xl px-5 py-2 flex items-center justify-center gap-x-2 transition-all ease-in-out duration-300 hover:bg-green-dark disabled:bg-gray-light disabled:text-gray-dark">
        <span class="text-sm lg:text-base">Speichern</span>
        <Checkmark classes="w-5 h-5 stroke-2"/>
    </button>
</section>
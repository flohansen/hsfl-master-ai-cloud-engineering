<script lang="ts">
    import Add from "../../assets/svg/Add.svelte";
    import Select from "$lib/forms/SelectProducts.svelte";
    import Count from "$lib/forms/Count.svelte";
    import {handleErrors} from "../../assets/helper/handleErrors";
    import Modal from "$lib/general/Modal.svelte";

    interface NewEntry {
        id: number,
        count: number,
        checked: boolean
    }

    interface ShoppingListEntry {
        productId: number,
        count: number,
    }

    let isOpen: boolean = false;
    let entry: NewEntry = { id: 0, count: 1, checked: false };

    export let listId: number;
    export let currentEntries: ShoppingListEntry[];

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

<Modal on:submit={submit} bind:isOpen>
    <h3 class="text-lg font-semibold lg:text-xl">Eintrag hinzufügen:</h3>

    <Select bind:justValue={entry.id}/>
    <Count bind:count={entry.count}/>
</Modal>

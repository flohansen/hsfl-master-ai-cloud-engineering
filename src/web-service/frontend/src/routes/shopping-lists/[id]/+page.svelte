<script lang="ts">
    import ShoppingListEntry from "$lib/shopping-list/ShoppingListEntry.svelte";
    import AddEntryModal from "$lib/shopping-list/AddEntryModal.svelte";
    import ViewSelect from "$lib/shopping-list/ViewSelect.svelte";
    import {onMount} from "svelte";
    import {handleErrors} from "../../../assets/helper/handleErrors";
    import BackLink from "$lib/general/BackLink.svelte";

    interface Data {
        list: { id: number, description: string, completed?: boolean },
        entries: { productId: number, count: number, checked?: boolean }[],
    }

    type ViewState = "detailed" | "compressed";

    let view: ViewState;
    let checkedEntriesCount: number = 0;
    export let data: Data;

    onMount(async () => {
        data.entries.forEach(entry => {
            if (entry.checked) {
                checkedEntriesCount++;
            }
        })
    });

    function updateCheckedEntriesCount(event: any): void {
        event.detail.state ? checkedEntriesCount++ : checkedEntriesCount --;

        if (data.entries.length === checkedEntriesCount && ! data.list.completed) {
            data.list.completed = true;
            updateShoppingList();
        }

        if (data.entries.length !== checkedEntriesCount && data.list.completed) {
            data.list.completed = false;
            updateShoppingList();
        }
    }

    function updateShoppingList(): void {
        const userId = 2; // TODO: add real user id

        if (! data.list.id || ! userId) return;

        const apiUrl = `/api/v1/shoppinglist/${data.list.id}/${userId}`;
        const requestOptions = {
            method: "PUT",
            headers: { 'Content-Type': 'application/json' },
            body: `{ "description": "${data.list.description}", "checked": ${data.list.completed} }`,
        };

        fetch(apiUrl, requestOptions)
            .then(handleErrors)
            .catch(error => console.error("Failed to fetch data:", error.message));
    }
</script>

<header>
    <BackLink
        url="/shopping-lists"
        label="Alle Einkaufslisten" />
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
                    listId={data.list.id}
                    view={view}
                    entry={entry}
                    on:updateCheckedEntriesCount={updateCheckedEntriesCount}/>
            {/each}
        </ul>

        {#if ! data.list.completed}
            <AddEntryModal
                listId="{data.list.id}"
                currentEntries="{data.entries}"/>
        {/if}
    </div>
</main>

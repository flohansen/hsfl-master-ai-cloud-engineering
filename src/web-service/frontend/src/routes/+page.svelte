<script lang="ts">
    import Header from "$lib/general/Header.svelte";
    import Euro from "../assets/svg/Euro.svelte";
    import ShoppingList from "$lib/shopping-list/ShoppingList.svelte";
    import {isAuthenticated} from "../store";

    interface Data {
        lists: { id: number, description: string }[],
        headline: string;
    }

    export let data: Data;
</script>

{#if isAuthenticated}
    <Header headline="{data.headline}"/>

    <main>
        <section class="mx-5 mt-8 rounded-2xl bg-blue-light/50 p-5 lg:p-8">
            <figure class="w-12 h-12 flex items-center justify-center rounded-full bg-blue-dark/40">
                <Euro classes="w-7 h-7 text-blue-dark"></Euro>
            </figure>
            <h2 class="text-lg font-semibold mt-5 lg:text-xl xl:text-2xl xl:mt-8">
                Spare bis zu 35% Prozent auf deinen nächsten Supermarkt-Einkauf!
            </h2>
            <p class="mt-3 text-base font-light lg:text-lg">
                Erstelle spielend Einkaufslisten, füge Produkte mühelos hinzu.
                PriceWhisper zeigt dir den günstigsten Supermarkt pro Produkt und
                berechnet automatisch den Gesamtpreis.
            </p>
        </section>

        <section class="mx-5 mt-12 lg:mt-16">
            <h2 class="text-gray-dark text-sm font-medium lg:text-base">
                Vorschau deiner Einkaufslisten:
            </h2>
            <ul class="grid grid-cols-1 gap-y-4 mt-4 lg:gap-y-6 lg:mt-6">
                {#each data.lists as list }
                    <ShoppingList
                        description={list.description}
                        id="{list.id}"
                        hideDeleteButton />
                {/each}
            </ul>
        </section>
    </main>
{/if}

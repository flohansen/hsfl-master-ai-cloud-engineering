<script lang="ts">
    import { onMount } from "svelte";
    import ShoppingListItem from "$lib/ShoppingListItem.svelte";
    import { page } from "$app/stores";
    import Header from "$lib/Header.svelte";

    let jsonData: any = null;

    onMount(async () => {
        try {
            const response = await fetch("http://localhost:8080/api/v1/shoppinglist/2");
            if (response.ok) {
                jsonData = await response.json();
            } else {
                throw new Error("Failed to fetch data");
            }
        } catch (error) {
            console.error(error);
        }
    });
</script>

<Header headline="{$page.data.headline}"/>

<main class="mt-8 sm:ml-20 md:ml-24 lg:max-w-4xl lg:mx-auto lg:mt-10 xl:max-w-5xl">
    <h2 class="px-5 text-gray-dark text-sm font-medium mt-6 lg:mt-10 lg:text-base">
        Offene Einkaufslisten
    </h2>
    <ul class="px-5 mt-4 grid grid-cols-1 gap-y-4 lg:gap-y-6 lg:mt-6">
        {#if jsonData}
            {#each jsonData.items as item}
                <ShoppingListItem description={item.description}/>
            {/each}
        {:else}
            <p>Loading data...</p>
        {/if}
    </ul>
</main>
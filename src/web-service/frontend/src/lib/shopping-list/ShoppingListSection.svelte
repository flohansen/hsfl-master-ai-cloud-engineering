<script lang="ts">
    import ShoppingList from "$lib/shopping-list/ShoppingList.svelte";

    export let lists: {id: number, description: string, complete?: boolean }[];
    export let label: string;
</script>

{#if lists.length !== 0}
    <section class="mb-8">
        <h2 class="px-5 text-gray-dark text-sm font-medium mt-6 lg:mt-10 lg:text-base">
            {label}
        </h2>
        <ul class="px-5 mt-4 grid grid-cols-1 gap-y-4 lg:gap-y-6 lg:mt-6">
            {#each lists as list}
                <ShoppingList
                    description={list.description}
                    id="{list.id}" />
            {/each}
        </ul>
    </section>
{/if}
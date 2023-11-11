<script lang="ts">
    import Chevron from "../../assets/svg/Chevron.svelte";
    import {onMount} from "svelte";
    import {handleErrors} from "../../assets/helper/handleErrors";

    let isOpen: boolean = false;
    let label: string = 'Eintrag auswählen';
    let products: { id: number, description: string }[] = [];
    export let entryId: number;

    onMount(async () => {
        const apiUrlProducts: string = '/api/v1/product';

        fetch(apiUrlProducts)
            .then(handleErrors)
            .then(data => products = data)
            .catch(error => console.error("Failed to fetch data:", error.message));
    });
</script>

<div class="relative my-5 lg:my-8">
    <p class="text-gray-dark text-sm font-medium mb-2 lg:mb-3 lg:text-base">
        Anzahl auswählen:
    </p>

    <button
        aria-haspopup="listbox"
        aria-expanded="{isOpen}"
        aria-controls="select-options"
        on:click={() => isOpen = ! isOpen}
        class="rounded-lg border px-3 py-2 w-full text-left text-green-dark/75 flex items-center justify-between transition-all duration-300 ease-in-out hover:bg-blue-light/25 lg:px-4 lg:py-3
            {isOpen ? 'border-green-dark' : 'border-green-dark/50'}">
        <span class="font-medium text-sm lg:text-base">{label}</span>
        <Chevron classes="w-4 h-4 transition-all duration-300 ease-in-out {isOpen ? 'rotate-180' : ''}"/>
    </button>

    <div class:hidden={! isOpen}>
        <ul
            id="select-options"
            role="listbox"
            class="absolute h-full min-h-[20vh] overflow-y-auto top-[4.4rem] w-full bg-gray-light rounded-lg shadow-gray-dark/30 shadow-lg lg:top-[5.75rem]">
            {#each products as product}
                <li role="option" aria-selected="false" class="px-4 group transition-all ease-in-out duration-300 hover:bg-gray-dark/25 lg:px-6">
                    <button
                        on:click={() => {entryId = product.id; label = product.description; isOpen = false}}
                        class="w-full text-left text-sm py-2.5 sm:py-3 border-t border-t-gray-dark/20 group-first:border-none lg:py-4 lg:text-base">
                        {product.description}
                    </button>
                </li>
            {/each}
        </ul>
    </div>
</div>
<script lang="ts">
    import Chevron from "../../assets/svg/Chevron.svelte";
    import {onMount} from "svelte";
    import {clickOutside} from "../../assets/helper/clickOutside";

    let isOpen: boolean = false;
    let products: { id: number, description: string }[] = [];
    let label: string = '';

    export let selectedProduct: number;

    const apiUrlProducts: string = '/api/v1/product';

    onMount(async () => {
        try {
            const response = await fetch(apiUrlProducts);
            response.ok
                ? products = await response.json()
                : console.error('Failed to fetch data');
        } catch (error) {
            console.error(error);
        }
    });

    function handleClickOutside() {
        if (! isOpen) return;
        isOpen = false;
    }

    function toggleOpen() {
        isOpen = ! isOpen;
    }

</script>

<div class="relative">
    <button
        aria-haspopup="listbox"
        aria-expanded="{isOpen}"
        aria-controls="select-options"
        on:click={toggleOpen}
        class="rounded-lg mt-4 border px-3 py-2 w-full text-left text-green-dark/75 flex items-center justify-between transition-all duration-300 ease-in-out hover:bg-blue-light/25
            {isOpen ? 'border-green-dark' : 'border-green-dark/50'}">
        <span class="font-medium text-sm">
            {#if label}
                {label}
            {:else}
                Eintrag auswählen
            {/if}
        </span>
        <Chevron classes="w-4 h-4 transition-all duration-300 ease-in-out  {isOpen ? 'rotate-180' : ''}"/>
    </button>

    <div class:hidden={!isOpen} use:clickOutside on:click_outside={handleClickOutside}>
        <ul
            id="select-options"
            role="listbox"
            class="absolute h-full min-h-[6rem] overflow-y-scroll top-10 w-full bg-white rounded-lg shadow-gray-dark/20 shadow-lg">
            {#each products as product}
                <li role="option" aria-selected="false">
                    <button
                        on:click={() => {selectedProduct = product.id; label = product.description; isOpen = false}}
                        class="w-full text-left text-sm px-3 py-2 transition-all ease-in-out duration-300 hover:text-gray-dark hover:bg-gray-light sm:py-3">
                        {product.description}
                    </button>
                </li>
            {/each}
        </ul>
    </div>
</div>
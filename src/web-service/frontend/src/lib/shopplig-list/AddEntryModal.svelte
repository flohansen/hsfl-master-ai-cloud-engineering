<script lang="ts">
    import Add from "../../assets/svg/Add.svelte";
    import {clickOutside} from "../../assets/helper/clickOutside";
    import Select from "$lib/general/Select.svelte";

    let isOpen = false;
    export let selectedProduct: number;

    function toggleModal() {
        isOpen = !isOpen;
    }

    function handleClickOutside() {
        if (! isOpen) return;
        isOpen = false;
    }
</script>

<button on:click={toggleModal} class="-ml-[2px] mt-6 text-green-dark hover:text-green-light flex items-center justify-center gap-x-4">
    <Add classes="w-6 h-6 transition-all ease-in-out duration-300"/>
    <span class="block transition-all ease-in-out duration-300 text-sm lg:text-base">
        Einträge hinzufügen
    </span>
</button>

<div class:hidden={! isOpen} class="bg-black/80 fixed inset-0 w-screen h-screen"></div>

<div
    use:clickOutside
    on:click_outside={handleClickOutside}
    class:hidden={! isOpen}
    class="fixed inset-x-4 h-[40vh] top-1/2 -translate-y-1/2 bg-gray-light rounded-xl p-4 sm:left-28 sm:right-8 lg:inset-x-0 lg:max-w-[52rem] lg:mx-auto xl:max-w-[62rem]">
    <h3 class="text-lg font-semibold mt-4">
        Eintrag hinzufügen:
    </h3>

    <Select bind:selectedProduct={selectedProduct}/>
</div>
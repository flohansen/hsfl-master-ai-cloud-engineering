<script lang="ts">
    import Checkmark from "../../assets/svg/Checkmark.svelte";
    import {clickOutside} from "../../assets/helper/clickOutside";
    import {createEventDispatcher} from 'svelte'

    const dispatch = createEventDispatcher()

    export let isOpen: boolean;

    function handleClickOutside(): void {
        if (! isOpen) return;
        isOpen = false;
    }

</script>

<div class:hidden={! isOpen} class="bg-black/80 fixed inset-0 w-screen h-screen"></div>

<section
    use:clickOutside
    on:click_outside={handleClickOutside}
    class:hidden={! isOpen}
    class="fixed inset-x-4 h-min top-1/2 -translate-y-1/2 bg-white rounded-xl px-4 py-6 sm:left-28 sm:right-8 lg:inset-x-0 lg:max-w-[52rem] lg:mx-auto lg:px-6 lg:px-10 xl:max-w-[62rem]">

    <slot></slot>

    <button
        on:click={() => {dispatch('submit') }}
        class="bg-green-light mt-8 mx-auto text-white rounded-xl px-5 py-2 flex items-center justify-center gap-x-2 transition-all ease-in-out duration-300 hover:bg-green-dark disabled:bg-gray-light disabled:text-gray-dark">
        <span class="text-sm lg:text-base">Speichern</span>
        <Checkmark classes="w-5 h-5 stroke-2"/>
    </button>
</section>
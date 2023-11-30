<script lang="ts">
    import {clickOutside} from "../../assets/helper/clickOutside";
    import {createEventDispatcher} from 'svelte'
    import SubmitButton from "$lib/forms/SubmitButton.svelte";
    import Close from "../../assets/svg/Close.svelte";

    const dispatch = createEventDispatcher()

    export let isOpen: boolean;
    export let submitLabel: string = 'Speichern';

    function handleClickOutside(): void {
        if (! isOpen) return;
        isOpen = false;
    }

</script>

<div class:hidden={! isOpen} class="bg-black/80 fixed inset-0 w-screen h-screen z-10"></div>

<section
    use:clickOutside
    on:click_outside={handleClickOutside}
    class:hidden={! isOpen}
    class="fixed inset-x-4 z-20 h-min top-1/2 -translate-y-1/2 bg-white rounded-xl px-4 py-6 sm:left-28 sm:right-8 lg:inset-x-0 lg:max-w-[52rem] lg:mx-auto lg:px-6 lg:px-10 xl:max-w-[62rem]">

    <slot></slot>

    <div class="md:flex md:items-center md:gap-x-4 md:justify-center">
        <button
            on:click={() => isOpen = !isOpen}
            class="border-[1.5px] border-green-light mt-8 text-green-light rounded-xl px-5 py-2 flex items-center justify-center gap-x-2 transition-all ease-in-out duration-300 hover:bg-green-dark hover:border-green-dark hover:text-white">
            <span class="text-sm lg:text-base">Abbrechen</span>
            <Close classes="w-5 h-5"/>
        </button>

        <SubmitButton
            label={submitLabel}
            centered={false}
            on:submit={() => dispatch('submit')}/>
    </div>
</section>
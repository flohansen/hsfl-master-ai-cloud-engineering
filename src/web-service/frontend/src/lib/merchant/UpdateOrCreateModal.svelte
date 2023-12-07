<script lang="ts">
    import {clickOutside} from "../../assets/helper/clickOutside";

    export let isOpen: boolean;

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

    <h3 class="text-lg text-center font-semibold lg:text-xl">
        Was genau möchtest du neu erstellen oder editieren?
    </h3>

    <div class="flex flex-col items-center mt-6 gap-y-2 gap-x-6 md:flex-row md:justify-center">
        <a
            href="/products/add"
            class="border-[1.5px] max-w-max border-green-light text-green-light rounded-xl px-5 py-2 flex items-center justify-center gap-x-2 transition-all ease-in-out duration-300 hover:bg-green-dark hover:border-green-dark hover:text-white">
            <span class="text-sm lg:text-base">Produkt erstellen oder editieren</span>
        </a>
        <a
            href="/prices/add"
            class="border-[1.5px] max-w-max border-green-light text-green-light rounded-xl px-5 py-2 flex items-center justify-center gap-x-2 transition-all ease-in-out duration-300 hover:bg-green-dark hover:border-green-dark hover:text-white">
            <span class="text-sm lg:text-base">Preis erstellen oder editieren</span>
        </a>
    </div>
</section>
<script lang="ts">
    import Close from "../../assets/svg/Close.svelte";

    export let url: string;
    export let label: string;
</script>

<a
    href="{url}"
    aria-label="{label}"
    class="group rounded-full border-[1.5px] border-green-dark w-8 h-8 flex items-center justify-center transition-all ease-in-out duration-300 cursor-pointer hover:bg-green-dark">
    <span class="text-green-dark group-hover:text-white">
        <Close classes="w-5 h-5"/>
    </span>
</a>
<script lang="ts">
    import Checkmark from "../../assets/svg/Checkmark.svelte";
    import {createEventDispatcher} from "svelte";

    export let label: string = 'Speichern';
    export let centered: boolean = true;
    export let type: "button" | "submit" | "reset" | null | undefined = 'button';
    const dispatch = createEventDispatcher();
</script>

<button
    type={type}
    on:click={() => dispatch('submit')}
    class="bg-green-light text-white rounded-xl px-5 py-2 flex items-center justify-center gap-x-2 transition-all ease-in-out duration-300 hover:bg-green-dark disabled:bg-gray-light disabled:text-gray-dark { centered ? 'mx-auto mt-8' : '' }">
    <span class="text-sm lg:text-base">{label}</span>
    <Checkmark classes="w-5 h-5 stroke-2"/>
</button>
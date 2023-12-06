<script lang="ts">
    import Checkmark from "../../assets/svg/Checkmark.svelte";
    import {createEventDispatcher} from "svelte";

    export let label: string;
    export let id: number;
    export let checked: boolean;

    const dispatch = createEventDispatcher();
</script>

<div class="flex items-center mr-4 mb-2 relative">
    <input
        on:click={() => { checked = ! checked; dispatch('updateShoppingListEntry') } }
        type="checkbox"
        id="input-{id}"
        name="Checkbox"
        checked="{checked}"
        class="z-10 cursor-pointer opacity-0 absolute h-5 w-5 peer"/>

    <div class="relative cursor-pointer z-0 border rounded-md border-green-dark/50 w-5 h-5 mr-4 peer-checked:border-none peer-checked:bg-green-light/25"></div>
    <Checkmark classes="z-0 w-4 h-4 cursor-pointer stroke-2 text-green-dark hidden absolute left-0.5 peer-checked:block"/>

    <label
        for="input-{id}"
        class="text-sm font-medium cursor-pointer lg:text-base { checked ? 'line-through opacity-50' : '' }">
        {label}
    </label>
</div>
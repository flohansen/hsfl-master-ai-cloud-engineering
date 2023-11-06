<script lang="ts">
    import Checkmark from "../../assets/svg/Checkmark.svelte";

    export let label: string;
    export let id: number;
</script>

<div class="flex items-center mr-4 mb-2 relative cursor-pointer">
    <input
        type="checkbox"
        id="input-{id}"
        name="Checkbox"
        class="opacity-0 absolute h-5 w-5 peer"/>

    <div class="border rounded-md border-green-dark/50 w-5 h-5 mr-4 peer-checked:border-none peer-checked:bg-green-light/25"></div>
    <Checkmark classes="w-4 h-4 stroke-2 text-green-dark hidden absolute left-0.5 peer-checked:block"/>

    <label
        for="input-{id}"
        class="text-sm font-medium lg:text-base">
        {label}
    </label>
</div>
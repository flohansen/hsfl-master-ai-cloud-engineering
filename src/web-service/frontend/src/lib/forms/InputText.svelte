<script lang="ts">
    export let value: string | number;
    export let label: string;
    export let fieldName: string;
    export let type: string = 'text';
</script>

<div class="w-full">
    <label
        for="{fieldName}"
        class="text-sm text-gray-dark font-medium block mb-2">
        {label}:*
    </label>
    {#if type === 'number'}
        <input
            id="{fieldName}"
            name="{fieldName}"
            type="number"
            required
            placeholder="{label}"
            bind:value={value}
            class="text-sm rounded-lg border px-3 py-2 w-full text-left text-green-dark/75 flex items-center justify-between transition-all duration-300 ease-in-out hover:bg-blue-light/25 lg:px-4 lg:py-3"/>
    {:else}
        <input
            id="{fieldName}"
            name="{fieldName}"
            type="text"
            required
            placeholder="{label}"
            bind:value={value}
            class="text-sm rounded-lg border px-3 py-2 w-full text-left text-green-dark/75 flex items-center justify-between transition-all duration-300 ease-in-out hover:bg-blue-light/25 lg:px-4 lg:py-3"/>
    {/if}
</div>
<script lang="ts">
    export let value: string | number | null;
    export let label: string | null;
    export let fieldName: string;
    export let type: string = 'text';
    export let readonly: boolean = false;

    function typeAction(node: any) {
        node.type = type;
    }
</script>

<div class="w-full">
    <label
        for="{fieldName}"
        class="text-sm text-gray-dark font-medium block mb-2">
        {label}: *
    </label>
    <input
        id="{fieldName}"
        name="{fieldName}"
        use:typeAction
        required
        readonly="{readonly}"
        placeholder="{label}"
        bind:value={value}
        class="text-sm rounded-lg border px-3 py-2 w-full text-left text-green-dark/75 flex items-center justify-between transition-all duration-300 ease-in-out hover:bg-blue-light/25 lg:px-4 lg:py-3 read-only:bg-gray-light"/>
</div>
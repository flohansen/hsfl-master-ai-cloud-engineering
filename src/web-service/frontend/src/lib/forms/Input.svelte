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
        class="text-sm rounded-lg border px-4 py-3 w-full text-left text-green-dark/75 flex items-center justify-between hover:bg-blue-light/25 read-only:bg-gray-light"/>
</div>
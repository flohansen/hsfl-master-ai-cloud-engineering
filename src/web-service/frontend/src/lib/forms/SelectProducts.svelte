<script lang="ts">
    import Chevron from "../../assets/svg/Chevron.svelte";
    import {onMount} from "svelte";
    import {handleErrors} from "../../assets/helper/handleErrors";
    import Select from 'svelte-select';

    let placeholder: string = 'Eintrag auswählen';
    let listOpen: boolean = false;
    const itemId = 'id';
    const label = 'description';
    const id = 'product';
    let items: { id: number, description: string }[] = [];

    export let justValue: number;

    onMount(async () => {
        const apiUrlProducts: string = '/api/v1/product';

        fetch(apiUrlProducts)
            .then(handleErrors)
            .then(data => items = data)
            .catch(error => console.error("Failed to fetch data:", error.message));
    });
</script>

<div class="relative my-5 lg:my-8">
    <label for="{id}" class="text-sm text-gray-dark font-medium block mb-2">
        {label}: *
    </label>

    <Select
        {id}
        {itemId}
        {label}
        {items}
        {placeholder}
        bind:listOpen
        bind:justValue
        clearable={false}
        listOffset={0}
        showChevron
        required
        --border="1px solid rgba(49, 112, 80, 0.75)"
        --border-radius="0.5rem"
        --border-focused="1px solid rgba(49, 112, 80, 1)"
        --placeholder-color="rgba(49, 112, 80, 0.75)"
        --input-color="rgba(49, 112, 80, 1)"
        --selected-item-color="rgba(49, 112, 80, 1)"
        --padding="0.85rem 1rem 0.85rem 1rem"
        --font-size="0.875rem"
        --list-border-radius="0 0 0.5rem 0.5rem"
        --item-hover-bg="rgba(143, 143, 143, 0.25)"
        --item-is-active-bg="rgba(49, 112, 80, 0.75)"
        --list-background="#F4F4F9"
        --chevron-width="1.75rem"
        --chevron-height="1rem">
        <Chevron slot="chevron-icon" classes="w-4 h-4 text-green-dark/75 transition-all ease-in-out duration-300 { listOpen ? 'rotate-180' : '' }"/>
    </Select>
</div>
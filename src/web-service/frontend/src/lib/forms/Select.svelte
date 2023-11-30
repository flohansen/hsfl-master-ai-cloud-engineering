<script lang="ts">
    import Chevron from "../../assets/svg/Chevron.svelte";
    import {onMount} from "svelte";
    import {handleErrors} from "../../assets/helper/handleErrors";
    import Select from 'svelte-select';

    let placeholder: string = 'Eintrag auswählen';
    const itemId = 'id';
    const label = 'description';

    let items: { id: number, description: string }[] = [];
    export let entryId: number;

    onMount(async () => {
        const apiUrlProducts: string = '/api/v1/product';

        fetch(apiUrlProducts)
            .then(handleErrors)
            .then(data => items = data)
            .catch(error => console.error("Failed to fetch data:", error.message));
    });
</script>

<div class="relative my-5 lg:my-8">
    <p class="text-gray-dark text-sm font-medium mb-2 lg:mb-3 lg:text-base">
        Anzahl auswählen:
    </p>
    <Select
        {itemId}
        {label}
        {items}
        {placeholder}
        clearable={false}
        listOffset={0}
        showChevron
        --border="1px solid rgba(49, 112, 80, 0.5)"
        --border-radius="0.5rem"
        --border-focused="1px solid rgba(49, 112, 80, 1)"
        --placeholder-color="rgba(49, 112, 80, 0.75)"
        --input-color="rgba(49, 112, 80, 1)"
        --selected-item-color="rgba(49, 112, 80, 1)"
        --font-size="0.875rem"
        --list-border-radius="0 0 0.5rem 0.5rem"
        --item-hover-bg="rgba(143, 143, 143, 0.25)"
        --item-is-active-bg="rgba(49, 112, 80, 0.75)"
        --list-background="#F4F4F9"
        --chevron-width="1.75rem"
        --chevron-height="1rem">
        <Chevron slot="chevron-icon" classes="w-4 h-4 text-green-dark/75 mr-3"/>
    </Select>
</div>
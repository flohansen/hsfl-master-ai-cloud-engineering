<script lang="ts">
    import Close from "../../assets/svg/Close.svelte";
    import Checkmark from "../../assets/svg/Checkmark.svelte";

    interface Product {
        id: number,
        description: string,
        ean: number,
    }

    export let productData: Product;
    export let successfulMessage: string;
    export let notSuccessfulMessage: string;
</script>

<div class="grid grid-cols-[1.5rem,auto] gap-x-2 mt-3">
    <figure class="w-6 h-6 rounded-full flex items-center justify-center { productData ? 'bg-green-light/25' : 'bg-red/25' }">
        {#if productData}
            <Checkmark classes="text-green-dark w-4 h-4"/>
        {:else}
            <Close classes="text-red w-4 h-4"/>
        {/if }
    </figure>

    <div class="text-sm text-gray-dark">
        <p>
            {#if productData}
                {successfulMessage}
            {:else}
                {notSuccessfulMessage}
            {/if}
        </p>

        {#if productData}
            <p class="mt-2">
                Produkt-ID: {productData.id}<br>
                Produkt-Beschreibung: {productData.description}<br>
                Produkt-EAN: {productData.ean}
            </p>
        {/if}
    </div>
</div>
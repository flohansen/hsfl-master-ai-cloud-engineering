<script lang="ts">
    import InputText from "$lib/forms/InputText.svelte";
    import {handleErrors} from "../../assets/helper/handleErrors";
    import Checkmark from "../../assets/svg/Checkmark.svelte";
    import Close from "../../assets/svg/Close.svelte";

    let eanSubmitted: boolean = false;
    export let productId: number;
    export let productEan: number;

    function findProduct(): void {
        if (! productEan) return;

        eanSubmitted = true;
        const apiUrl: string = `/api/v1/product/ean/${productEan}`
        fetch(apiUrl)
            .then(handleErrors)
            .then(data => productId = data.id)
            .catch(error => console.error("Failed to fetch data:", error.message));
    }
</script>

<section class="mb-6 lg:mb-8">
    <InputText
        fieldName="productEan"
        label="EAN des Produktes"
        type="number"
        readonly={eanSubmitted}
        bind:value={productEan} />

    {#if ! eanSubmitted}
        <button
            on:click={findProduct}
            class="mt-10 bg-green-light mt-8 mx-auto text-white rounded-xl px-5 py-2 flex items-center justify-center gap-x-2 transition-all ease-in-out duration-300 hover:bg-green-dark disabled:bg-gray-light disabled:text-gray-dark">
            <span class="text-sm lg:text-base">Produkt suchen</span>
        </button>
    {:else}
        <div class="flex items-start gap-x-2 mt-3">
            <figure class="w-6 h-6 rounded-full flex items-center justify-center { productId ? 'bg-green-light/25' : 'bg-red/25' }">
                {#if productId}
                    <Checkmark classes="text-green-dark w-4 h-4"/>
                {:else}
                    <Close classes="text-red w-4 h-4"/>
                {/if }
            </figure>
            <span class="text-sm text-gray-dark">
                {#if productId}
                    Produkt gefunden
                {:else}
                    Produkt mit angegebener EAN konnte nicht gefunden werden.<br>
                    Das Produkt muss neu hinzugefügt werden.
                {/if}
            </span>
        </div>
    {/if}
</section>
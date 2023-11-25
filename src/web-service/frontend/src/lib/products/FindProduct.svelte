<script lang="ts">
    import InputText from "$lib/forms/InputText.svelte";
    import {handleErrors} from "../../assets/helper/handleErrors";
    import Checkmark from "../../assets/svg/Checkmark.svelte";
    import Close from "../../assets/svg/Close.svelte";
    import FetchFeedback from "$lib/products/FetchFeedback.svelte";

    interface Product {
        id: number,
        description: string,
        ean: number,
    }

    interface Price {
        price: number,
    }

    export let eanSubmitted: boolean = false;
    export let productEan: number;
    export let productData: Product;
    export let shouldFindPrice: boolean = false;
    export let priceIsAlreadyCreated: boolean = false;
    let priceData: Price;

    function findProduct(): void {
        if (! productEan) return;

        eanSubmitted = true;
        const apiUrl: string = `/api/v1/product/ean/${productEan}`

        fetch(apiUrl)
            .then(handleErrors)
            .then(data => { productData = data; findPrice()})
            .catch(error => console.error("Failed to fetch data:", error.message));
    }

    function findPrice() {
        if (! shouldFindPrice || ! productData.id) return;

        const userId: number = 2 // TODO: add current user id
        const apiUrl: string = `/api/v1/price/${productData.id}/${userId}`

        fetch(apiUrl)
            .then(handleErrors)
            .then(data => {
                if (! data) return;
                priceData = data;
                priceIsAlreadyCreated = true;
            })
            .catch(error => console.error("Failed to fetch data:", error.message));
    }
</script>

<section class="mb-6 lg:mb-8">
    {#if ! productData}
        <InputText
            fieldName="productEan"
            label="EAN des Produktes"
            type="number"
            readonly={eanSubmitted}
            bind:value={productEan} />
    {/if}

    {#if ! eanSubmitted}
        <button
            on:click={findProduct}
            class="mt-10 bg-green-light mt-8 mx-auto text-white rounded-xl px-5 py-2 flex items-center justify-center gap-x-2 transition-all ease-in-out duration-300 hover:bg-green-dark disabled:bg-gray-light disabled:text-gray-dark">
            <span class="text-sm lg:text-base">Produkt suchen</span>
        </button>
    {:else}
        <FetchFeedback
            productData={productData}
            priceData={priceData}
            successfulMessage="Produkt mit angegebener EAN wurde gefunden"
            notSuccessfulMessage="Produkt mit angegebener EAN konnte nicht gefunden werden. Das Produkt muss neu hinzugefügt werden." />
    {/if}
</section>
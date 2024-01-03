<script lang="ts">
    import Input from "$lib/forms/Input.svelte";
    import FetchFeedback from "$lib/products/FetchFeedback.svelte";
    import {handleErrors} from "../../assets/helper/handleErrors";
    import {validateEan} from "../../assets/helper/validateEan";

    interface FeedbackOption {
        type: FeedbackType,
        label: string,
    }

    const feedbackOptions: FeedbackOption[] = [
        { type: 'successful', label: 'Produkt mit angegebener EAN wurde gefunden.' },
        { type: 'unsuccessful', label: 'Angegebene EAN ist nicht EAN-8 oder EAN-13 valide.' },
        { type: 'notFound', label: 'Produkt mit angegebener EAN konnte nicht gefunden werden.' },
    ]

    type FeedbackType = "successful" | "unsuccessful" | "notFound";
    let currentFeedbackOption: FeedbackOption;

    export let eanSubmitted: boolean = false;
    export let productEan: number;

    export let shouldFindPrice: boolean = false;
    export let priceIsAlreadyCreated: boolean = false;

    export let productData: { id: number, description: string, ean: number };
    export let priceData: { price: number };

    $: showFeedback = !!currentFeedbackOption;

    function getOptionByType(typeToFind: FeedbackType): FeedbackOption {
        return feedbackOptions.find(option => option.type === typeToFind) ?? feedbackOptions[0];
    }

    function findProduct(): void {
        if (! productEan || ! validateEan(productEan)) {
            currentFeedbackOption = getOptionByType('unsuccessful');
            return;
        }

        eanSubmitted = true;
        const apiUrl: string = `/api/v1/product/ean/${productEan}`

        fetch(apiUrl)
            .then(handleErrors)
            .then(data => {
                productData = data;
                findPrice();
                currentFeedbackOption = productData
                    ? getOptionByType('successful')
                    : getOptionByType('notFound');
            })
            .catch(error => console.error("Failed to fetch data:", error.message));
    }

    function findPrice() {
        if (! shouldFindPrice || ! productData) return;

        const userId: number = 2 // TODO: add current user id
        const apiUrl: string = `/api/v1/price/${productData.id}/${userId}`

        fetch(apiUrl)
            .then(handleErrors)
            .then(data => {
                if (! data) return;
                priceData = data;
                priceIsAlreadyCreated = true;
            })
            .catch(error => {console.error("Failed to fetch data:", error.message)});
    }
</script>

<section class="mb-6 lg:mb-8">
    {#if ! productData}
        <Input
            fieldName="productEan"
            label="EAN des Produktes"
            type="number"
            readonly={eanSubmitted}
            bind:value={productEan} />
    {/if}

    {#if showFeedback}
        <FetchFeedback
            feedback={currentFeedbackOption.type}
            label={currentFeedbackOption.label}>
            {#if productData}
                <p class="mt-2">
                    Produkt-ID: {productData.id}<br>
                    Produkt-Beschreibung: {productData.description}<br>
                    Produkt-EAN: {productData.ean}
                    {#if priceData}
                        <br>Produkt-Preis: {priceData.price} €
                    {/if}
                </p>
            {/if}
        </FetchFeedback>
    {/if}

    {#if ! eanSubmitted}
        <button
            on:click={findProduct}
            class="mt-10 bg-green-light mx-auto text-white rounded-xl px-5 py-2 flex items-center justify-center gap-x-2 transition-all ease-in-out duration-300 hover:bg-green-dark disabled:bg-gray-light disabled:text-gray-dark">
            <span class="text-sm lg:text-base">Produkt suchen</span>
        </button>
    {/if}
</section>
<script lang="ts">
    import SubmitButton from "$lib/forms/SubmitButton.svelte";
    import InputText from "$lib/forms/InputText.svelte";
    import {handleErrors} from "../../assets/helper/handleErrors";
    import FetchFeedback from "$lib/products/FetchFeedback.svelte";

    interface Product {
        id: number,
        description: string,
        ean: number,
    }

    let productDescription: number;
    export let productData: Product;
    export let productSubmitted: boolean = false;
    export let productEan: number;

    function submit(): void {
        const apiUrl: string = `/api/v1/product/`
        const requestOptions = {
            method: "POST",
            headers: { 'Content-Type': 'application/json' },
            body: `{"description": "${productDescription}", "ean": ${productEan}}`,
        };

        fetch(apiUrl, requestOptions)
            .then(handleErrors)
            .then(data => { productData = data; productSubmitted = true })
            .catch(error => console.error("Failed to fetch data:", error.message));
    }
</script>

<section>
    {#if ! productSubmitted}
        <InputText
            fieldName="productName"
            label="Name des Produktes"
            bind:value={productDescription} />

        <div class="mt-10">
            <SubmitButton on:submit={submit}/>
        </div>
    {:else}
        <FetchFeedback
            productData={productData}
            successfulMessage="Neues Produkt wurde erfolgreich hinzugefügt."
            notSuccessfulMessage="Es ist leider ein Fehler aufgetreten. Überprüfe bitte deine Angaben" />
    {/if}
</section>
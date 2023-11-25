<script lang="ts">
    import SubmitButton from "$lib/forms/SubmitButton.svelte";
    import InputText from "$lib/forms/InputText.svelte";
    import {handleErrors} from "../../assets/helper/handleErrors";

    let productPrice: number;
    export let formSubmitted: boolean = false;
    export let productId: number;

    function submit(): void {
        if (! productId || ! productPrice) return;

        const userId: number = 2; // TODO: add current user id
        const apiUrl: string = `/api/v1/price/${productId}/${userId}`
        const requestOptions = {
            method: "POST",
            headers: { 'Content-Type': 'application/json' },
            body: `{"price": ${productPrice}}`,
        };

        fetch(apiUrl, requestOptions)
            .then(handleErrors)
            .then(() => formSubmitted = true)
            .catch(error => console.error("Failed to fetch data:", error.message));
    }
</script>

<section class="mt-6 lg:mt-8">
    <InputText
        fieldName="productPrice"
        label="Preis des Produktes"
        type="number"
        bind:value={productPrice} />

    <div class="mt-10">
        <SubmitButton on:submit={submit}/>
    </div>
</section>
<script lang="ts">
    import {page} from "$app/stores";
    import SubmitButton from "$lib/forms/SubmitButton.svelte";
    import {handleErrors} from "../../../assets/helper/handleErrors";
    import Badge from "$lib/general/Badge.svelte";
    import InputText from "$lib/forms/InputText.svelte";
    import CloseButton from "$lib/general/CloseButton.svelte";
    import BackLink from "$lib/general/BackLink.svelte";
    import Checkmark from "../../../assets/svg/Checkmark.svelte";
    import FindProduct from "$lib/products/FindProduct.svelte";

    interface Product {
        id: number,
        description: string,
        ean: number,
    }

    let productPrice: number;
    let formSubmitted: boolean = false;

    let productData: Product = { id: 0, description: '', ean: 0 };

    function submit(): void {
        if (! productPrice || ! productData.ean) return;
        createProduct();
    }

    function createProduct(): void {
        if (! productData.id) return;

        const apiUrl: string = `/api/v1/product/`
        const requestOptions = {
            method: "POST",
            headers: { 'Content-Type': 'application/json' },
            body: `{"description": "${productData.description}", "ean": ${productData.ean}}`,
        };

        fetchData(apiUrl, requestOptions);
    }

    function fetchData(apiUrl: string, requestOptions: object): void {
        fetch(apiUrl, requestOptions)
            .then(handleErrors)
            .then(()=> formSubmitted = true)
            .catch(error => console.error("Failed to fetch data:", error.message));
    }
</script>

<header>
    {#if ! formSubmitted}
        <h1 class="font-bold text-xl md:text-2xl xl:text-3xl">
            {$page.data.metaTitle}
        </h1>
        <CloseButton
            url="/profile"
            label="Erstellen eines Produktes abbrechen" />
    {:else}
        <BackLink
            url="/"
            label="Zur Startseite" />
    {/if}
</header>

<main>
    <div class="mx-5 bg-white rounded-xl p-4 lg:p-6">
        {#if ! formSubmitted}
            <FindProduct
                bind:productId={productData.id}
                bind:productEan={productData.ean}/>

            {#if productData.ean}
                <section>
                    <div class="flex flex-col gap-y-6 lg:gap-y-8">
                        {#if ! productData.id}
                            <InputText
                                fieldName="productName"
                                label="Name des Produktes"
                                bind:value={productData.description} />
                        {/if}
                        <InputText
                            fieldName="productPrice"
                            label="Preis des Produktes"
                            type="number"
                            bind:value={productPrice} />
                    </div>

                    <div class="mt-10">
                        <SubmitButton on:submit={submit} />
                    </div>
                </section>
            {/if}
        {:else}
            <Badge />
            <h2 class="font-semibold text-lg lg:text-xl">
                Dein Produkt und dein Preis wurden erfolgreich gespeichert.
            </h2>
        {/if}
    </div>
</main>
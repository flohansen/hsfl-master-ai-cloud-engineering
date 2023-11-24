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

    let productPrice: number;
    let productId: number;
    let productDescription: number;
    let productEan: number;

    let formSubmitted: boolean = false;
    let eanSubmitted: boolean = false;

    function submit(): void {
        if (! productPrice || ! productEan) return;
        createProduct();
    }

    function createProduct(): void {
        const apiUrl: string = `/api/v1/product/`
        const requestOptions = {
            method: "POST",
            headers: { 'Content-Type': 'application/json' },
            body: `{"description": "${productDescription}", "ean": ${productEan}}`,
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
                bind:eanSubmitted={eanSubmitted}
                bind:productId={productId}
                bind:productEan={productEan}/>

            {#if eanSubmitted}
                <section>
                    <div class="flex flex-col gap-y-6 lg:gap-y-8">
                        {#if ! productId}
                            <InputText
                                fieldName="productName"
                                label="Name des Produktes"
                                bind:value={productDescription}/>
                        {/if}
                        <InputText
                            fieldName="productPrice"
                            label="Preis des Produktes"
                            type="number"
                            bind:value={productPrice} />
                    </div>

                    <div class="mt-10">
                        <SubmitButton on:submit={submit}/>
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
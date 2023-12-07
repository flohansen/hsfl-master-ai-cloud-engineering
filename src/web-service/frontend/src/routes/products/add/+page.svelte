<script lang="ts">
    import {page} from "$app/stores";
    import Badge from "$lib/general/Badge.svelte";
    import CloseButton from "$lib/general/CloseButton.svelte";
    import BackLink from "$lib/general/BackLink.svelte";
    import FindProduct from "$lib/products/FindProduct.svelte";
    import SubmitButton from "$lib/forms/SubmitButton.svelte";
    import InputText from "$lib/forms/InputText.svelte";
    import {handleErrors} from "../../../assets/helper/handleErrors";

    interface Product {
        id: number,
        description: string,
        ean: number,
    }

    let productEan: number;
    let productData: Product;
    let productDescription: string;

    let formSubmitted: boolean = false;
    let eanSubmitted: boolean = false;

    $: productDescription = productData ? productData.description : '';

    function submit(): void {
        if (! productDescription) return;
        productData ? updateProduct() : createProduct();
    }

    function updateProduct(): void {
        if (! productData.id) return;
        const apiUrl: string = `/api/v1/product/${productData.id}`
        fetchContent(apiUrl, "PUT");
    }

    function createProduct(): void {
        const apiUrl: string = `/api/v1/product/`
        fetchContent(apiUrl, "POST");
    }

    function fetchContent(apiUrl: string, method: string): void {
        const requestOptions = {
            method: method,
            headers: { 'Content-Type': 'application/json' },
            body: `{"description": "${productDescription}", "ean": ${productEan}}`,
        };

        fetch(apiUrl, requestOptions)
            .then(handleErrors)
            .then(() => formSubmitted = true)
            .catch(error => console.error("Failed to fetch data:", error.message));
    }
</script>

<header>
    {#if ! formSubmitted}
        <h1 class="font-bold text-xl w-[90%] md:text-2xl xl:text-3xl">
            {$page.data.headline}
        </h1>
        <CloseButton
            url="/merchants"
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
                bind:productData={productData}
                bind:productEan={productEan} />

            {#if eanSubmitted}
                <section>
                    <InputText
                        fieldName="productName"
                        label="Name des Produktes"
                        bind:value={productDescription} />

                    <div class="mt-10">
                        <SubmitButton on:submit={submit}/>
                    </div>
                </section>
            {/if}
        {:else}
            <Badge />
            <h2 class="font-semibold mb-6 text-lg lg:text-xl lg:mb-8">
                Dein Produkt wurde erfolgreich gespeichert.
            </h2>
            <BackLink
                url="/prices/add"
                label="Preis erstellen"
                reverse />
        {/if}
    </div>
</main>
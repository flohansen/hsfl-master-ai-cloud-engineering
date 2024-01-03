<script lang="ts">
    import { page } from "$app/stores";
    import Badge from "$lib/general/Badge.svelte";
    import CloseButton from "$lib/general/CloseButton.svelte";
    import BackLink from "$lib/general/BackLink.svelte";
    import FindProduct from "$lib/products/FindProduct.svelte";
    import SubmitButton from "$lib/forms/SubmitButton.svelte";
    import { handleErrors } from "../../../assets/helper/handleErrors";
    import Input from "$lib/forms/Input.svelte";
    import { isAuthenticated } from "../../../store";

    interface Product {
        id: number,
        description: string,
        ean: number,
    }

    let productEan: number;
    let productPrice: number | null;
    let productData: Product;
    let priceData: { price: number };
    let priceIsAlreadyCreated: boolean = false;

    let formSubmitted: boolean = false;

    $: productPrice = priceData ? priceData.price : null;

    function submit(): void {
        if (! productPrice || ! productData) return;
        priceIsAlreadyCreated ? updatePrice() : fetchContent("POST");
    }

    function updatePrice(): void {
        if (! productData.id) return;
        fetchContent("PUT");
    }

    function fetchContent(method: string): void {
        const userId: number = 2; // TODO: add current user id
        const apiUrl: string = `/api/v1/price/${productData.id}/${userId}`

        const requestOptions = {
            method: method,
            headers: { 'Content-Type': 'application/json' },
            body: `{"price": ${productPrice}}`,
        };

        fetch(apiUrl, requestOptions)
            .then(handleErrors)
            .then(() => formSubmitted = true)
            .catch(error => console.error("Failed to fetch data:", error.message));
    }
</script>

{#if $isAuthenticated}
    <header>
        {#if ! formSubmitted}
            <h1 class="font-bold text-xl w-[90%] md:text-2xl xl:text-3xl">
                {$page.data.headline}
            </h1>
            <CloseButton
                url="/merchants"
                label="Erstellen eines Preises abbrechen" />
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
                    bind:productData={productData}
                    bind:productEan={productEan}
                    bind:priceData={priceData}
                    bind:priceIsAlreadyCreated={priceIsAlreadyCreated}
                    shouldFindPrice />

                {#if productData}
                    <section>
                        <Input
                            fieldName="productPrice"
                            label="Preis des Produktes"
                            type="number"
                            bind:value={productPrice} />

                        <div class="mt-10">
                            <SubmitButton on:submit={submit}/>
                        </div>
                    </section>
                {/if}
            {:else}
                <Badge />
                <h2 class="font-semibold text-lg lg:text-xl">
                    Dein Preis wurde erfolgreich gespeichert.
                </h2>
            {/if}
        </div>
    </main>
{/if}
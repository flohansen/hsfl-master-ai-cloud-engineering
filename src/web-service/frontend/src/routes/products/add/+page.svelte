<script lang="ts">
    import {page} from "$app/stores";
    import Badge from "$lib/general/Badge.svelte";
    import CloseButton from "$lib/general/CloseButton.svelte";
    import BackLink from "$lib/general/BackLink.svelte";
    import FindProduct from "$lib/products/FindProduct.svelte";
    import CreateProduct from "$lib/products/CreateProduct.svelte";
    import CreatePrice from "$lib/products/CreatePrice.svelte";

    interface Product {
        id: number,
        description: string,
        ean: number,
    }

    let productEan: number;
    let productData: Product;

    let formSubmitted: boolean = false;
    let eanSubmitted: boolean = false;
    let productSubmitted: boolean = false;
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
            {#if ! productSubmitted}
                <FindProduct
                    bind:eanSubmitted={eanSubmitted}
                    bind:productData={productData}
                    bind:productEan={productEan} />
            {/if}

            {#if eanSubmitted && ! formSubmitted}
                <CreateProduct
                    productEan={productEan}
                    bind:productData={productData}
                    bind:productSubmitted={productSubmitted}/>
            {/if}

            {#if productData}
                <CreatePrice
                    bind:productId={productData.id}
                    bind:formSubmitted={formSubmitted}/>
            {/if}
        {:else}
            <Badge />
            <h2 class="font-semibold text-lg lg:text-xl">
                Dein Produkt und dein Preis wurden erfolgreich gespeichert.
            </h2>
        {/if}
    </div>
</main>
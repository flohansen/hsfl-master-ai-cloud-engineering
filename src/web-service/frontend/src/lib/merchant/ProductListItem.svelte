<script lang="ts">
    import {handleErrors} from "../../assets/helper/handleErrors";
    import Trash from "../../assets/svg/Trash.svelte";

    export let price: { userId: number, productId: number, price: number } | undefined;
    export let product: { id: number, description: string, ean: number };

    function deletePrice() : void {
        const token: string | null = sessionStorage.getItem('access_token');

        if (! price || ! token) return;

        const apiUrl: string = `/api/v1/price/${price.productId}/${price.userId}`
        const requestOptions = {
            method: "DELETE",
            headers: { 'Authorization': `Bearer ${token}` },
        };

        fetch(apiUrl, requestOptions)
            .then(handleErrors)
            .then(()=> location.reload())
            .catch(error => console.error("Failed to fetch data:", error.message));
    }
</script>

<li class="w-full group cursor-pointer">
    <section class="bg-white max-h-max h-full transition-all ease-in-out duration-300 rounded-2xl w-full flex items-center justify-between group-hover:bg-blue-light/25">
        <div class="px-3 py-5 lg:p-6">
            <h3 class="font-semibold text-base transition-all ease-in-out duration-300 group-hover:text-blue-dark lg:text-lg">
                {product.description}
            </h3>
            <p class="text-sm mt-1">EAN: <span class="text-gray-dark">{product.ean}</span></p>
            {#if price}
                <p class="text-sm mt-1">Preis: <span class="text-gray-dark">{price.price} Є</span></p>
            {/if}
        </div>

        <button
            aria-label="Produkt und Preis löschen"
            on:click={deletePrice}
            class="h-full p-3 lg:p-6 border-l border-l-blue-light hidden group-hover:block text-blue-dark/50 transition-all ease-in-out duration-300 hover:text-blue-dark">
            <Trash classes="w-5 h-5" />
        </button>
    </section>
</li>
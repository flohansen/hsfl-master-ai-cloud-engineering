<script lang="ts">
    import {onMount} from "svelte";
    import {handleErrors} from "../../assets/helper/handleErrors";
    import Trash from "../../assets/svg/Trash.svelte";

    interface Product {
        id: number,
        description: string,
        ean: number,
    }

    export let price: { userId: number, productId: number, price: number };
    let product: Product = { id: 0, description: '', ean: 0 };

    const apiUrl = `/api/v1/product/${price.productId}`;

    onMount(async () => {
        fetch(apiUrl)
            .then(handleErrors)
            .then(data => product = data)
            .catch(error => console.error("Failed to fetch data:", error.message));
    });

    function deletePrice() : void {
        const apiUrl: string = `/api/v1/price/${price.productId}/${price.userId}`
        const requestOptions = {
            method: "DELETE",
            headers: { 'Content-Type': 'application/json' },
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
            <p class="text-sm mt-1">Preis: <span class="text-gray-dark">{price.price}</span></p>
        </div>

        <button
            aria-label="Produkt und Preis löschen"
            on:click={deletePrice}
            class="h-full p-3 lg:p-6 border-l border-l-blue-light hidden group-hover:block text-blue-dark/50 transition-all ease-in-out duration-300 hover:text-blue-dark">
            <Trash classes="w-5 h-5" />
        </button>
    </section>
</li>
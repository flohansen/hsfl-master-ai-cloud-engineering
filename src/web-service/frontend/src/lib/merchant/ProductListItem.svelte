<script lang="ts">
    import {onMount} from "svelte";
    import {handleErrors} from "../../assets/helper/handleErrors";

    interface Product {
        id: number,
        description: string,
        ean: number,
    }

    export let price: { productId: number, price: number };
    let product: Product = { id: 0, description: '', ean: 0 };

    const apiUrl = `/api/v1/product/${price.productId}`;

    onMount(async () => {
        fetch(apiUrl)
            .then(handleErrors)
            .then(data => product = data)
            .catch(error => console.error("Failed to fetch data:", error.message));
    });
</script>

<li class="bg-white w-full px-3 py-5 rounded-2xl flex items-center justify-between transition-all ease-in-out duration-300 group hover:bg-blue-light/25 lg:px-6 lg:py-8">
    <div class="w-full flex flex-col gap-x-4 gap-y-2 md:flex-wrap md:flex-row md:items-center md:justify-between">
        <div>
            <h3 class="font-semibold text-base transition-all ease-in-out duration-300 group-hover:text-blue-dark lg:text-lg">
                {product.description}
            </h3>
            <p class="text-sm mt-1">EAN: <span class="text-gray-dark">{product.ean}</span></p>
        </div>
        <p class="text-right">{price.price} €</p>
    </div>
</li>
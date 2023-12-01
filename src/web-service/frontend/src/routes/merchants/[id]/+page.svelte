<script lang="ts">
    import BackLink from "$lib/general/BackLink.svelte";

    interface Data {
        merchant: { name: string },
    }

    export let data: Data;
</script>

<header>
    <BackLink url="/merchants" label="Alle Supermärkte" />
</header>

<main>
    <div class="mx-5 bg-white rounded-xl p-4 lg:p-6">
        <h1 class="text-lg mb-4 font-semibold lg:text-xl xl:text-2xl">
            {data.merchant.name}
        </h1>
        TODO: show all linked products and its prices
    </div>
</main>

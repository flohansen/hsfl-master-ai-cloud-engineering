<script lang="ts">
    import {page} from "$app/stores";
    import Header from "$lib/general/Header.svelte";

    interface Merchant {
        id: number;
        name: string;
        email: string;
        role?: number;
        productsCount?: number;
    }

    interface Data {
        merchants: Merchant[],
        headline: string;
    }

    export let data: Data;
</script>

<Header headline="{$page.data.headline}"/>

<main>
    <ul class="px-5 mt-4 grid grid-cols-1 gap-y-4 lg:gap-y-6 lg:mt-6">
        {#if data.merchants}
            {#each data.merchants as merchant}
                <li class="bg-white w-full rounded-2xl flex items-center justify-between transition-all ease-in-out duration-300 group hover:bg-blue-light/25">
                    <a href="/merchants/{merchant.id}" class="w-full px-3 py-5 lg:px-6 lg:py-8">
                        <div class="text-left">
                            <h3 class="font-semibold text-base transition-all ease-in-out duration-300 group-hover:text-blue-dark lg:text-lg">
                                {merchant.name}
                            </h3>
                            <p class="text-xs text-gray-dark mt-1 lg:text-sm">
                                {#if merchant.productsCount}
                                    {merchant.productsCount} Produkte
                                {:else}
                                    0 Produkte
                                {/if}
                            </p>
                        </div>
                    </a>
                </li>
            {/each}
        {:else}
            <p>Es konnten keine Daten geladen werden.</p>
        {/if}
    </ul>
</main>

<script lang="ts">
    let totalPrice: number = 31.50;
    let totalMerchant: string = "Marktkauf";
    export let showNotice: boolean = false;
</script>

{#if showNotice}
    <div class="mt-10 border border-green-dark rounded-xl p-4 md:max-w-screen-sm md:mx-auto lg:p-6">
        <p class="text-green-dark text-sm lg:text-base">
            Dein Gesamteinkauf kostet <strong>{totalPrice} Euro.</strong>
        </p>
        <p class="text-gray-dark mt-3 text-sm flex flex-wrap items-center gap-2 lg:text-base lg:mt-4">
            Am günstigsten bei
            <strong class="text-green-dark font-semibold block px-3 py-1 rounded-2xl bg-green-light/25">
                {totalMerchant}
            </strong>
        </p>
    </div>
{/if}
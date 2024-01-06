<script lang="ts">
    import { page } from "$app/stores";
    import { goto } from "$app/navigation";
    import Badge from "$lib/general/Badge.svelte";
    import Add from "../../assets/svg/Add.svelte";
    import DataList from "$lib/profle/DataList.svelte";
    import DeleteAccountModal from "$lib/profle/DeleteAccountModal.svelte";

    interface Data {
        user: { id: number, email: string, name: string, role?: number },
        metaTitle: string,
        headline: string,
    }

    export let data: Data;
    let userRole: string = resolveUserRole();
    let successfulDeleted: boolean = false;

    function resolveUserRole(): string {
        if (! data.user.role) return '';

        switch (data.user.role) {
            case 1:
                return 'Händler:in';
            case 2:
                return 'Administrator:in';
            default:
                return 'Kund:in';
        }
    }

    function logout(): void {
        sessionStorage.removeItem('access_token');
        sessionStorage.removeItem('user_id');
        goto('/profile/login');
    }
</script>

<header>
    {#if ! successfulDeleted}
        <h1 class="font-bold text-xl md:text-2xl xl:text-3xl">
            {$page.data.metaTitle}
        </h1>
    {:else}
        <a href="/" class="flex gap-x-2 items-center text-gray-dark transition-all duration-300 ease-in-out hover:text-green-dark lg:gap-x-4">
            <svg class="w-6 h-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M6.75 15.75L3 12m0 0l3.75-3.75M3 12h18" />
            </svg>
            <p class="text-sm lg:text-base">Zu Startseite</p>
        </a>
    {/if}
    <button
        on:click={logout} class="flex gap-x-2 items-center text-gray-dark transition-all duration-300 ease-in-out hover:text-green-dark lg:gap-x-4 flex-row-reverse">
        <svg class="w-6 h-6 rotate-180" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M6.75 15.75L3 12m0 0l3.75-3.75M3 12h18" />
        </svg>
        <p class="text-sm lg:text-base">Logout</p>
    </button>
</header>

<main>
    <div class="mx-5 bg-white rounded-xl p-4 lg:p-6">
        {#if data.user.name}
            {#if ! successfulDeleted}
                <DataList
                    name="{data.user.name}"
                    email="{data.user.email}"
                    role="{userRole}" />

                <section class="mt-28 border-t-gray-dark/25 border-t pt-5 flex flex-col items-center gap-x-6 gap-y-4 md:flex-row">
                    {#if userRole !== 'Kund:in'}
                        <a href="/products/add"
                           class="text-green-dark flex items-center gap-x-2 font-medium transition-all ease-in-out duration-300 hover:text-green-light">
                            <Add classes="w-5 h-5"/>
                            Produkt hinzufügen
                        </a>
                        <a href="/prices/add"
                           class="text-green-dark flex items-center gap-x-2 font-medium transition-all ease-in-out duration-300 hover:text-green-light">
                            <Add classes="w-5 h-5"/>
                            Preis hinzufügen
                        </a>
                    {/if}
                    <DeleteAccountModal
                        bind:successfulDeleted={successfulDeleted} />
                </section>
            {:else}
                <Badge />
                <h2 class="font-semibold text-lg lg:text-xl">
                    Dein Account wurde erfolgreich gelöscht.
                </h2>
            {/if}
        {:else}
            <p>Leider bist du nicht eingeloggt.</p>
        {/if}
    </div>
</main>
<script lang="ts">
    import Profile from "../../assets/svg/Profile.svelte";
    import Trash from "../../assets/svg/Trash.svelte";
    import {handleErrors} from "../../assets/helper/handleErrors";
    import Badge from "$lib/general/Badge.svelte";
    import {page} from "$app/stores";

    interface Data {
        user: { id: number, email: string, name: string, role: number },
        metaTitle: string,
        headline: string,
    }

    export let data: Data;
    let userRole: string = resolveUserRole(data.user.role);
    let successfulDeleted: boolean = false;

    function resolveUserRole(role: number): string {
        switch (role) {
            case 1:
                return 'Händler:in';
            case 2:
                return 'Administrator:in';
            default:
                return 'Kund:in';
        }
    }

    function deleteAccount() : void {
        if (! data.user) return;

        const apiUrl: string = `/api/v1/user/${data.user.id}`
        const requestOptions = {
            method: "DELETE",
            headers: { 'Content-Type': 'application/json' },
        };

        fetch(apiUrl, requestOptions)
            .then(handleErrors)
            .then(()=> successfulDeleted = true)
            .catch(error => console.error("Failed to fetch data:", error.message));
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
</header>

<main>
    <div class="mx-5 bg-white rounded-xl p-4 lg:p-6">
        {#if data.user.name}
            {#if ! successfulDeleted}
                <section class="mt-10 lg:flex lg:items-center lg:justify-center lg:gap-x-16">
                    <figure class="mx-auto w-28 h-28 rounded-full bg-green-light/25 flex items-center justify-center lg:mx-0">
                        <Profile classes="w-12 h-12 text-green-dark"/>
                    </figure>
                    <div>
                        <h2 class="text-center mt-4 font-semibold text-lg mb-12 md:text-xl lg:mb-6 lg:text-left xl:text-2xl">
                            {data.user.name}
                        </h2>
                        <dl class="mx-auto max-w-[30rem] lg:max-w-none lg:mx-0">
                            <div class="mb-4 md:grid md:grid-cols-[auto,1fr] md:gap-x-4">
                                <dt class="text-gray-dark mb-2 md:mb-0">E-Mail:</dt>
                                <dd>{data.user.email}</dd>
                            </div>

                            <div class="mb-4 md:grid md:grid-cols-[auto,1fr] md:gap-x-4">
                                <dt class="text-gray-dark mb-2 md:mb-0">Benutzerrolle:</dt>
                                <dd>{userRole}</dd>
                            </div>
                        </dl>
                    </div>
                </section>

                <section class="mt-28 border-t-gray-dark/25 border-t pt-5">
                    <button
                        on:click={deleteAccount}
                        class="ml-auto mr-0 text-green-dark flex items-center gap-x-2 font-medium transition-all ease-in-out duration-300 hover:text-green-light">
                        <Trash classes="w-5 h-5"/>
                        Account löschen
                    </button>
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

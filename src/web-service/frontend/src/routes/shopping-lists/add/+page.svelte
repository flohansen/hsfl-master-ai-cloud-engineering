<script lang="ts">
    import {page} from "$app/stores";
    import Placeholder from "../../../assets/svg/Placeholder.svelte";
    import SubmitButton from "$lib/forms/SubmitButton.svelte";
    import {handleErrors} from "../../../assets/helper/handleErrors";
    import Checkmark from "../../../assets/svg/Checkmark.svelte";
    import Badge from "$lib/general/Badge.svelte";

    let listHeadline: string = '';
    let formSubmitted: boolean = false;

    function submit(): void {
        if ( listHeadline === '') return;

        const userId: number = 2; // TODO: dynamic user id of current logged in user
        const apiUrl: string = `/api/v1/shoppinglist/${userId}`
        const requestOptions = {
            method: "POST",
            headers: { 'Content-Type': 'application/json' },
            body: `{"description": "${listHeadline}"}`,
        };

        fetch(apiUrl, requestOptions)
            .then(handleErrors)
            .then(()=> formSubmitted = true)
            .catch(error => console.error("Failed to fetch data:", error.message));
    }
</script>

<header>
    {#if ! formSubmitted}
        <h1 class="font-bold text-xl md:text-2xl xl:text-3xl">
            {$page.data.metaTitle}
        </h1>

        <a
            href="/shopping-lists"
            aria-label="Erstellen der Einkaufsliste abbrechen"
            class="group rounded-full border-[1.5px] border-green-dark w-8 h-8 flex items-center justify-center transition-all ease-in-out duration-300 cursor-pointer hover:bg-green-dark">
            <span class="text-green-dark group-hover:text-white">
                <svg class="w-5 h-5" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                </svg>
            </span>
        </a>
    {:else}
        <a href="/shopping-lists" class="flex gap-x-2 items-center text-gray-dark transition-all duration-300 ease-in-out hover:text-green-dark lg:gap-x-4">
            <svg class="w-6 h-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M6.75 15.75L3 12m0 0l3.75-3.75M3 12h18" />
            </svg>
            <p class="text-sm lg:text-base">Zu deinen Einkaufslisten</p>
        </a>
    {/if}
</header>

<main>
    <div class="mx-5 bg-white rounded-xl p-4 lg:p-6">
        {#if ! formSubmitted}
            <section class="grid grid-cols-[3.5rem,auto] items-center gap-x-4 lg:gap-x-6 lg:grid-cols-[4rem,auto]">
                <figure class="bg-green-light/25 rounded-full w-14 h-14 flex items-center justify-center lg:w-16 lg:h-16">
                    <Placeholder classes="w-6 h-6 text-green-dark"/>
                </figure>
                <div class="w-full">
                    <label
                        for="listName"
                        class="text-sm text-gray-dark font-medium block mb-2">
                        Name der Einkaufsliste:*
                    </label>
                    <input
                        id="listName"
                        name="listName"
                        type="text"
                        required
                        placeholder="Name der Einkaufsliste"
                        bind:value={listHeadline}
                        class="text-sm rounded-lg border px-3 py-2 w-full text-left text-green-dark/75 flex items-center justify-between transition-all duration-300 ease-in-out hover:bg-blue-light/25 lg:px-4 lg:py-3"/>
                </div>
            </section>

            <div class="mt-10">
                <SubmitButton on:submit={submit} />
            </div>
        {:else}
            <Badge />
            <h2 class="font-semibold text-lg lg:text-xl">
                Deine Einkaufsliste wurde erfolgreich gespeichert.
            </h2>
        {/if}
    </div>
</main>
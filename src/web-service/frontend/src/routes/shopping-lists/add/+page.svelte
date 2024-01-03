<script lang="ts">
    import { page } from "$app/stores";
    import Placeholder from "../../../assets/svg/Placeholder.svelte";
    import SubmitButton from "$lib/forms/SubmitButton.svelte";
    import { handleErrors } from "../../../assets/helper/handleErrors";
    import Badge from "$lib/general/Badge.svelte";
    import CloseButton from "$lib/general/CloseButton.svelte";
    import BackLink from "$lib/general/BackLink.svelte";
    import Input from "$lib/forms/Input.svelte";
    import { isAuthenticated } from "../../../store";

    let listHeadline: string = '';
    let formSubmitted: boolean = false;

    function submit(): void {
        if ( listHeadline === '' ||  ! isAuthenticated) return;

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

{#if isAuthenticated}
    <header>
        {#if ! formSubmitted}
            <h1 class="font-bold text-xl md:text-2xl xl:text-3xl">
                {$page.data.metaTitle}
            </h1>

            <CloseButton url="/shopping-lists" label="Erstellen der Einkaufsliste abbrechen" />
        {:else}
            <BackLink url="/shopping-lists" label="Zu deinen Einkaufslisten" />
        {/if}
    </header>

    <main>
        <div class="mx-5 bg-white rounded-xl p-4 lg:p-6">
            {#if ! formSubmitted}
                <section class="grid grid-cols-[3.5rem,auto] items-center gap-x-4 lg:gap-x-6 lg:grid-cols-[4rem,auto]">
                    <figure class="bg-green-light/25 rounded-full w-14 h-14 flex items-center justify-center lg:w-16 lg:h-16">
                        <Placeholder classes="w-6 h-6 text-green-dark"/>
                    </figure>
                    <Input
                        fieldName="listName"
                        type="text"
                        label="Name der Einkaufsliste"
                        bind:value={listHeadline} />
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
{/if}
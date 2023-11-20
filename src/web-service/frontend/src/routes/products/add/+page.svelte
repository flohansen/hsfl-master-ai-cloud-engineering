<script lang="ts">
    import {page} from "$app/stores";
    import SubmitButton from "$lib/forms/SubmitButton.svelte";
    import {handleErrors} from "../../../assets/helper/handleErrors";
    import Badge from "$lib/general/Badge.svelte";
    import InputText from "$lib/forms/InputText.svelte";
    import CloseButton from "$lib/general/CloseButton.svelte";
    import BackLink from "$lib/general/BackLink.svelte";

    let productName: string = '';
    let formSubmitted: boolean = false;

    function submit(): void {
        if ( productName === '') return;

    }
</script>

<header>
    {#if ! formSubmitted}
        <h1 class="font-bold text-xl md:text-2xl xl:text-3xl">
            {$page.data.metaTitle}
        </h1>

        <CloseButton url="/profile" label="Erstellen eines Produktes abbrechen" />
    {:else}
        <BackLink url="/" label="Zur Startseite" />
    {/if}
</header>

<main>
    <div class="mx-5 bg-white rounded-xl p-4 lg:p-6">
        {#if ! formSubmitted}
            <section>
                <InputText
                    fieldName="productName"
                    label="Name des Produktes"
                    bind:value={productName} />
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
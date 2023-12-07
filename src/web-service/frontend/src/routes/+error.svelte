<script lang="ts">
    import { page } from '$app/stores';

    let headline: string = "Die Seite konnte leider nicht gefunden werden!";

    if ($page.error !== null && $page.error.message !== "Not Found") {
        headline = "Leider ist ein serverseitiger Fehler aufgetreten!";
    }
</script>

<svelte:head>
    <title>Fehler | Price Whisper</title>
</svelte:head>

<main class="my-0 flex flex-col justify-center h-screen">
    <div class="mx-5 bg-white rounded-xl p-4 lg:p-6">
        <h1 class="text-center font-semibold mb-6 text-2xl lg:text-3xl lg:mb-8">
            {headline}
        </h1>
        <p class="text-center mb-5 lg:mb-8">
            Leider konnte die von dir gewünschte Seite nicht gefunden oder aufgerufen werden.<br>
            Bitte überprüfe die URL auf Tippfehler und versuche es erneut.
        </p>
        <a
            href="/"
            class="mx-auto w-max border-[1.5px] bg-green-light text-white rounded-xl px-5 py-2 flex items-center justify-center gap-x-2 transition-all ease-in-out duration-300 hover:bg-green-dark">
            <span class="text-sm lg:text-base">Zur Starteite</span>
        </a>
    </div>
</main>
<script lang="ts">
    import Euro from "../assets/svg/Euro.svelte";
    import {page} from "$app/stores";
    import Header from "$lib/Header.svelte";
</script>

<Header headline="{$page.data.headline}"/>

<main class="mt-8 sm:ml-20 md:ml-24 lg:max-w-4xl lg:mx-auto lg:mt-10 xl:max-w-5xl">
    <section class="mx-5 rounded-2xl bg-blue-light/50 p-5 lg:p-8">
        <figure class="w-12 h-12 flex items-center justify-center rounded-full bg-blue-dark/40">
            <Euro classes="w-7 h-7 text-blue-dark"></Euro>
        </figure>
        <h2 class="text-lg font-semibold mt-5 lg:text-xl xl:text-2xl xl:mt-8">
            Spare bis zu 35% Prozent auf deinen nächsten Supermarkt-Einkauf!
        </h2>
        <p class="mt-3 text-base font-light lg:text-lg">
            Tempor qui cillum fugiat aliqua Lorem sint laborum duis. Esse elit dolore dolor quis mollit.
            Eu excepteur occaecat ea adipisicing ipsum proident nisi.
        </p>
    </section>
</main>


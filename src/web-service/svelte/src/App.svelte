<header class="mx-5 mt-16 flex items-center justify-between">
    <h1 class="font-bold text-lg">PriceWhisper</h1>
</header>

<main>
    <section class="mx-5 mt-8 rounded-2xl bg-blue-light/50 p-5">
        <figure class="w-12 h-12 flex items-center justify-center rounded-full bg-blue-dark/40">
            <Euro classes="w-7 h-7 text-blue-dark"></Euro>
        </figure>
        <h2 class="text-lg font-semibold mt-5">
            Spare bis zu 35% Prozent auf deinen nächsten Supermarkt-Einkauf!
        </h2>
        <p class="mt-3 text-base font-light">
            Tempor qui cillum fugiat aliqua Lorem sint laborum duis. Esse elit dolore dolor quis mollit.
            Eu excepteur occaecat ea adipisicing ipsum proident nisi.
        </p>
    </section>

    <NavBar/>
</main>

<script lang="ts">
    import Euro from "./assets/svg/Euro.svelte";
    import NavBar from "./lib/navigation/NavBar.svelte";
</script>



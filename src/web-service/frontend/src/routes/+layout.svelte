<script>
    import "../styles/app.css";
    import Header from "$lib/Header.svelte";
    import NavBar from "$lib/navigation/NavBar.svelte";
</script>

<Header headline="PriceWhisper"/>

<main class="sm:ml-20 md:ml-24 lg:max-w-4xl lg:mx-auto xl:max-w-5xl">
    <slot></slot>
</main>

<NavBar />
<header class="px-5 mt-16 flex items-center justify-between sm:ml-20 md:ml-24 lg:max-w-4xl lg:mx-auto xl:max-w-5xl">
    <h1 class="font-bold text-xl md:text-2xl xl:text-3xl">
        Deine Einkaufslisten
    </h1>
</header>

<main class="sm:ml-20 md:ml-24 lg:max-w-4xl lg:mx-auto xl:max-w-5xl">
    <h2 class="px-5 text-gray-dark text-sm font-medium mt-10 lg:text-base">
        Offene Einkaufslisten
    </h2>
    <ul class="px-5 mt-4 grid grid-cols-1 gap-y-4 lg:gap-y-6 lg:mt-6">
        <ShoppingListItem description="Geburtstagskuchen"/>
        <ShoppingListItem description="Frühstück mit Anne"/>
    </ul>
</main>

<NavBar/>
<script lang="ts">
    import NavBar from "./lib/navigation/NavBar.svelte";
    import ShoppingListItem from "./lib/ShoppingListItem.svelte";
</script>



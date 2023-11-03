<script>
    import ShoppingListItem from "$lib/ShoppingListItem.svelte";
</script>

<h2 class="px-5 text-gray-dark text-sm font-medium mt-6 lg:mt-10 lg:text-base">
    Offene Einkaufslisten
</h2>
<ul class="px-5 mt-4 grid grid-cols-1 gap-y-4 lg:gap-y-6 lg:mt-6">
    <ShoppingListItem description="Geburtstagskuchen"/>
    <ShoppingListItem description="Frühstück mit Anne"/>
</ul>

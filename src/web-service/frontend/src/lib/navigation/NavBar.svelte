<script lang="ts">
    import NavItem from "./NavItem.svelte";
    import Euro from "../../assets/svg/Euro.svelte";
</script>

<aside class="fixed bottom-0 bg-white w-full rounded-t-[2.5rem] px-5 py-4 drop-shadow-navBar sm:left-0 sm:top-0 sm:w-auto sm:rounded-none sm:rounded-r-2xl sm:px-4">
    <div class="hidden pt-10 pb-8 sm:block border-b-2 border-b-gray-dark/20">
        <figure class="w-12 h-12 bg-blue-dark/90 rounded-full flex items-center justify-center">
            <Euro classes="text-white w-6 h-6"/>
        </figure>
    </div>
    <nav aria-label="Hauptnavigation" class="flex items-center justify-between sm:flex-col sm:gap-6 sm:pt-8">
        <NavItem icon="home"/>
        <NavItem icon="list"/>
        <NavItem icon="cart"/>
        <NavItem icon="profile"/>
    </nav>
</aside>

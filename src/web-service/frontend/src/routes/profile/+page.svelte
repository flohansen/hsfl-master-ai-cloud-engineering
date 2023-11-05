<script>
import {page} from "$app/stores";
import Header from "$lib/Header.svelte";
</script>

<Header headline="{$page.data.headline}"/>

<main class="mt-8 sm:ml-20 md:ml-24 lg:max-w-4xl lg:mx-auto lg:mt-10 xl:max-w-5xl">
    <p class="mx-5">TODO: Add content about current user</p>
</main>

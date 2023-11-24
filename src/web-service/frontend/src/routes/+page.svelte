<script lang="ts">
    import {page} from "$app/stores";
    import Header from "$lib/general/Header.svelte";
    import Startpage from "$lib/startpage/Startpage.svelte";
</script>

<Header headline="{$page.data.headline}"/>

<main>
    <Startpage />
</main>


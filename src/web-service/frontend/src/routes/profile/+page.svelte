<script lang="ts">
    import {page} from "$app/stores";
    import Header from "$lib/general/Header.svelte";
</script>

<Header headline="{$page.data.headline}"/>

<main>
    <p class="mx-5">TODO: Add content about current user</p>
</main>

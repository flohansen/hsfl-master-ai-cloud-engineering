<script>
    import "../styles/app.css";
    import Header from "$lib/Header.svelte";
    import NavBar from "$lib/navigation/NavBar.svelte";
    import { page } from '$app/stores';
</script>

<svelte:head>
    <title>{$page.data.metaTitle} | Price Whisper</title>
</svelte:head>

<slot></slot>

<NavBar />

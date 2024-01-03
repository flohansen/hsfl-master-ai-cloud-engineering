<script lang="ts">
    import "../styles/app.css";
    import NavBar from "$lib/navigation/NavBar.svelte";
    import { page } from '$app/stores';
    import { onMount } from "svelte";
    import { setAuthenticationStatus } from "../store";
    import { decodeToken } from "../assets/helper/decodeToken";

    onMount(() => {
        const unsubscribe = page.subscribe(() => {
            if ($page.url.pathname === '/profile/login' || $page.url.pathname === '/profile/register') {
                return;
            }

            const token: string | null = sessionStorage.getItem('access_token');

            if (! token || isExpired(token)) {
                setAuthenticationStatus(false);
                window.location.href = '/profile/login';
                return;
            }
            setAuthenticationStatus(true);
        });

        return unsubscribe;
    });

    function isExpired(token: string): boolean {
        const decodedToken: string = decodeToken(token);
        const currentTimestamp: number = Math.floor(Date.now() / 1000);

        return decodedToken.exp < currentTimestamp;
    }
</script>

<svelte:head>
    <title>{$page.data.metaTitle} | Price Whisper</title>
</svelte:head>

<slot></slot>
<NavBar/>

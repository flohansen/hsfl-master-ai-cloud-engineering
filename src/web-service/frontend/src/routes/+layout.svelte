<script lang="ts">
    import "../styles/app.css";
    import NavBar from "$lib/navigation/NavBar.svelte";
    import { page } from '$app/stores';
    import { onMount } from "svelte";
    import { setAuthenticationStatus } from "../store";

    let currentPathname: string = $page.url.pathname;

    onMount(() => {
        const unsubscribe = page.subscribe(() => {
            currentPathname = $page.url.pathname;
            console.log('Current Pathname:', currentPathname);

            if (currentPathname === '/profile/login' || currentPathname === '/profile/register') {
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
        const decodedToken: string = parseJwt(token);
        const currentTimestamp: number = Math.floor(Date.now() / 1000);

        return decodedToken.exp < currentTimestamp;
    }

    function parseJwt(token: string): string {
        const base64Url: string = token.split('.')[1];
        const base64: string = base64Url.replace(/-/g, '+').replace(/_/g, '/');
        const jsonPayload: string = decodeURIComponent(atob(base64).split('').map(c => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2)).join(''));

        return JSON.parse(jsonPayload);
    }
</script>

<svelte:head>
    <title>{$page.data.metaTitle} | Price Whisper</title>
</svelte:head>

<slot></slot>
<NavBar/>

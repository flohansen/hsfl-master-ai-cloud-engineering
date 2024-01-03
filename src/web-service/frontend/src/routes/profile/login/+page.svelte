<script lang="ts">
    import { page } from "$app/stores";
    import { writable } from 'svelte/store';
    import { decodeToken } from "../../../assets/helper/decodeToken";
    import { setAuthenticationStatus } from "../../../store";
    import Close from "../../../assets/svg/Close.svelte";
    import Profile from "../../../assets/svg/Profile.svelte";
    import Input from "$lib/forms/Input.svelte";
    import SubmitButton from "$lib/forms/SubmitButton.svelte";

    let userMail: string;
    let userPassword: string;
    let errorMessage = writable('');

    async function login(event: Event): Promise<void> {
        event.preventDefault();

        if (! userMail || ! userPassword) return;

        const apiUrl: string = '/api/v1/authentication/login/';
        const requestOptions = {
            method: "POST",
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email: userMail, password: userPassword }),
        };

        fetch(apiUrl, requestOptions)
            .then(response => response.json())
            .then(data => {
                const { access_token } = data;
                const decodedToken: string = decodeToken(access_token);
                if (decodedToken) {
                    sessionStorage.setItem("access_token", access_token);
                    sessionStorage.setItem('user_id', decodedToken.id);
                    setAuthenticationStatus(true);
                    window.location.href = '/';
                }
            })
            .catch(error => {
                errorMessage.set("Ungültige E-Mail oder ungültiges Passwort. Bitte versuche es erneut.");
                console.error("Failed to fetch data:", error.message);
            });
    }
</script>

<header>
    <h1 class="font-bold text-xl md:text-2xl xl:text-3xl">
        {$page.data.headline}
    </h1>
</header>

<main>
    <div class="mx-5 bg-white rounded-xl p-4 lg:p-6">
        <section class="my-10 flex flex-col items-center">
            <figure class="mx-auto w-28 h-28 rounded-full bg-green-light/25 flex items-center justify-center">
                <Profile classes="w-12 h-12 text-green-dark"/>
            </figure>
            <form method="POST" class="mt-10 w-full max-w-screen-sm" on:submit={login}>
                <div class="grid grid-cols-1 gap-y-6 {$errorMessage ? 'mb-1' : 'mb-10'}">
                    <Input
                        fieldName="userMail"
                        type="text"
                        label="Deine E-Mail Adresse "
                        bind:value={userMail} />
                    <Input
                        fieldName="userPassword"
                        type="password"
                        label="Dein Passwort "
                        bind:value={userPassword} />
                </div>
                {#if $errorMessage}
                    <div class="grid grid-cols-[1.5rem,auto] items-center gap-x-2 mt-3">
                        <figure class="w-6 h-6 rounded-full flex items-center justify-center bg-red/25">
                            <Close classes="text-red w-4 h-4"/>
                        </figure>

                        <div class="text-sm text-gray-dark">
                            <p>{$errorMessage}</p>
                        </div>
                    </div>
                {/if}

                <SubmitButton
                    type="submit"
                    label="Anmelden" />
            </form>
        </section>
    </div>
</main>

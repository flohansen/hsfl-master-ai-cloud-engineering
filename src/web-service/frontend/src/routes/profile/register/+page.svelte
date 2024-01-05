<script lang="ts">
    import { page } from "$app/stores";
    import Profile from "../../../assets/svg/Profile.svelte";
    import Input from "$lib/forms/Input.svelte";
    import SubmitButton from "$lib/forms/SubmitButton.svelte";
    import { writable } from 'svelte/store';
    import Close from "../../../assets/svg/Close.svelte";
    import SelectUserRole from "$lib/forms/SelectUserRole.svelte";
    import BackLink from "$lib/general/BackLink.svelte";
    import Badge from "$lib/general/Badge.svelte";

    let userName: string;
    let userMail: string;
    let userPassword: string;
    let userRole: number;
    let errorMessage = writable('');
    let formSubmitted: boolean = false;

    async function register(event: Event): Promise<void> {
        event.preventDefault();

        if (! userMail || ! userPassword || ! userName || !userRole) return;

        const apiUrl: string = '/api/v1/authentication/register/';
        const requestOptions = {
            method: "POST",
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email: userMail, password: userPassword, name: userName, role: userRole }),
        };

        fetch(apiUrl, requestOptions)
            .then(response => response.json())
            .then(() => formSubmitted = true)
            .catch(error => {
                errorMessage.set("Leider ist etwas schief gelaufen. Bitte versuche es erneut.");
                console.error("Failed to fetch data:", error.message);
            });
    }
</script>

<header>
    {#if ! formSubmitted}
        <h1 class="font-bold text-xl md:text-2xl xl:text-3xl">
            {$page.data.headline}
        </h1>
    {:else}
        <BackLink
            url="/profile/login"
            label="Zum Login" />
    {/if}
</header>

<main>
    <div class="mx-5 bg-white rounded-xl p-4 lg:p-6">
        {#if ! formSubmitted}
            <section class="my-10 flex flex-col items-center">
                <figure class="mx-auto w-28 h-28 rounded-full bg-green-light/25 flex items-center justify-center">
                    <Profile classes="w-12 h-12 text-green-dark"/>
                </figure>
                <form method="POST" class="mt-10 w-full max-w-screen-sm" on:submit={register}>
                    <div class="grid grid-cols-1 gap-y-6 {$errorMessage ? 'mb-1' : 'mb-10'}">
                        <Input
                            fieldName="userName"
                            type="text"
                            label="Dein Benutzername "
                            bind:value={userName} />
                        <Input
                            fieldName="userMail"
                            type="email"
                            label="Deine E-Mail Adresse "
                            bind:value={userMail} />
                        <Input
                            fieldName="userPassword"
                            type="password"
                            label="Dein Passwort "
                            bind:value={userPassword} />

                        <SelectUserRole bind:justValue={userRole} />
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
                        label="Registieren" />
                </form>
                <p class="text-center text-sm text-gray-dark mt-6">
                    Bist du schon registiert?<br>
                    <a href="/profile/login" class="font-semibold text-green-dark transition-all ease-in-out duration-300 hover:text-green-light">
                        Dann kannst du dich hier einloggen.
                    </a>
                </p>
            </section>
        {:else}
            <Badge />
            <h2 class="font-semibold mb-6 text-lg lg:text-xl lg:mb-8">
                Dein Account wurde erfolgreich angelegt.
            </h2>
            <BackLink
                url="/profile/login"
                label="Zum Login"
                reverse />
        {/if}
    </div>
</main>

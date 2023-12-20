<script lang="ts">
    import {page} from "$app/stores";
    import Profile from "../../../assets/svg/Profile.svelte";
    import Input from "$lib/forms/Input.svelte";
    import SubmitButton from "$lib/forms/SubmitButton.svelte";

    let userMail: string;
    let userPassword: string;
</script>

<header>
    <h1 class="font-bold text-xl md:text-2xl xl:text-3xl">
        {$page.data.metaTitle}
    </h1>
</header>

<main>
    <div class="mx-5 bg-white rounded-xl p-4 lg:p-6">
        <section class="my-10 flex flex-col items-center">
            <figure class="mx-auto w-28 h-28 rounded-full bg-green-light/25 flex items-center justify-center">
                <Profile classes="w-12 h-12 text-green-dark"/>
            </figure>
            <form class="mt-10 w-full max-w-screen-sm">
                <div class="grid grid-cols-1 gap-y-6">
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
                <SubmitButton label="Anmelden" />
            </form>
        </section>
    </div>
</main>

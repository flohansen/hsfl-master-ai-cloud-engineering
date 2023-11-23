<script lang="ts">
    import Profile from "../../assets/svg/Profile.svelte";
    export let name: string;
    export let email: string;
    export let role: string;
</script>

<section class="mt-10 flex flex-col items-center">
    <figure class="mx-auto w-28 h-28 rounded-full bg-green-light/25 flex items-center justify-center">
        <Profile classes="w-12 h-12 text-green-dark"/>
    </figure>
    <div>
        <h2 class="text-center mt-4 font-semibold text-lg mb-12 md:text-xl xl:text-2xl">
            {name}
        </h2>
        <dl class="flex flex-col justify-center">
            <div class="mb-4 flex gap-x-4">
                <dt class="text-gray-dark mb-2 md:mb-0">E-Mail:</dt>
                <dd>{email}</dd>
            </div>

            <div class="mb-4 flex gap-x-4">
                <dt class="text-gray-dark mb-2 md:mb-0">Benutzerrolle:</dt>
                <dd>{role}</dd>
            </div>
        </dl>
    </div>
</section>
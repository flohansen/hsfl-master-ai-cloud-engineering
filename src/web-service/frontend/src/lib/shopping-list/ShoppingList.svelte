<script lang="ts">
    import Trash from "../../assets/svg/Trash.svelte";
    import {handleErrors} from "../../assets/helper/handleErrors";

    export let description: string = 'Name der Einkaufsliste';
    export let id: number = 0;
    export let hideDeleteButton: boolean = false;

    function deleteList() : void {
        const token: string | null = sessionStorage.getItem('access_token');
        if (hideDeleteButton || ! token) return;

        const apiUrl: string = `/api/v1/shoppinglist/${id}`
        const requestOptions = {
            method: "DELETE",
            headers: { 'Authorization': `Bearer ${token}` },
        };

        fetch(apiUrl, requestOptions)
            .then(handleErrors)
            .then(()=> location.reload())
            .catch(error => console.error("Failed to fetch data:", error.message));
    }
</script>

<li class="bg-white w-full rounded-2xl flex items-center justify-between transition-all ease-in-out duration-300 group hover:bg-blue-light/25">
    <a href="/shopping-lists/{id}" class="w-full p-3 lg:p-6 flex items-center gap-x-4 lg:gap-x-6">
        <figure class="bg-green-light/25 rounded-full w-14 h-14 flex items-center justify-center lg:w-16 lg:h-16">
            <span class="text-2xl lg:text-3xl">🥗</span>
        </figure>
        <div class="text-left">
            <h3 class="font-semibold text-base transition-all ease-in-out duration-300 group-hover:text-blue-dark lg:text-lg">
                {description}
            </h3>
            <p class="text-xs text-gray-dark mt-1 lg:text-sm">
                Hier ist Platz für eine Beschreibung.
            </p>
        </div>
    </a>

    {#if ! hideDeleteButton}
        <button
            aria-label="Einkaufsliste löschen"
            on:click={deleteList}
            class="h-full p-3 lg:p-6 border-l border-l-blue-light hidden group-hover:block text-blue-dark/50 transition-all ease-in-out duration-300 hover:text-blue-dark">
            <Trash classes="w-5 h-5" />
        </button>
    {/if}
</li>
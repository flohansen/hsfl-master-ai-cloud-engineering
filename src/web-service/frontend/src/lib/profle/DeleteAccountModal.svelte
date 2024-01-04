<script lang="ts">
    import Trash from "../../assets/svg/Trash.svelte";
    import Modal from "$lib/general/Modal.svelte";
    import {handleErrors} from "../../assets/helper/handleErrors.js";

    let isOpen: boolean = false;
    export let successfulDeleted: boolean = false;
    export let  user: { id: number, email: string, name: string, role: number };

    function deleteAccount() : void {
        if (! user) return;

        const apiUrl: string = `/api/v1/user/${user.id}`
        const requestOptions = {
            method: "DELETE",
            headers: { 'Content-Type': 'application/json' },
        };

        fetch(apiUrl, requestOptions)
            .then(handleErrors)
            .then(()=> successfulDeleted = true)
            .catch(error => console.error("Failed to fetch data:", error.message));
    }
</script>

<button
    on:click={() => isOpen = !isOpen}
    class="mx-auto text-green-dark flex items-center gap-x-2 font-medium transition-all ease-in-out duration-300 hover:text-green-light md:ml-auto md:mr-0">
    <Trash classes="w-5 h-5"/>
    Account löschen
</button>

<Modal
    submitLabel="Account löschen"
    on:submit={deleteAccount}
    bind:isOpen>
    <h3 class="text-lg font-semibold text-center lg:text-xl">
        Möchtest du wirklich deinen Account löschen?
    </h3>
</Modal>
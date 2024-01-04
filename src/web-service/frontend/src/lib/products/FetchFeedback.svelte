<script lang="ts">
    import Close from "../../assets/svg/Close.svelte";
    import Checkmark from "../../assets/svg/Checkmark.svelte";

    type FeedbackType = "successful" | "unsuccessful" | "notFound";
    export let label: string;
    export let feedback: FeedbackType;

    $: feedbackClass = feedback === 'successful'
        ? 'bg-green-light/25'
        : (feedback === 'unsuccessful' ? 'bg-red/25' : 'bg-gray-light');
</script>

<div class="grid grid-cols-[1.5rem,auto] gap-x-2 mt-3">
    <figure class={`w-6 h-6 rounded-full flex items-center justify-center ${feedbackClass}`}>
        {#if feedback === 'successful'}
            <Checkmark classes="text-green-dark w-4 h-4"/>
        {:else if feedback === 'notFound'}
            <span class="text-gray-dark">?</span>
        {:else}
            <Close classes="text-red w-4 h-4"/>
        {/if}
    </figure>

    <div class="text-sm text-gray-dark">
        <p>{label}</p>
        <slot></slot>
    </div>
</div>

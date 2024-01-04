<script lang="ts">
    export let url: string;
    export let label: string;
    export let reverse: boolean = false;
</script>

<a href="{url}" class="flex gap-x-2 items-center text-gray-dark transition-all duration-300 ease-in-out hover:text-green-dark lg:gap-x-4 { reverse ? 'flex-row-reverse' : '' }">
    <svg class="w-6 h-6 { reverse ? 'rotate-180' : '' }" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" d="M6.75 15.75L3 12m0 0l3.75-3.75M3 12h18" />
    </svg>
    <p class="text-sm lg:text-base">{label}</p>
</a>
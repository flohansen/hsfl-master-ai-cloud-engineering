<script lang="ts">
    export let url: string;
    export let label: string;
</script>

<a
    href="{url}"
    aria-label="{label}"
    class="group rounded-full border-[1.5px] border-green-dark w-8 h-8 flex items-center justify-center transition-all ease-in-out duration-300 cursor-pointer hover:bg-green-dark">
    <span class="text-green-dark group-hover:text-white">
        <svg class="w-5 h-5" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
        </svg>
    </span>
</a>
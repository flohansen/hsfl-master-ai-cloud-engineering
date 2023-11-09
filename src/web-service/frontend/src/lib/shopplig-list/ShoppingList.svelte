<script lang="ts">
    export let description: string = 'Name der Einkaufsliste';
    export let id: number = 0;
</script>

<li>
    <a href="/shopping-lists/{id}" class="bg-white w-full rounded-2xl p-3 flex items-center gap-x-4 transition-all ease-in-out duration-300 group hover:bg-blue-light/25 lg:p-6 lg:gap-x-6">
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
</li>
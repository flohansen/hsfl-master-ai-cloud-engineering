<button
        aria-label="{currentOption.ariaLabel}"
        aria-current="{isActive}"
        class="{isActive ? 'bg-green-light/25 text-green-dark hover:bg-green-light/50' : ' bg-gray-light text-black hover:bg-gray-dark/25'}
                w-14 h-14 rounded-full flex items-center justify-center transition-all ease-in-out duration-300">
        <svelte:component this="{currentOption.component}" classes="w-6 h-6"/>
</button>

<script lang="ts">
    import Home from "../../assets/svg/Home.svelte";
    import List from "../../assets/svg/List.svelte";
    import Cart from "../../assets/svg/Basket.svelte";
    import Profile from "../../assets/svg/Profile.svelte";

    interface NavIcons {
        label: string;
        component: typeof Home | typeof List | typeof  Cart | typeof  Profile,
        ariaLabel: string
    }

    const options: NavIcons[] = [
        { label: 'home', component: Home, ariaLabel: "Zur Startseite" },
        { label: 'list', component: List, ariaLabel: "Zu deinen Einkauftslisten" },
        { label: 'cart', component: Cart, ariaLabel: "Zu allen Händlern und deren Produkte" },
        { label: 'profile', component: Profile, ariaLabel: "Zu deinem Profil" },
    ]

    function getOptionByLabel(labelToFind: string): NavIcons {
        return options.find(option => option.label === labelToFind) ?? options[0];
    }

    export let icon = 'home';
    let currentOption = getOptionByLabel(icon);
    let isActive: boolean = false;
</script>
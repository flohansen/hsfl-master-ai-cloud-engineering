import type { PageLoad } from './$types';

export const load: PageLoad = () : object => {
    return {
        metaTitle: 'Login',
        headline: 'Logge dich ein',
    };
}
import type { PageLoad } from './$types';

export const load: PageLoad = () : object => {
    return {
        metaTitle: 'Registrierung',
        headline: 'Registrierung eines neuen Accounts',
    };
}
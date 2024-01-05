import type { PageLoad } from './$types';
import { fetchHelper } from "../assets/helper/fetchHelper";
import { checkAuthentication } from "../assets/helper/checkAuthentication";

export const load: PageLoad = async () : Promise<object> => {
    await checkAuthentication();

    const apiUrl: string = `/api/v1/shoppinglist/${sessionStorage.getItem('user_id')}`;
    const lists: object[] = await fetchHelper(apiUrl);

    return {
        lists: lists ? lists.slice(0, 3) : [],
        metaTitle: 'Startseite',
        headline: 'Price Whisper',
    };
};

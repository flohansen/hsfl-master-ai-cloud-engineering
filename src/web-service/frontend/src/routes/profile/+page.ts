import type { PageLoad } from './$types';
import { fetchHelper } from "../../assets/helper/fetchHelper";
import { checkAuthentication } from "../../assets/helper/checkAuthentication";

export const load: PageLoad = async () : Promise<object> => {
    await checkAuthentication();

    const apiUrl: string = `/api/v1/user/${sessionStorage.getItem('user_id')}`;
    const user: object[] = await fetchHelper(apiUrl);

    return {
        user: user ?? [],
        metaTitle: 'Dein Profil',
        headline: 'Dein Profil',
    };
};

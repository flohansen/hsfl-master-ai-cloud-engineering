import type { PageLoad } from './$types';
import { checkAuthentication } from "../../../assets/helper/checkAuthentication";

export const load: PageLoad = async () : Promise<object> => {
    await checkAuthentication();

    return {
        metaTitle: 'Neue Einkaufsliste erstellen',
    };
};
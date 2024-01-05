import { fetchHelper } from "../../assets/helper/fetchHelper";

export const load = async (): Promise<object> => {
    const apiUrl: string = `/api/v1/user/${sessionStorage.getItem('user_id')}`;
    const user: object[] = await fetchHelper(apiUrl);

    return {
        user: user ?? [],
        metaTitle: 'Deine Profil-Einstellungen',
        headline: 'Dein Profil',
    };
};

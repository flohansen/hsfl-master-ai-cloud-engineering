import { handleErrors } from "../../assets/helper/handleErrors";
import { isAuthenticated } from "../../store";

export const load = async (): Promise<Promise<object> | undefined> => {
    if (! isAuthenticated) return;

    const token: string | null = sessionStorage.getItem('access_token');
    const userId: string | null  = sessionStorage.getItem('user_id');

    if (! token || ! userId) return;

    const apiUrl: string = `/api/v1/user/${userId}`;
    const requestOptions: object = {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${token}`
        },
    };

    return fetch(apiUrl, requestOptions)
        .then(handleErrors)
        .then(user => {
            return {
                user: user ?? [],
                metaTitle: 'Deine Profil-Einstellungen',
                headline: 'Dein Profil',
            };
        })
        .catch(error => {
            console.error("Failed to fetch user data:", error.message);
            return {
                user: [],
                metaTitle: 'Deine Profil-Einstellungen',
                headline: 'Dein Profil',
            };
        });
};
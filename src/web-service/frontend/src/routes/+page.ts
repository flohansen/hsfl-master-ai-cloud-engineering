import { handleErrors } from "../assets/helper/handleErrors";
import { isAuthenticated } from "../store";

export const load = (): Promise<object> | undefined => {
    const userId: string | null = sessionStorage.getItem('user_id');
    const token: string | null = sessionStorage.getItem('access_token');

    if (! token || ! userId || ! isAuthenticated) return;

    const requestOptions: object = {
        method: "GET",
        headers: {
            'Authorization': `Bearer ${token}`
        },
    };

    const apiUrl: string = `/api/v1/shoppinglist/${userId}`;

    return fetch(apiUrl, requestOptions)
        .then(handleErrors)
        .then(lists => {
            return {
                lists: lists.slice(0, 3),
                metaTitle: 'Startseite',
                headline: 'Price Whisper',
            };
        })
        .catch(error => {
            console.error("Failed to fetch shopping lists data:", error.message);
            return {
                metaTitle: 'Error',
                headline: 'Price Whisper',
            };
        });
};

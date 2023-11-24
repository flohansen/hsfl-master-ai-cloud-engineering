import { handleErrors } from '../../assets/helper/handleErrors';

// Loads all shopping lists of the current user.
export const load = (): Promise<object> => {
    const id: number = 2; // TODO: dynamic user id of current logged in user
    const apiUrl: string = `/api/v1/shoppinglist/${id}`;

    return fetch(apiUrl)
        .then(handleErrors)
        .then(lists => {
            return {
                lists,
                metaTitle: 'Auflistung deiner Einkaufslisten',
                headline: 'Deine Einkaufslisten',
            };
        })
        .catch(error => {
            console.error("Failed to fetch shopping lists data:", error.message);
            return {
                lists: [],
                metaTitle: 'Error',
                headline: 'Leider ist ein Fehler aufgetreten',
            };
        });
};
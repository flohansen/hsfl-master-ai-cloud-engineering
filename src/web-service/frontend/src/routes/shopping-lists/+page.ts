import { handleErrors } from '../../assets/helper/handleErrors';
import { isAuthenticated } from "../../store";

interface List {
    id: number;
    userId: number;
    description: string;
    completed?: boolean;
}

export const load = (): Promise<object> | undefined => {
    const userId: string | null = sessionStorage.getItem('user_id');
    const token: string | null = sessionStorage.getItem('access_token');

    if (! token || ! userId || ! isAuthenticated) return;

    const apiUrl: string = `/api/v1/shoppinglist/${userId}`;

    const requestOptions: object = {
        method: "GET",
        headers: {
            'Authorization': `Bearer ${token}`
        },
    };

    return fetch(apiUrl, requestOptions)
        .then(handleErrors)
        .then(lists => {
            return {
                completedLists: filterListsByCompletedState(true, lists),
                incompleteLists: filterListsByCompletedState(false, lists),
                metaTitle: 'Auflistung deiner Einkaufslisten',
                headline: 'Deine Einkaufslisten',
            };
        })
        .catch(error => {
            console.error("Failed to fetch shopping lists data:", error.message);
            return {
                completedLists: [],
                incompleteLists: [],
                metaTitle: 'Error',
                headline: 'Leider ist ein Fehler aufgetreten',
            };
        });
};

function filterListsByCompletedState(completedState: boolean, lists: List[]): List[] {
    return lists.filter((list: List) =>
        (list.completed === completedState) || (list.completed === undefined && !completedState));
}

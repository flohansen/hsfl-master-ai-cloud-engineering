import { handleErrors } from '../../assets/helper/handleErrors';
import { isAuthenticated } from "../../store";

interface List {
    id: number;
    userId: number;
    description: string;
    completed?: boolean;
}

// Loads all shopping lists of the current user.
export const load = (): Promise<object> | undefined => {
    if (! isAuthenticated) {
        return;
    }

    const id: number = 2; // TODO: dynamic user id of the current logged-in user
    const apiUrl: string = `/api/v1/shoppinglist/${id}`;

    return fetch(apiUrl)
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

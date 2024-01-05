import type { PageLoad } from './$types';
import { checkAuthentication } from "../../assets/helper/checkAuthentication";
import { fetchHelper } from "../../assets/helper/fetchHelper";

interface List {
    id: number;
    userId: number;
    description: string;
    completed?: boolean;
}

export const load: PageLoad = async () : Promise<object> => {
    await checkAuthentication();

    const apiUrl: string = `/api/v1/shoppinglist/${sessionStorage.getItem('user_id')}`;
    const lists: any = await fetchHelper(apiUrl);

    if (! lists) return data;

    return data(
        filterListsByCompletedState(true, lists),
        filterListsByCompletedState(false, lists)
    );
};

function filterListsByCompletedState(completedState: boolean, lists: List[]): List[] {
    return lists.filter((list: List) =>
        (list.completed === completedState) || (list.completed === undefined && !completedState));
}

const data = (completedLists: object[] = [], incompleteLists: object[] = []): object => {
    return {
        completedLists: completedLists,
        incompleteLists: incompleteLists,
        metaTitle: 'Deine Einkaufslisten',
        headline: 'Deine Einkaufslisten',
    };
};

import { handleErrors } from '../../../assets/helper/handleErrors';

// Loads shopping list data and entries for a given shopping list id.
export const load = async (context: { params: { id: string } }): Promise<object> => {
    const { id } = context.params;
    const apiUrlList: string = `/api/v1/shoppinglist/${id}/2`;
    const apiUrlEntries: string = `/api/v1/shoppinglistentries/${id}`;

    const [list, entries] = await Promise.all([
        fetch(apiUrlList).then(handleErrors),
        fetch(apiUrlEntries).then(handleErrors),
    ]);

    if (entries === null || list === null) {
        console.warn('Entries or Shoppinglist not found');
    }

    return {
        list: list ?? [],
        entries: entries ?? [],
        metaTitle: 'Liste: ' + list?.description,
    };
};
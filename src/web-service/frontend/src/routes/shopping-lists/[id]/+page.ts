export const load = async (context: { params: { id: string } }): Promise<object> => {
    const { id } = context.params;

    const responseList = await fetch(`/api/v1/shoppinglist/${id}/2`);
    const list = await responseList.json();

    const responseEntries = await fetch(`/api/v1/shoppinglistentries/${id}`);
    const entries = await responseEntries.json();

    return {
        list,
        entries,
        metaTitle: 'Liste: ' + list.description,
    };
};
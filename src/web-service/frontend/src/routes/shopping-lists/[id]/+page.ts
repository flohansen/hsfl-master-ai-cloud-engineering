export const load = async (context: { params: { id: string } }): Promise<object> => {
    const { id } = context.params;
    const responseList = await fetch(`/api/v1/shoppinglist/${id}/2`);
    const responseProducts = await fetch(`/api/v1/shoppinglistentries/${id}`);
    const list = responseList.json();

    return {
        list,
    };
};
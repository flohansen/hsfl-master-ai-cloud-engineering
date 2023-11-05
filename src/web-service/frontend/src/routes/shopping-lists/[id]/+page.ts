export const load = async (context: { params: { id: string } }): Promise<object> => {
    const { id } = context.params;
    const responseList = await fetch(`http://127.0.0.1:8080/api/v1/shoppinglist/${id}/2`);
    // const responseProducts = await fetch(`http://127.0.0.1:8080/api/v1/shoppinglistentries/${id}`);

    return {
        responseList,
    };
};
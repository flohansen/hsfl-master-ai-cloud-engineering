export const load = async (): Promise<object> => {
    const res = await fetch("http://127.0.0.1:8080/api/v1/shoppinglist/2");
    const lists = await res.json();

    return {
        lists,
        metaTitle: 'Auflistung deiner Einkaufslisten',
        headline: 'Deine Einkaufslisten',
    };
};
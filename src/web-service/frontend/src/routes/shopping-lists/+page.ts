export const load = async (): Promise<object> => {
    const res = await fetch("/api/v1/shoppinglist/2");
    const lists = await res.json();

    return {
        lists,
        metaTitle: 'Auflistung deiner Einkaufslisten',
        headline: 'Deine Einkaufslisten',
    };
};
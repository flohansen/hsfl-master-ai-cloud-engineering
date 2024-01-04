import { handleErrors } from '../../../assets/helper/handleErrors';
import { sortProducts } from "../../../assets/helper/sortProducts";
import { isAuthenticated } from "../../../store";

interface List {
    id: number,
    description: string,
    completed?: boolean,
}

interface Entry {
    productId: number,
    count: number,
    checked?: boolean
}

interface Product {
    id: number,
    description: string,
    ean: number,
}

export const load = async (context: { params: { id: string } }): Promise<Promise<object> | undefined> => {
    if (! isAuthenticated) return;

    const { id } = context.params;
    const token: string | null = sessionStorage.getItem('access_token');

    if (! token || ! id) return;

    const apiUrlList: string = `/api/v1/shoppinglist/${id}/2`;
    const apiUrlEntries: string = `/api/v1/shoppinglistentries/${id}`;

    const requestOptions: object = {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${token}`
        },
    };

    try {
        const [list, entries] = await Promise.all([
            fetch(apiUrlList, requestOptions).then(handleErrors) as Promise<List>,
            fetch(apiUrlEntries, requestOptions).then(handleErrors) as Promise<Entry[]>,
        ]);

        let sortedProducts: Product[] = await sortProducts(entries);

        return {
            list: list,
            entries: entries ?? [],
            products: sortedProducts,
            metaTitle: 'Liste: ' + list?.description,
        };
    } catch (error) {
        return {
            merchant: null,
            prices: [],
            products: [],
            metaTitle: 'Leider ist ein Fehler aufgetreten.',
        };
    }
};
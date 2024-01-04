import { handleErrors } from '../../../assets/helper/handleErrors';
import { isAuthenticated } from "../../../store";
import {sortProducts} from "../../../assets/helper/sortProducts";

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
    const apiUrlList: string = `/api/v1/shoppinglist/${id}/2`;
    const apiUrlEntries: string = `/api/v1/shoppinglistentries/${id}`;

    try {
        const [list, entries] = await Promise.all([
            fetch(apiUrlList).then(handleErrors) as Promise<List>,
            fetch(apiUrlEntries).then(handleErrors) as Promise<Entry[]>,
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
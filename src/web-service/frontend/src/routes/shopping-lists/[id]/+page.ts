import { handleErrors } from '../../../assets/helper/handleErrors';
import { sortProducts } from "../../../assets/helper/sortProducts";
import { checkAuthentication } from "../../../assets/helper/checkAuthentication";

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

export const load = async (context: { params: { id: string } }): Promise<Promise<object>> => {
    await checkAuthentication();

    const { id } = context.params;
    const token: string | null = sessionStorage.getItem('access_token');

    if (! token || ! id) return data;

    const apiUrlList: string = `/api/v1/shoppinglist/${id}/2`;
    const apiUrlEntries: string = `/api/v1/shoppinglistentries/${id}`;
    const requestOptions: object = { headers: { 'Authorization': `Bearer ${token}` }};

    try {
        const [list, entries] = await Promise.all([
            fetch(apiUrlList, requestOptions).then(handleErrors) as Promise<List>,
            fetch(apiUrlEntries, requestOptions).then(handleErrors) as Promise<Entry[]>,
        ]);

        let sortedProducts: Product[] = await sortProducts(entries);
        return data(list, entries ?? [], sortedProducts);

    } catch (error) {
        return data;
    }
};

const data = (list: any = [], entries: object[] = [], products: object[] = []): object => {
    return {
        list: list,
        entries: entries,
        products: products,
        metaTitle: 'Liste: ' + list?.description,
    };
};
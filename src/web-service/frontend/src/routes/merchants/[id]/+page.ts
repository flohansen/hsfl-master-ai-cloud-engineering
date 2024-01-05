import type { PageLoad } from './$types';
import { handleErrors } from '../../../assets/helper/handleErrors';
import { sortProducts } from "../../../assets/helper/sortProducts";
import { checkAuthentication } from "../../../assets/helper/checkAuthentication";

interface Merchant {
    id: number;
    name: string;
    role?: number;
    productsCount?: number;
}

interface Price {
    userId: number;
    productId: number;
    price: number;
}

interface Product {
    id: number,
    description: string,
    ean: number,
}

export const load: PageLoad = async (context: { params: { id: string } }) : Promise<object> => {
    await checkAuthentication();

    const { id } = context.params;
    const token: string | null = sessionStorage.getItem('access_token');
    const apiUrlMerchant: string = `/api/v1/user/${id}`;
    const apiUrlPrices: string = `/api/v1/price/user/${id}`;

    if (! token || ! id) return data;

    const requestOptions: object = { headers: { 'Authorization': `Bearer ${token}` }};

    try {
        const [merchant, prices] = await Promise.all([
            fetch(apiUrlMerchant, requestOptions).then(handleErrors) as Promise<Merchant>,
            fetch(apiUrlPrices, requestOptions).then(handleErrors) as Promise<Price[]>,
        ]);

        let sortedProducts: Product[] = await sortProducts(prices);

        return data(merchant,
            prices ?? [],
            sortedProducts ?? []);
    } catch (error) {
        return data;
    }
};

const data = (merchant: any = [], prices: object[] = [], sortedProducts: object[] = []): object => {
    return {
        merchant: merchant,
        prices: prices ?? [],
        products: sortedProducts,
        metaTitle: merchant?.name ?? 'Leider ist ein Fehler aufgetreten.',
    };
};
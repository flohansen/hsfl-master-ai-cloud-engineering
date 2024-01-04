import { handleErrors } from '../../../assets/helper/handleErrors';
import { isAuthenticated } from "../../../store";
import {sortProducts} from "../../../assets/helper/sortProducts";

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

export const load = async (context: { params: { id: string } }): Promise<Promise<object> | undefined> => {
    if (! isAuthenticated) return;

    const { id } = context.params;
    const apiUrlMerchant: string = `/api/v1/user/${id}`;
    const apiUrlPrices: string = `/api/v1/price/user/${id}`;

    try {
        const [merchant, prices] = await Promise.all([
            fetch(apiUrlMerchant).then(handleErrors) as Promise<Merchant>,
            fetch(apiUrlPrices).then(handleErrors) as Promise<Price[]>,
        ]);

        let sortedProducts: Product[] = await sortProducts(prices);

        return {
            merchant: merchant,
            prices: prices ?? [],
            products: sortedProducts,
            metaTitle: merchant?.name
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


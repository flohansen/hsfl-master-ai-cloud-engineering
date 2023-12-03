import { handleErrors } from '../../../assets/helper/handleErrors';

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

export const load = async (context: { params: { id: string } }): Promise<object> => {
    const { id } = context.params;
    const apiUrlMerchant: string = `/api/v1/user/${id}`;
    const apiUrlPrices: string = `/api/v1/price/user/${id}`;

    try {
        const [merchant, price] = await Promise.all([
            fetch(apiUrlMerchant).then(handleErrors) as Promise<Merchant>,
            fetch(apiUrlPrices).then(handleErrors) as Promise<Price[]>,
        ]);

        return {
            merchant: merchant,
            prices: price,
            metaTitle: merchant?.name
        };
    } catch (error) {
        return {
            merchant: null,
            prices: [],
            metaTitle: 'Leider ist ein Fehler aufgetreten.',
        };
    }
};

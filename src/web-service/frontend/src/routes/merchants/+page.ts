import { handleErrors } from "../../assets/helper/handleErrors";

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

export const load = async (): Promise<object> => {
    const apiUrlMerchants = '/api/v1/user/role/1';

    try {
        const merchantsResponse: Response = await fetch(apiUrlMerchants);
        const merchants: Merchant[] = await handleErrors(merchantsResponse);
        const prices: Price[] = await getPricesArray(merchants);

        merchants.forEach((merchant) => {
            merchant.productsCount = calculateProductsCount(merchant.id, prices);
        });

        return {
            merchants,
            metaTitle: 'Auflistung der Supermärkte',
            headline: 'Alle verfügbaren Supermärkte',
        };
    } catch (error) {
        return {
            merchants: [],
            metaTitle: 'Auflistung der Supermärkte',
            headline: 'Alle verfügbaren Supermärkte',
        };
    }
};

async function fetchPrices(url: string): Promise<Price[]> {
    const response: Response = await fetch(url);

    return handleErrors(response);
}

async function getPricesArray(merchants: Merchant[]): Promise<Price[]> {
    const pricesPromises: Promise<Price[]>[] = merchants.map(
        (merchant: Merchant) => fetchPrices(`/api/v1/price/user/${merchant.id}`));
    const pricesArray: Price[][] = await Promise.all(pricesPromises);

    return pricesArray.flat();
}

function calculateProductsCount(merchantId: number, prices: Price[]): number {
    const merchantProducts: Price[] = prices.filter(
        (price: Price): boolean => price.userId === merchantId);

    return merchantProducts.length;
}

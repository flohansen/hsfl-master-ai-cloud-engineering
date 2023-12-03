import {handleErrors} from "../../assets/helper/handleErrors";

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
    const apiUrlUsers: string = '/api/v1/user/role/1';
    const apiUrlPrices: string = '/api/v1/price';

    const [merchants, prices] = await Promise.all([
        fetch(apiUrlUsers).then(handleErrors),
        fetch(apiUrlPrices).then(handleErrors),
    ]);

    if (merchants === null || prices === null) {
        console.warn('Products or users not found');
    }

    return {
        merchants: merchants ? getMerchantsContent(merchants, prices) : [],
        metaTitle: 'Auflistung der Supermärkte',
        headline: 'Alle verfügbaren Supermärkte',
    };
};

function getMerchantsContent(merchants: Merchant[], prices: Price[]): Merchant[] {
    merchants.forEach(merchant => {
        merchant.productsCount = calculateProductsCount(merchant.id, prices);
    });

    return merchants;
}

function calculateProductsCount(merchantId: number, prices: Price[]): number {
    const merchantProducts: Price[] = prices.filter(price => price.userId === merchantId);
    return merchantProducts.length;
}
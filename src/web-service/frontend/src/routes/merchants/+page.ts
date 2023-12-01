import {handleErrors} from "../../assets/helper/handleErrors";

interface User {
    id: number;
    name: string;
    email: string;
    role?: number;
    productsCount?: number;
}

interface Price {
    userId: number;
    productId: number;
    price: number;
}

export const load = async (): Promise<object> => {
    const apiUrlUsers: string = '/api/v1/user';
    const apiUrlPrices: string = '/api/v1/price';

    const [users, prices] = await Promise.all([
        fetch(apiUrlUsers).then(handleErrors),
        fetch(apiUrlPrices).then(handleErrors),
    ]);

    if (users === null || prices === null) {
        console.warn('Products or users not found');
    }

    return {
        merchants: users ? getMerchantsContent(users, prices) : [],
        metaTitle: 'Auflistung der Supermärkte',
        headline: 'Alle verfügbaren Supermärkte',
    };
};

function getMerchantsContent(users: User[], prices: Price[]): User[] {
    const merchants: User[] = filterUserByRole(users);

    merchants.forEach(merchant => {
        merchant.productsCount = calculateProductsCount(merchant.id, prices);
    });

    return merchants;

}

function calculateProductsCount(merchantId: number, prices: Price[]): number {
    const merchantProducts: Price[] = prices.filter(price => price.userId === merchantId);
    return merchantProducts.length;
}

function filterUserByRole(users: User[]): User[] {
    return users.filter(user => user.role === 1);
}
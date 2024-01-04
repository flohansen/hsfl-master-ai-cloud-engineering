import { handleErrors } from "../../assets/helper/handleErrors";
import { isAuthenticated } from "../../store";

interface Merchant {
    id: number;
    name: string;
    role?: number;
    productsCount?: number;
}

export const load = async (): Promise<Promise<object> | undefined> => {
    if (! isAuthenticated) return;

    const token: string | null = sessionStorage.getItem('access_token');
    const apiUrlMerchants: string = "/api/v1/user/role/1";

    if (! token) return;

    const requestOptions: object = {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${token}`
        },
    };

    try {
        const merchantsResponse: Response = await fetch(apiUrlMerchants, requestOptions);

        const merchants: Merchant[] = await handleErrors(merchantsResponse);

        for (const merchant of merchants) {
            merchant.productsCount = await calculateProductsCount(merchant.id);
        }

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

async function calculateProductsCount(merchantId: number): Promise<number> {
    const token: string | null = sessionStorage.getItem('access_token');

    if (! token || ! merchantId) return 0;

    const apiUrl: string = `/api/v1/price/user/${merchantId}`;
    const requestOptions: object = {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${token}`
        },
    };

    try {
        const response: Response = await fetch(apiUrl, requestOptions);
        const data: any = await handleErrors(response);
        return data !== null && data.length !== 0 ? data.length : 0;
    } catch (error) {
        return 0;
    }
}

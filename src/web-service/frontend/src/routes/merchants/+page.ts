import type { PageLoad } from './$types';
import { handleErrors } from "../../assets/helper/handleErrors";
import { fetchHelper } from "../../assets/helper/fetchHelper";
import { checkAuthentication } from "../../assets/helper/checkAuthentication";

export const load: PageLoad = async () : Promise<object> => {
    await checkAuthentication();

    const apiUrl: string = "/api/v1/user/role/1";
    const merchants: any = await fetchHelper(apiUrl);

    if (! merchants) return data;

    for (const merchant of merchants) {
        merchant.productsCount = await calculateProductsCount(merchant.id);
    }

    return data(merchants);
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

const data = (merchants: object[] = []): object => {
    return {
        merchants,
        metaTitle: 'Auflistung der Supermärkte',
        headline: 'Alle verfügbaren Supermärkte',
    };
};
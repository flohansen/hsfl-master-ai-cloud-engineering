import { handleErrors } from "../../assets/helper/handleErrors";

interface Merchant {
    id: number;
    name: string;
    role?: number;
    productsCount?: number;
}

export const load = async (): Promise<object> => {
    const apiUrlMerchants: string = "/api/v1/user/role/1";

    try {
        const merchantsResponse: Response = await fetch(apiUrlMerchants);
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
    const apiUrl: string = `/api/v1/price/user/${merchantId}`;

    try {
        const response: Response = await fetch(apiUrl);
        const data: any = await handleErrors(response);

        return data !== null && data.length !== 0 ? data.length : 0;
    } catch (error) {
        return 0;
    }
}

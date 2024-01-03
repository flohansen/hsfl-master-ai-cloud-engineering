import { handleErrors } from '../../../assets/helper/handleErrors';

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

export const load = async (context: { params: { id: string } }): Promise<object> => {
    const { id } = context.params;
    const apiUrlList: string = `/api/v1/shoppinglist/${id}/2`;
    const apiUrlEntries: string = `/api/v1/shoppinglistentries/${id}`;

    try {
        const [list, entries] = await Promise.all([
            fetch(apiUrlList).then(handleErrors) as Promise<List>,
            fetch(apiUrlEntries).then(handleErrors) as Promise<Entry[]>,
        ]);

        const uniqueProductIds: number[] = Array.from(new Set(entries.map(entry => entry.productId)));

        // Fetch products for each unique product ID
        const productsPromises: Promise<Product>[] = uniqueProductIds.map(productId =>
            fetch(`/api/v1/product/${productId}`).then(handleErrors) as Promise<Product>
        );

        const products: Product[] = await Promise.all(productsPromises);

        // Sort products by description
        const sortedProducts :Product[] = products.sort(
            (a: Product, b: Product) => a.description.localeCompare(b.description));

        return {
            list: list,
            entries: entries,
            products: sortedProducts,
            metaTitle: 'Liste: ' + list?.description,
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
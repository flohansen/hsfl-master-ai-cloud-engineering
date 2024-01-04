import { handleErrors } from "./handleErrors";

interface Product {
    id: number,
    description: string,
    ean: number,
}

export async function sortProducts(items: any): Promise<Product[]> {
    if (! items) return [];

    const uniqueProductIds: number[] = Array.from(new Set(items.map((item: { productId: number }) => item.productId)));
    const productsPromises: Promise<Product>[] = uniqueProductIds.map(productId =>
        fetch(`/api/v1/product/${productId}`).then(handleErrors) as Promise<Product>
    );

    const products: Product[] = await Promise.all(productsPromises);
    return products.sort(
        (a: Product, b: Product) => a.description.localeCompare(b.description));
}
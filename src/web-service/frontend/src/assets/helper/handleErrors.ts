/**
 * Handles errors in a fetch response and checks for correct response header
 */
export async function handleErrors(response: Response) : Promise<any> {
    if (! response.ok) {
        throw new Error(response.statusText);
    }

    const contentType: string | null = response.headers.get('content-type');
    return contentType && contentType.includes('application/json') ? response.json() : null;
}
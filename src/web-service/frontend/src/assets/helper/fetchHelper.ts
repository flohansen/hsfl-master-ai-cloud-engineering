import { handleErrors } from "./handleErrors";

export async function fetchHelper(apiUrl: string): Promise<object[]> {
    const userId: string | null = sessionStorage.getItem('user_id');
    const token: string | null = sessionStorage.getItem('access_token');
    let content: object[] = [];

    if (! token || ! userId) return content;

    const requestOptions: object = {
        headers: { 'Authorization': `Bearer ${token}` }
    };

    try {
        const response: Response = await fetch(apiUrl, requestOptions);
        return await handleErrors(response);
    } catch (error) {
        console.log("No data fetched");
    }

    return content;
}

import { decodeToken } from "./decodeToken";
import { goto } from '$app/navigation';

export async function checkAuthentication(): Promise<void> {
    const token: string | null = sessionStorage.getItem('access_token');
    const userId: string | null = sessionStorage.getItem('user_id');

    if (! userId || ! token || isExpired(token)) {
        await goto('/profile/login');
    }
}

function isExpired(token: string): boolean {
    const decodedToken: any = decodeToken(token);
    const currentTimestamp: number = Math.floor(Date.now() / 1000);

    return decodedToken.exp < currentTimestamp;
}
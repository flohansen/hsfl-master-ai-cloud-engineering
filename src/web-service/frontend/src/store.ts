import {type Writable, writable} from 'svelte/store';

export const isAuthenticated: Writable<boolean> = writable<boolean>(false);

export function setAuthenticationStatus(status: boolean): void {
    isAuthenticated.set(status);
}

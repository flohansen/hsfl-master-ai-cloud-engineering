export let isAuthenticated: boolean = false;

export function setAuthenticationStatus(status: boolean) {
    isAuthenticated = status;
}

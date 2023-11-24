/**
 * Handles click events and dispatches a custom 'click_outside' event
 * if the target is outside the provided HTML element.
 * It returns an object with a destroy method to remove the event listener.
 */
export function clickOutside(node: HTMLElement): {destroy(): void} {
    const handleClick = (event: MouseEvent): void => {
        if (node && ! node.contains(event.target as Node) && ! event.defaultPrevented) {
            node.dispatchEvent(new CustomEvent('click_outside', { detail: node }));
        }
    };

    document.addEventListener('click', handleClick, true);

    return {
        destroy(): void {
            document.removeEventListener('click', handleClick, true);
        },
    };
}
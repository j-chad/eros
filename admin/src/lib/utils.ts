export function debounce<T extends unknown[]>(fn: (...args: T) => void, delay = 150) {
    let timer: number;

    return (...args: T) => {
        clearTimeout(timer);
        timer = window.setTimeout(() => fn(...args), delay);
    };
}
import { writable } from 'svelte/store';

export const SendMessage = writable<((message: string) => void) | null>(null);
export const SendAction = writable<((message: object) => void) | null>(null);

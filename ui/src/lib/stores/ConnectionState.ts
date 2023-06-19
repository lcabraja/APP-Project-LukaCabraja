import { writable } from 'svelte/store';

export type ConnectionState = {
	status: string;
	initials: string;
	fill: string;
	background: string;
};

export const Connection = writable<ConnectionState>({
	status: 'uninitialized',
	initials: 'UN',
	fill: 'fill-token',
	background: 'bg-gray-400'
});

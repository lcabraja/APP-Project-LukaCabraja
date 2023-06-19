import { dev } from '$app/environment';
import path from 'path';

function getBackendURL(subpath: string) {
	path.join('/', subpath);
	if (dev) {
		return 'http:
	} else {
		return '/';
	}
}

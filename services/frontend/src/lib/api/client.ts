export const API_BASE = 'http://localhost:8080/api/v1';

export async function fetchWithAuth(endpoint: string, options: RequestInit = {}) {
	const token = typeof window !== 'undefined' ? localStorage.getItem('findex_token') : null;
	const headers = new Headers(options.headers || {});
	
	if (token) {
		headers.set('Authorization', `Bearer ${token}`);
	}

	const res = await fetch(`${API_BASE}${endpoint}`, {
		...options,
		headers
	});

	const json = await res.json();
	if (!res.ok) {
		throw new Error(json.error || 'Terjadi kesalahan pada server');
	}

	return json;
}

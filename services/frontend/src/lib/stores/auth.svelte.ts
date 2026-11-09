// Svelte 5 Runes Auth Store
export interface UserProfile {
	id: string;
	company_id: string;
	name: string;
	email: string;
	role: 'SUPERADMIN' | 'HR_ADMIN' | 'MANAGER' | 'EMPLOYEE';
	allow_wfh: boolean;
}

class AuthStore {
	token = $state<string | null>(null);
	user = $state<UserProfile | null>(null);

	constructor() {
		if (typeof window !== 'undefined') {
			const savedToken = localStorage.getItem('findex_token');
			const savedUser = localStorage.getItem('findex_user');
			if (savedToken && savedUser) {
				try {
					this.token = savedToken;
					this.user = JSON.parse(savedUser);
				} catch (e) {
					this.logout();
				}
			}
		}
	}

	setSession(token: string, user: UserProfile) {
		this.token = token;
		this.user = user;
		if (typeof window !== 'undefined') {
			localStorage.setItem('findex_token', token);
			localStorage.setItem('findex_user', JSON.stringify(user));
		}
	}

	logout() {
		this.token = null;
		this.user = null;
		if (typeof window !== 'undefined') {
			localStorage.removeItem('findex_token');
			localStorage.removeItem('findex_user');
		}
	}

	get isAuthenticated() {
		return !!this.token;
	}

	get isAdmin() {
		return this.user?.role === 'SUPERADMIN' || this.user?.role === 'HR_ADMIN';
	}
}

export const auth = new AuthStore();

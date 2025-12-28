const STORAGE_KEY = 'admin_api_key';

class AuthState {
    apiKey = $state(localStorage.getItem(STORAGE_KEY));

    login(key: string) {
        this.apiKey = key;
        localStorage.setItem(STORAGE_KEY, key);
    }

    logout() {
        this.apiKey = null;
        localStorage.removeItem(STORAGE_KEY);
    }

    get isAuthenticated() {
        return this.apiKey !== null && this.apiKey !== '';
    }
}

export const auth = new AuthState();
import config from '../config/api.js';

class AuthService {
  constructor() {
    this.baseURL = config.API_BASE_URL;
    this.timeout = config.TIMEOUT;
  }

  async apiCall(endpoint, options = {}, retryCount = 0) {
    const url = `${this.baseURL}${endpoint}`;
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), this.timeout);

    const defaultHeaders = {
      'Content-Type': 'application/json',
    };

    const fetchConfig = {
      headers: {
        ...defaultHeaders,
        ...options.headers,
      },
      signal: controller.signal,
      credentials: 'include',
      ...options,
    };

    try {
      const response = await fetch(url, fetchConfig);
      clearTimeout(timeoutId);

      if (response.status === config.HTTP_STATUS.UNAUTHORIZED) {
        this.clearAuthData();
        throw new Error('Authentication failed. Please login again.');
      }

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.message || `HTTP ${response.status}: ${response.statusText}`);
      }

      return {
        success: true,
        data,
        status: response.status,
      };

    } catch (error) {
      clearTimeout(timeoutId);

      if (error.name === 'AbortError') {
        throw new Error('Request timeout. Please try again.');
      }

      if (error.name === 'TypeError' && error.message.includes('fetch')) {
        if (retryCount < config.RETRY_ATTEMPTS) {
          console.warn(`Retrying request to ${endpoint} (${retryCount + 1}/${config.RETRY_ATTEMPTS})`);
          await new Promise(resolve => setTimeout(resolve, 1000 * (retryCount + 1)));
          return this.apiCall(endpoint, options, retryCount + 1);
        }
        throw new Error('Unable to connect to server. Please check your internet connection.');
      }

      throw error;
    }
  }

  async login(email, password) {
    try {
      const result = await this.apiCall(config.AUTH_ENDPOINTS.LOGIN, {
        method: 'POST',
        body: JSON.stringify({ email, password }),
      });

      if (result.success) {
        const { user } = result.data;
        if (user) {
          localStorage.setItem(config.STORAGE_KEYS.USER, JSON.stringify(user));
        }
        return { success: true, user };
      }

      return result;
    } catch (error) {
      return { success: false, error: error.message || 'Login failed' };
    }
  }

  async register(email, password) {
    try {
      const result = await this.apiCall(config.AUTH_ENDPOINTS.REGISTER, {
        method: 'POST',
        body: JSON.stringify({ email, password }),
      });

      if (result.success) {
        const { user } = result.data;
        if (user) {
          localStorage.setItem(config.STORAGE_KEYS.USER, JSON.stringify(user));
        }
        return { success: true, user };
      }

      return result;
    } catch (error) {
      return { success: false, error: error.message || 'Registration failed' };
    }
  }

  async logout() {
    this.clearAuthData();
  }


  isAuthenticated() {
    return !!localStorage.getItem(config.STORAGE_KEYS.USER);
  }

  clearAuthData() {
    localStorage.removeItem(config.STORAGE_KEYS.USER);
  }
}

export default new AuthService();
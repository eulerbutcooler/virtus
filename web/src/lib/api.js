const BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

export const api = {
  async fetch(endpoint, options = {}) {
    const token = localStorage.getItem('virtus_token');

    const headers = {
      'Content-Type': 'application/json',
      ...options.headers,
    };

    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    let response = await fetch(`${BASE_URL}${endpoint}`, {
      ...options,
      headers,
    });

    if (response.status === 401 && token) {
      try {
        const refreshToken = localStorage.getItem('virtus_refresh_token');
        if (!refreshToken) throw new Error('no refresh token');
        const refreshResponse = await fetch(`${BASE_URL}/auth/refresh`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ refresh_token: refreshToken }),
        });
        if (refreshResponse.ok) {
          const json = await refreshResponse.json();
          const data = json.data ?? json;
          localStorage.setItem('virtus_token', data.access_token);
          localStorage.setItem('virtus_refresh_token', data.refresh_token);
          headers['Authorization'] = `Bearer ${data.access_token}`;
          response = await fetch(`${BASE_URL}${endpoint}`, {
            ...options,
            headers,
          });
        } else {
          localStorage.removeItem('virtus_token');
          localStorage.removeItem('virtus_refresh_token');
          window.location.href = '/login';
        }
      } catch (err) {
        localStorage.removeItem('virtus_token');
        localStorage.removeItem('virtus_refresh_token');
        window.location.href = '/login';
      }
    }

    if (!response.ok) {
      const errorData = await response.json().catch(() => null);
      throw new Error(errorData?.error?.message || response.statusText);
    }

    if (response.status === 204) return null;
    const json = await response.json();
    return json.data ?? json; // unwrap envelope { data: ... }
  },

  get(endpoint, options) { return this.fetch(endpoint, { ...options, method: 'GET' }) },
  post(endpoint, body, options) { return this.fetch(endpoint, { ...options, method: 'POST', body: JSON.stringify(body) }) },
  patch(endpoint, body, options) { return this.fetch(endpoint, { ...options, method: 'PATCH', body: JSON.stringify(body) }) },
  delete(endpoint, options) { return this.fetch(endpoint, { ...options, method: 'DELETE' }) },
};

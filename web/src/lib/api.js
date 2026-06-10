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
        const refreshResponse = await fetch(`${BASE_URL}/auth/refresh`, { method: 'POST' });
        if (refreshResponse.ok) {
          const data = await refreshResponse.json();
          localStorage.setItem('virtus_token', data.access_token);
          headers['Authorization'] = `Bearer ${data.access_token}`;
          response = await fetch(`${BASE_URL}${endpoint}`, {
            ...options,
            headers,
          });
        } else {
          localStorage.removeItem('virtus_token');
          window.location.href = '/login';
        }
      } catch (err) {
        localStorage.removeItem('virtus_token');
        window.location.href = '/login';
      }
    }

    if (!response.ok) {
      const errorData = await response.json().catch(() => null);
      throw new Error(errorData?.message || response.statusText);
    }

    if (response.status === 204) return null;
    return response.json();
  },

  get(endpoint, options) { return this.fetch(endpoint, { ...options, method: 'GET' }) },
  post(endpoint, body, options) { return this.fetch(endpoint, { ...options, method: 'POST', body: JSON.stringify(body) }) },
  patch(endpoint, body, options) { return this.fetch(endpoint, { ...options, method: 'PATCH', body: JSON.stringify(body) }) },
  delete(endpoint, options) { return this.fetch(endpoint, { ...options, method: 'DELETE' }) },
};

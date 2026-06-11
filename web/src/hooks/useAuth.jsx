import { createContext, useContext, useState, useEffect, useCallback } from 'react';
import { api } from '../lib/api';

const AuthContext = createContext(null);

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const initAuth = async () => {
      const token = localStorage.getItem('virtus_token');
      if (token) {
        try {
          const userData = await api.get('/me');
          setUser(userData);
        } catch (error) {
          console.error('Failed to fetch user', error);
          localStorage.removeItem('virtus_token');
        }
      }
      setLoading(false);
    };
    initAuth();
  }, []);

  const login = async (email, password) => {
    setLoading(true);
    try {
      const data = await api.post('/auth/login', { email, password });

      localStorage.setItem('virtus_token', data.token.access_token);
      localStorage.setItem('virtus_refresh_token', data.token.refresh_token);
      setUser(data.user);
    } finally {
      setLoading(false);
    }
  };

  const register = async (name, email, password) => {
    setLoading(true);
    try {
      const data = await api.post('/auth/register', { name, email, password });

      localStorage.setItem('virtus_token', data.token.access_token);
      localStorage.setItem('virtus_refresh_token', data.token.refresh_token);
      setUser(data.user);
    } finally {
      setLoading(false);
    }
  };

  const logout = useCallback(() => {
    localStorage.removeItem('virtus_token');
    localStorage.removeItem('virtus_refresh_token');
    setUser(null);
  }, []);

  const value = {
    user,
    isAuthenticated: !!user,
    loading,
    login,
    register,
    logout
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  return useContext(AuthContext);
}

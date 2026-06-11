import { useState, useEffect, useCallback } from 'react';
import { api } from '../lib/api';

export function useAdminUsers(limit = 50, offset = 0) {
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const fetchUsers = useCallback(async () => {
    try {
      setLoading(true);
      const data = await api.get(`/admin/users?limit=${limit}&offset=${offset}`);
      setUsers(Array.isArray(data) ? data : (data.items || []));
    } catch (err) {
      setError(err);
    } finally {
      setLoading(false);
    }
  }, [limit, offset]);

  useEffect(() => {
    fetchUsers();
  }, [fetchUsers]);

  const verifyUser = async (id) => {
    await api.post(`/admin/users/${id}/verify`);
    await fetchUsers();
  };

  const deleteUser = async (id) => {
    await api.delete(`/admin/users/${id}`);
    await fetchUsers();
  };

  return { users, loading, error, verifyUser, deleteUser, refresh: fetchUsers };
}

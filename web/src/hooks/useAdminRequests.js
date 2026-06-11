import { useState, useEffect, useCallback } from 'react';
import { api } from '../lib/api';

export function useAdminRequests(limit = 50, offset = 0) {
  const [requests, setRequests] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const fetchRequests = useCallback(async () => {
    try {
      setLoading(true);
      const data = await api.get(`/admin/requests?limit=${limit}&offset=${offset}`);
      setRequests(Array.isArray(data) ? data : (data.items || []));
    } catch (err) {
      setError(err);
    } finally {
      setLoading(false);
    }
  }, [limit, offset]);

  useEffect(() => {
    fetchRequests();
  }, [fetchRequests]);

  const verifyRequest = async (id) => {
    await api.post(`/admin/requests/${id}/verify`);
    await fetchRequests();
  };

  const rejectRequest = async (id, note) => {
    await api.post(`/admin/requests/${id}/reject`, { note });
    await fetchRequests();
  };

  return { requests, loading, error, verifyRequest, rejectRequest, refresh: fetchRequests };
}

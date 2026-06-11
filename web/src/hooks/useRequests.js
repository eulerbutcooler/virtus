import { useState, useEffect, useCallback } from 'react';
import { api } from '../lib/api';

export function useRequests(limit = null) {
  const [requests, setRequests] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const fetchRequests = useCallback(async () => {
    try {
      setLoading(true);
      const data = await api.get(`/requests${limit ? `?limit=${limit}` : ''}`);
      setRequests(Array.isArray(data) ? data : (data.items || []));
    } catch (err) {
      setError(err);
    } finally {
      setLoading(false);
    }
  }, [limit]);

  useEffect(() => {
    fetchRequests();
  }, [fetchRequests]);

  const createRequest = async (payload) => {
    const data = await api.post('/requests', payload);
    await fetchRequests();
    return data;
  };

  const updateRequest = async (id, payload) => {
    const data = await api.patch(`/requests/${id}`, payload);
    await fetchRequests();
    return data;
  };

  const deleteRequest = async (id) => {
    await api.delete(`/requests/${id}`);
    await fetchRequests();
  };

  return { requests, loading, error, createRequest, updateRequest, deleteRequest, refresh: fetchRequests };
}

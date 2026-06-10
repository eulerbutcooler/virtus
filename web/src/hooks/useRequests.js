import { useState, useEffect } from 'react';
import { api } from '../lib/api';

export function useRequests(limit = null) {
  const [requests, setRequests] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    async function fetchData() {
      try {
        setLoading(true);
        const data = await api.get(`/requests${limit ? `?limit=${limit}` : ''}`);
        setRequests(data.items || []);
      } catch (err) {
        setError(err);
      } finally {
        setLoading(false);
      }
    }
    fetchData();
  }, [limit]);

  return { requests, loading, error };
}

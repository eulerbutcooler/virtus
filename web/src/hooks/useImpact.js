import { useState, useEffect, useCallback } from 'react';
import { api } from '../lib/api';

export function useImpact() {
  const [impacts, setImpacts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const fetchImpacts = useCallback(async () => {
    try {
      setLoading(true);
      const data = await api.get('/impact');
      setImpacts(Array.isArray(data) ? data : (data.items || []));
    } catch (err) {
      setError(err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchImpacts();
  }, [fetchImpacts]);

  const createImpact = async (payload) => {

    await api.post('/impact', payload);
    await fetchImpacts();
  };

  return { impacts, loading, error, createImpact, refresh: fetchImpacts };
}

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
      setImpacts(data.items || []);
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
    // payload: { delivery_id, interval, outcome, satisfaction_score, description }
    await api.post('/impact', payload);
    await fetchImpacts();
  };

  return { impacts, loading, error, createImpact, refresh: fetchImpacts };
}

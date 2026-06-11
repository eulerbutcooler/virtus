import { useState, useEffect, useCallback } from 'react';
import { api } from '../lib/api';

export function useContributions(limit = null) {
  const [contributions, setContributions] = useState([]);
  const [total, setTotal] = useState({ amount: 0, count: 0 });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const fetchContributions = useCallback(async () => {
    try {
      setLoading(true);
      const [totalData, listData] = await Promise.all([
        api.get('/contributions/total'),
        api.get(`/contributions${limit ? `?limit=${limit}` : ''}`)
      ]);
      setTotal(totalData);
      setContributions(Array.isArray(listData) ? listData : (listData.items || []));
    } catch (err) {
      setError(err);
    } finally {
      setLoading(false);
    }
  }, [limit]);

  useEffect(() => {
    fetchContributions();
  }, [fetchContributions]);

  const createContribution = async (amount, currency = 'USD') => {
    const data = await api.post('/contributions', { amount, currency });
    await fetchContributions();
    return data; // contains client_secret
  };

  const optimisticUpdate = (tempContribution) => {
    setContributions(prev => [tempContribution, ...prev]);
    setTotal(prev => ({ amount: prev.amount + tempContribution.amount, count: prev.count + 1 }));
  };

  return { contributions, total, loading, error, createContribution, refresh: fetchContributions, optimisticUpdate };
}

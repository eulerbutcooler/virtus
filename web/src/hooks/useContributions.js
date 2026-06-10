import { useState, useEffect } from 'react';
import { api } from '../lib/api';

export function useContributions(limit = null) {
  const [contributions, setContributions] = useState([]);
  const [total, setTotal] = useState({ amount: 0, count: 0 });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    async function fetchData() {
      try {
        setLoading(true);
        const [totalData, listData] = await Promise.all([
          api.get('/contributions/total'),
          api.get(`/contributions${limit ? `?limit=${limit}` : ''}`)
        ]);
        setTotal(totalData);
        setContributions(listData.items || []);
      } catch (err) {
        setError(err);
      } finally {
        setLoading(false);
      }
    }
    fetchData();
  }, [limit]);

  return { contributions, total, loading, error };
}

import { useState, useEffect } from 'react';
import { api } from '../lib/api';

export function usePool() {
  const [pool, setPool] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    async function fetchPool() {
      try {
        const data = await api.get('/pool');
        setPool(data);
      } catch (err) {
        setError(err);
      } finally {
        setLoading(false);
      }
    }
    fetchPool();
  }, []);

  return { pool, loading, error };
}

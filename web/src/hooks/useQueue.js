import { useState, useEffect } from 'react';
import { api } from '../lib/api';

export function useQueue(requestId) {
  const [queueEntry, setQueueEntry] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    async function fetchQueue() {
      if (!requestId) {
        setLoading(false);
        return;
      }
      try {
        setLoading(true);
        const data = await api.get(`/queue/${requestId}`);
        setQueueEntry(data);
      } catch (err) {
        setError(err);
      } finally {
        setLoading(false);
      }
    }
    fetchQueue();
  }, [requestId]);

  return { queueEntry, loading, error };
}

export function useFullQueue(limit = 20, offset = 0) {
  const [queue, setQueue] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    async function fetchFullQueue() {
      try {
        setLoading(true);
        const data = await api.get(`/queue?limit=${limit}&offset=${offset}`);
        setQueue(Array.isArray(data) ? data : (data.items || []));
      } catch (err) {
        setError(err);
      } finally {
        setLoading(false);
      }
    }
    fetchFullQueue();
  }, [limit, offset]);

  return { queue, loading, error };
}

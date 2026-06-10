import { useState, useEffect } from 'react';
import { api } from '../lib/api';

export function useQueue(requestId) {
  const [queueEntry, setQueueEntry] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    async function fetchQueue() {
      // If we don't have a requestId, we might just fetch the user's active request from somewhere else first
      // But for now, if requestId is undefined, we could skip or fetch a generic /queue/me
      // The API design says /queue/{requestID}. Let's assume we fetch the user's active request first if needed.
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

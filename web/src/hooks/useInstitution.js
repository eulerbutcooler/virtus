import { useState, useEffect, useCallback } from 'react';
import { api } from '../lib/api';

export function useInstitution() {
  const [institution, setInstitution] = useState(null);
  const [contributions, setContributions] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const fetchInstitutionData = useCallback(async () => {
    try {
      setLoading(true);
      const meData = await api.get('/institutions/me');
      setInstitution(meData);
      if (meData?.id) {
        const contribData = await api.get(`/institutions/${meData.id}/contributions`);
        setContributions(contribData.items || []);
      }
    } catch (err) {
      setError(err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchInstitutionData();
  }, [fetchInstitutionData]);

  const updateInstitution = async (id, payload) => {
    await api.patch(`/institutions/${id}`, payload);
    await fetchInstitutionData();
  };

  const createContribution = async (id, payload) => {
    // payload: { amount, currency, category_tags, region_tags }
    const data = await api.post(`/institutions/${id}/contributions`, payload);
    await fetchInstitutionData();
    return data; // returns client_secret for Stripe
  };

  return { institution, contributions, loading, error, updateInstitution, createContribution, refresh: fetchInstitutionData };
}

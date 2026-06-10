import { useState, useEffect, useCallback } from 'react';
import { api } from '../lib/api';

export function useAdminFulfillments(limit = 50, offset = 0) {
  const [fulfillments, setFulfillments] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const fetchFulfillments = useCallback(async () => {
    try {
      setLoading(true);
      const data = await api.get(`/admin/fulfillments?limit=${limit}&offset=${offset}`);
      setFulfillments(data.items || []);
    } catch (err) {
      setError(err);
    } finally {
      setLoading(false);
    }
  }, [limit, offset]);

  useEffect(() => {
    fetchFulfillments();
  }, [fetchFulfillments]);

  const beginFulfillment = async (requestId, payload) => {
    // payload: { vendor_name, vendor_ref, actual_cost, notes }
    await api.post(`/admin/fulfillments`, { request_id: requestId, ...payload });
    await fetchFulfillments();
  };

  const updateFulfillment = async (id, payload) => {
    await api.patch(`/admin/fulfillments/${id}`, payload);
    await fetchFulfillments();
  };

  const cancelFulfillment = async (id) => {
    await api.post(`/admin/fulfillments/${id}/cancel`);
    await fetchFulfillments();
  };

  const shipDelivery = async (fulfillmentId, trackingNumber, carrier) => {
    await api.post(`/admin/deliveries`, { fulfillment_id: fulfillmentId, tracking_number: trackingNumber, carrier });
    await fetchFulfillments();
  };

  const verifyDelivery = async (deliveryId, proofPhotoUrl, deliveredAt) => {
    await api.post(`/admin/deliveries/${deliveryId}/verify`, { proof_photo_url: proofPhotoUrl, delivered_at: deliveredAt });
    await fetchFulfillments();
  };

  const failDelivery = async (deliveryId) => {
    await api.post(`/admin/deliveries/${deliveryId}/fail`);
    await fetchFulfillments();
  };

  return { 
    fulfillments, loading, error, 
    beginFulfillment, updateFulfillment, cancelFulfillment,
    shipDelivery, verifyDelivery, failDelivery,
    refresh: fetchFulfillments 
  };
}

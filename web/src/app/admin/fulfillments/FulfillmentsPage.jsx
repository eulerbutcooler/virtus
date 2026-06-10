import { useState } from 'react';
import { useAdminFulfillments } from '../../../hooks/useAdminFulfillments';
import { useAdminRequests } from '../../../hooks/useAdminRequests';
import PageHeader from '../../../components/layout/PageHeader';
import Table from '../../../components/ui/Table';
import Badge from '../../../components/ui/Badge';
import Button from '../../../components/ui/Button';
import Modal from '../../../components/ui/Modal';
import Input from '../../../components/ui/Input';
import { formatCurrency, formatDate } from '../../../lib/formatters';

export default function FulfillmentsPage() {
  const { fulfillments, loading, beginFulfillment, cancelFulfillment, shipDelivery, verifyDelivery, failDelivery } = useAdminFulfillments();
  const { requests } = useAdminRequests(); // To select funded requests for fulfillment
  
  const [beginModalOpen, setBeginModalOpen] = useState(false);
  const [shipModalOpen, setShipModalOpen] = useState(false);
  const [verifyModalOpen, setVerifyModalOpen] = useState(false);
  const [selectedFulfillment, setSelectedFulfillment] = useState(null);
  
  const [actionLoading, setActionLoading] = useState(false);

  // Form states
  const [beginForm, setBeginForm] = useState({ request_id: '', vendor_name: '', vendor_ref: '', actual_cost: '', notes: '' });
  const [shipForm, setShipForm] = useState({ tracking_number: '', carrier: '' });
  const [verifyForm, setVerifyForm] = useState({ proof_photo_url: '', delivered_at: new Date().toISOString() });

  const fundedRequests = requests?.filter(r => r.status === 'funded') || [];

  const handleBegin = async (e) => {
    e.preventDefault();
    setActionLoading(true);
    try {
      const { request_id, actual_cost, ...rest } = beginForm;
      await beginFulfillment(request_id, { actual_cost: Number(actual_cost), ...rest });
      setBeginModalOpen(false);
      setBeginForm({ request_id: '', vendor_name: '', vendor_ref: '', actual_cost: '', notes: '' });
    } catch (err) {
      console.error(err);
    } finally {
      setActionLoading(false);
    }
  };

  const handleShip = async (e) => {
    e.preventDefault();
    if (!selectedFulfillment) return;
    setActionLoading(true);
    try {
      await shipDelivery(selectedFulfillment.id, shipForm.tracking_number, shipForm.carrier);
      setShipModalOpen(false);
      setSelectedFulfillment(null);
      setShipForm({ tracking_number: '', carrier: '' });
    } catch (err) {
      console.error(err);
    } finally {
      setActionLoading(false);
    }
  };

  const handleVerify = async (e) => {
    e.preventDefault();
    if (!selectedFulfillment) return;
    setActionLoading(true);
    try {
      // Find delivery ID from fulfillment if nested, assume the API routes take delivery_id or fulfillment_id.
      // Wait, shipDelivery creates a delivery. The prompt says `/admin/deliveries/{id}/verify`.
      // Let's assume fulfillment returns its delivery object inside `fulfillment.delivery`.
      const deliveryId = selectedFulfillment.delivery?.id;
      if (deliveryId) {
        await verifyDelivery(deliveryId, verifyForm.proof_photo_url, verifyForm.delivered_at);
      }
      setVerifyModalOpen(false);
      setSelectedFulfillment(null);
      setVerifyForm({ proof_photo_url: '', delivered_at: new Date().toISOString() });
    } catch (err) {
      console.error(err);
    } finally {
      setActionLoading(false);
    }
  };

  const handleCancel = async (f) => {
    if (confirm('Are you sure you want to cancel this fulfillment?')) {
      setActionLoading(true);
      try {
        await cancelFulfillment(f.id);
      } catch (err) {
        console.error(err);
      } finally {
        setActionLoading(false);
      }
    }
  };

  const getStatusBadge = (status) => {
    const map = {
      processing: 'warning',
      purchased: 'info',
      shipped: 'info',
      delivered: 'success',
      failed: 'error',
      canceled: 'error'
    };
    return <Badge variant={map[status] || 'neutral'}>{status}</Badge>;
  };

  const columns = [
    { key: 'request', label: 'Request', render: f => <span style={{ color: 'var(--text-primary)' }}>{f.request?.item_name || 'Unknown'}</span> },
    { key: 'vendor', label: 'Vendor', render: f => f.vendor_name },
    { key: 'cost', label: 'Cost', render: f => formatCurrency(f.actual_cost) },
    { key: 'status', label: 'Status', render: f => getStatusBadge(f.status) },
    { key: 'date', label: 'Date', render: f => formatDate(f.created_at) },
    { 
      key: 'actions', 
      label: 'Actions', 
      render: f => (
        <div style={{ display: 'flex', gap: 'var(--space-2)' }}>
          {f.status === 'processing' || f.status === 'purchased' ? (
             <Button size="sm" variant="ghost" onClick={() => { setSelectedFulfillment(f); setShipModalOpen(true); }} disabled={actionLoading}>Ship</Button>
          ) : null}
          {f.status === 'shipped' && f.delivery ? (
             <Button size="sm" variant="ghost" style={{ color: 'var(--success)' }} onClick={() => { setSelectedFulfillment(f); setVerifyModalOpen(true); }} disabled={actionLoading}>Verify Delivery</Button>
          ) : null}
          {f.status !== 'delivered' && f.status !== 'canceled' && f.status !== 'failed' && (
             <Button size="sm" variant="ghost" style={{ color: 'var(--error)' }} onClick={() => handleCancel(f)} disabled={actionLoading}>Cancel</Button>
          )}
        </div>
      ) 
    }
  ];

  return (
    <div style={{ padding: 'var(--space-6)', maxWidth: '1200px', margin: '0 auto' }}>
      <PageHeader 
        title="Fulfillments & Deliveries" 
        subtitle="Manage purchasing and shipping of funded requests" 
        action={<Button variant="primary" onClick={() => setBeginModalOpen(true)}>+ Begin Fulfillment</Button>}
      />
      
      {loading ? (
        <div style={{ height: '300px', animation: 'shimmer 1.5s infinite', background: 'var(--bg-surface)', borderRadius: 'var(--radius-lg)' }}></div>
      ) : (
        <Table columns={columns} data={fulfillments} />
      )}

      {/* Begin Fulfillment Modal */}
      <Modal open={beginModalOpen} onClose={() => setBeginModalOpen(false)} title="Begin Fulfillment">
        <form onSubmit={handleBegin} style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-4)' }}>
          <select 
            value={beginForm.request_id} 
            onChange={e => setBeginForm(p => ({ ...p, request_id: e.target.value }))}
            required
            style={{ padding: 'var(--space-2)', background: 'var(--bg-surface)', border: '1px solid var(--border-default)', color: 'var(--text-primary)', borderRadius: 'var(--radius-sm)' }}
          >
            <option value="">Select a funded request...</option>
            {fundedRequests.map(r => <option key={r.id} value={r.id}>{r.item_name} ({formatCurrency(r.estimated_cost)})</option>)}
          </select>
          <Input placeholder="Vendor Name" value={beginForm.vendor_name} onChange={e => setBeginForm(p => ({ ...p, vendor_name: e.target.value }))} required />
          <Input placeholder="Vendor Order Reference" value={beginForm.vendor_ref} onChange={e => setBeginForm(p => ({ ...p, vendor_ref: e.target.value }))} required />
          <Input type="number" placeholder="Actual Cost ($)" value={beginForm.actual_cost} onChange={e => setBeginForm(p => ({ ...p, actual_cost: e.target.value }))} required min="1" />
          <textarea placeholder="Notes (optional)" value={beginForm.notes} onChange={e => setBeginForm(p => ({ ...p, notes: e.target.value }))} style={{ padding: 'var(--space-2)', background: 'var(--bg-surface)', border: '1px solid var(--border-default)', color: 'var(--text-primary)', borderRadius: 'var(--radius-sm)', minHeight: '80px', fontFamily: 'inherit' }} />
          <Button type="submit" variant="primary" disabled={actionLoading}>Submit</Button>
        </form>
      </Modal>

      {/* Ship Delivery Modal */}
      <Modal open={shipModalOpen} onClose={() => setShipModalOpen(false)} title="Ship Fulfillment">
        <form onSubmit={handleShip} style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-4)' }}>
          <p style={{ color: 'var(--text-secondary)' }}>Log shipping details for <strong>{selectedFulfillment?.request?.item_name}</strong>.</p>
          <Input placeholder="Tracking Number" value={shipForm.tracking_number} onChange={e => setShipForm(p => ({ ...p, tracking_number: e.target.value }))} required />
          <Input placeholder="Carrier (e.g. UPS, FedEx)" value={shipForm.carrier} onChange={e => setShipForm(p => ({ ...p, carrier: e.target.value }))} required />
          <Button type="submit" variant="primary" disabled={actionLoading}>Mark as Shipped</Button>
        </form>
      </Modal>

      {/* Verify Delivery Modal */}
      <Modal open={verifyModalOpen} onClose={() => setVerifyModalOpen(false)} title="Verify Delivery">
        <form onSubmit={handleVerify} style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-4)' }}>
          <p style={{ color: 'var(--text-secondary)' }}>Provide proof of delivery for <strong>{selectedFulfillment?.request?.item_name}</strong>.</p>
          <Input placeholder="Proof Photo URL" value={verifyForm.proof_photo_url} onChange={e => setVerifyForm(p => ({ ...p, proof_photo_url: e.target.value }))} required />
          <Input type="datetime-local" value={verifyForm.delivered_at.slice(0, 16)} onChange={e => setVerifyForm(p => ({ ...p, delivered_at: new Date(e.target.value).toISOString() }))} required />
          <Button type="submit" variant="primary" disabled={actionLoading}>Verify Delivery</Button>
        </form>
      </Modal>
    </div>
  );
}

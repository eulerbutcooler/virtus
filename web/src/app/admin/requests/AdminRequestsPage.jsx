import { useState } from 'react';
import { useAdminRequests } from '../../../hooks/useAdminRequests';
import PageHeader from '../../../components/layout/PageHeader';
import Table from '../../../components/ui/Table';
import Badge from '../../../components/ui/Badge';
import Button from '../../../components/ui/Button';
import Modal from '../../../components/ui/Modal';
import { formatCurrency, formatDate } from '../../../lib/formatters';

export default function AdminRequestsPage() {
  const { requests, loading, verifyRequest, rejectRequest } = useAdminRequests();
  const [selectedRequest, setSelectedRequest] = useState(null);
  const [rejectModalOpen, setRejectModalOpen] = useState(false);
  const [rejectNote, setRejectNote] = useState('');
  const [actionLoading, setActionLoading] = useState(false);

  const handleVerify = async (request) => {
    setActionLoading(true);
    try {
      await verifyRequest(request.id);
    } catch (err) {
      console.error(err);
    } finally {
      setActionLoading(false);
    }
  };

  const handleReject = async () => {
    if (!selectedRequest || !rejectNote.trim()) return;
    setActionLoading(true);
    try {
      await rejectRequest(selectedRequest.id, rejectNote);
      setRejectModalOpen(false);
      setSelectedRequest(null);
      setRejectNote('');
    } catch (err) {
      console.error(err);
    } finally {
      setActionLoading(false);
    }
  };

  const getStatusBadge = (status) => {
    const map = {
      pending: 'warning',
      queued: 'warning',
      funded: 'success',
      delivered: 'success',
      completed: 'success',
      rejected: 'error'
    };
    return <Badge variant={map[status] || 'neutral'}>{status}</Badge>;
  };

  const columns = [
    { key: 'user_id', label: 'User ID', render: r => <span style={{ fontFamily: 'monospace', color: 'var(--text-secondary)' }}>{r.user_id.substring(0, 8)}</span> },
    { key: 'item_name', label: 'Item Name' },
    { key: 'status', label: 'Status', render: r => getStatusBadge(r.status) },
    { key: 'estimated_cost', label: 'Est. Cost', render: r => formatCurrency(r.estimated_cost) },
    { key: 'created_at', label: 'Date', render: r => formatDate(r.created_at) },
    { 
      key: 'actions', 
      label: 'Actions', 
      render: r => (
        <div style={{ display: 'flex', gap: 'var(--space-2)' }}>
          {r.status === 'pending' && (
            <>
              <Button size="sm" variant="ghost" style={{ color: 'var(--success)' }} onClick={() => handleVerify(r)} disabled={actionLoading}>Verify</Button>
              <Button size="sm" variant="ghost" style={{ color: 'var(--error)' }} onClick={() => { setSelectedRequest(r); setRejectModalOpen(true); }} disabled={actionLoading}>Reject</Button>
            </>
          )}
        </div>
      ) 
    }
  ];

  return (
    <div style={{ padding: 'var(--space-6)', maxWidth: '1200px', margin: '0 auto' }}>
      <PageHeader title="Requests Management" subtitle="Verify or reject incoming member requests" />
      
      {loading ? (
        <div style={{ height: '300px', animation: 'shimmer 1.5s infinite', background: 'var(--bg-surface)', borderRadius: 'var(--radius-lg)' }}></div>
      ) : (
        <Table columns={columns} data={requests} />
      )}

      <Modal open={rejectModalOpen} onClose={() => setRejectModalOpen(false)} title="Reject Request">
        <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-4)' }}>
          <p style={{ color: 'var(--text-secondary)' }}>
            Please provide a reason for rejecting the request for <strong>{selectedRequest?.item_name}</strong>. This note will be visible to the member.
          </p>
          <textarea
            value={rejectNote}
            onChange={(e) => setRejectNote(e.target.value)}
            placeholder="Rejection reason..."
            required
            style={{
              background: 'var(--bg-surface)',
              border: '1px solid var(--border-default)',
              borderRadius: 'var(--radius-sm)',
              color: 'var(--text-primary)',
              padding: 'var(--space-3)',
              minHeight: '100px',
              fontFamily: 'inherit'
            }}
          />
          <div style={{ display: 'flex', justifyContent: 'flex-end', gap: 'var(--space-3)', marginTop: 'var(--space-4)' }}>
            <Button variant="ghost" onClick={() => setRejectModalOpen(false)} disabled={actionLoading}>Cancel</Button>
            <Button variant="danger" onClick={handleReject} disabled={actionLoading || !rejectNote.trim()}>
              {actionLoading ? 'Processing...' : 'Reject Request'}
            </Button>
          </div>
        </div>
      </Modal>
    </div>
  );
}

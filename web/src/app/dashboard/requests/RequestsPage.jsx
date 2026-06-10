import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useRequests } from '../../../hooks/useRequests';
import PageHeader from '../../../components/layout/PageHeader';
import Table from '../../../components/ui/Table';
import Badge from '../../../components/ui/Badge';
import Button from '../../../components/ui/Button';
import Modal from '../../../components/ui/Modal';
import RequestForm from '../../../components/forms/RequestForm';
import EmptyState from '../../../components/ui/EmptyState';
import { formatCurrency, formatDate } from '../../../lib/formatters';

export default function RequestsPage() {
  const { requests, loading, deleteRequest, updateRequest } = useRequests();
  const [deleteModalOpen, setDeleteModalOpen] = useState(false);
  const [editModalOpen, setEditModalOpen] = useState(false);
  const [selectedRequest, setSelectedRequest] = useState(null);
  const [actionLoading, setActionLoading] = useState(false);

  const handleDelete = async () => {
    if (!selectedRequest) return;
    setActionLoading(true);
    try {
      await deleteRequest(selectedRequest.id);
      setDeleteModalOpen(false);
      setSelectedRequest(null);
    } catch (error) {
      console.error('Failed to delete', error);
    } finally {
      setActionLoading(false);
    }
  };

  const handleEditSubmit = async (formData) => {
    if (!selectedRequest) return;
    setActionLoading(true);
    try {
      await updateRequest(selectedRequest.id, formData);
      setEditModalOpen(false);
      setSelectedRequest(null);
    } catch (error) {
      console.error('Failed to update', error);
    } finally {
      setActionLoading(false);
    }
  };

  const getUrgencyBadge = (urgency) => {
    const map = {
      low: 'neutral',
      standard: 'info',
      high: 'warning',
      critical: 'error'
    };
    return <Badge variant={map[urgency] || 'neutral'}>{urgency}</Badge>;
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
    { key: 'item_name', label: 'Item Name' },
    { key: 'item_category', label: 'Category', render: r => <span style={{ textTransform: 'capitalize' }}>{r.item_category}</span> },
    { key: 'urgency', label: 'Urgency', render: r => getUrgencyBadge(r.urgency) },
    { key: 'status', label: 'Status', render: r => getStatusBadge(r.status) },
    { key: 'estimated_cost', label: 'Est. Cost', render: r => formatCurrency(r.estimated_cost) },
    { key: 'created_at', label: 'Date', render: r => formatDate(r.created_at) },
    {
      key: 'actions',
      label: 'Actions',
      render: (r) => (
        <div style={{ display: 'flex', gap: 'var(--space-2)' }}>
          <Button size="sm" variant="ghost" onClick={() => { setSelectedRequest(r); setEditModalOpen(true); }}>Edit</Button>
          <Button size="sm" variant="ghost" style={{ color: 'var(--error)' }} onClick={() => { setSelectedRequest(r); setDeleteModalOpen(true); }}>Delete</Button>
        </div>
      )
    }
  ];

  return (
    <div style={{ padding: 'var(--space-6)', maxWidth: '1200px', margin: '0 auto' }}>
      <PageHeader 
        title="My Requests" 
        subtitle="Manage your requests in the community pool" 
        action={
          <Link to="/dashboard/requests/new">
            <Button variant="primary">+ New Request</Button>
          </Link>
        }
      />

      {loading ? (
        <div style={{ height: '300px', animation: 'shimmer 1.5s infinite', background: 'var(--bg-surface)', borderRadius: 'var(--radius-lg)' }}></div>
      ) : requests.length === 0 ? (
        <EmptyState 
          title="No requests yet" 
          description="Submit your first request to join the community queue."
          action={<Link to="/dashboard/requests/new"><Button variant="primary">Submit Request</Button></Link>}
        />
      ) : (
        <Table columns={columns} data={requests} />
      )}

      {/* Edit Modal */}
      <Modal open={editModalOpen} onClose={() => setEditModalOpen(false)} title="Edit Request">
        {selectedRequest && (
          <RequestForm 
            initialData={selectedRequest} 
            onSubmit={handleEditSubmit} 
            isLoading={actionLoading} 
          />
        )}
      </Modal>

      {/* Delete Confirmation Modal */}
      <Modal open={deleteModalOpen} onClose={() => setDeleteModalOpen(false)} title="Delete Request">
        <p style={{ marginBottom: 'var(--space-6)' }}>
          Are you sure you want to delete the request for <strong>{selectedRequest?.item_name}</strong>? This action cannot be undone.
        </p>
        <div style={{ display: 'flex', justifyContent: 'flex-end', gap: 'var(--space-3)' }}>
          <Button variant="ghost" onClick={() => setDeleteModalOpen(false)} disabled={actionLoading}>Cancel</Button>
          <Button variant="danger" onClick={handleDelete} disabled={actionLoading}>
            {actionLoading ? 'Deleting...' : 'Delete'}
          </Button>
        </div>
      </Modal>
    </div>
  );
}

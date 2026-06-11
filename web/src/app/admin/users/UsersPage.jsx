import { useState } from 'react';
import { useAdminUsers } from '../../../hooks/useAdminUsers';
import PageHeader from '../../../components/layout/PageHeader';
import Table from '../../../components/ui/Table';
import Badge from '../../../components/ui/Badge';
import Button from '../../../components/ui/Button';
import Modal from '../../../components/ui/Modal';
import { formatDate } from '../../../lib/formatters';

export default function UsersPage() {
  const { users, loading, verifyUser, deleteUser } = useAdminUsers();
  const [selectedUser, setSelectedUser] = useState(null);
  const [deleteModalOpen, setDeleteModalOpen] = useState(false);
  const [actionLoading, setActionLoading] = useState(false);

  const handleVerify = async (user) => {
    setActionLoading(true);
    try {
      await verifyUser(user.id);
    } catch (err) {
      console.error(err);
    } finally {
      setActionLoading(false);
    }
  };

  const handleDelete = async () => {
    if (!selectedUser) return;
    setActionLoading(true);
    try {
      await deleteUser(selectedUser.id);
      setDeleteModalOpen(false);
      setSelectedUser(null);
    } catch (err) {
      console.error(err);
    } finally {
      setActionLoading(false);
    }
  };

  const columns = [
    { key: 'name', label: 'Name' },
    { key: 'email', label: 'Email' },
    { key: 'role', label: 'Role', render: u => <Badge variant={u.role === 'admin' ? 'info' : 'neutral'}>{u.role}</Badge> },
    { key: 'verified', label: 'Verified', render: u => <Badge variant={u.verified ? 'success' : 'warning'}>{u.verified ? 'Yes' : 'No'}</Badge> },
    { key: 'joined_at', label: 'Joined', render: u => formatDate(u.joined_at) },
    {
      key: 'actions',
      label: 'Actions',
      render: u => (
        <div style={{ display: 'flex', gap: 'var(--space-2)' }}>
          {!u.verified && <Button size="sm" variant="ghost" onClick={() => handleVerify(u)} disabled={actionLoading}>Verify</Button>}
          {u.role !== 'admin' && (
            <Button size="sm" variant="ghost" style={{ color: 'var(--error)' }} onClick={() => { setSelectedUser(u); setDeleteModalOpen(true); }} disabled={actionLoading}>
              Delete
            </Button>
          )}
        </div>
      )
    }
  ];

  return (
    <div style={{ padding: 'var(--space-6)', maxWidth: '1200px', margin: '0 auto' }}>
      <PageHeader title="Users" subtitle="Manage members and institutions in Virtus" />

      {loading ? (
        <div style={{ height: '300px', animation: 'shimmer 1.5s infinite', background: 'var(--bg-surface)', borderRadius: 'var(--radius-lg)' }}></div>
      ) : (
        <Table columns={columns} data={users} />
      )}

      <Modal open={deleteModalOpen} onClose={() => setDeleteModalOpen(false)} title="Delete User">
        <p style={{ marginBottom: 'var(--space-6)' }}>
          Are you sure you want to delete <strong>{selectedUser?.email}</strong>? This action is permanent and will remove their history.
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

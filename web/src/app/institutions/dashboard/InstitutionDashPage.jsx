import { useState } from 'react';
import { useInstitution } from '../../../hooks/useInstitution';
import PageHeader from '../../../components/layout/PageHeader';
import Card from '../../../components/ui/Card';
import Badge from '../../../components/ui/Badge';
import Table from '../../../components/ui/Table';
import Button from '../../../components/ui/Button';
import Modal from '../../../components/ui/Modal';
import Input from '../../../components/ui/Input';
import { formatCurrency, formatDate } from '../../../lib/formatters';

export default function InstitutionDashPage() {
  const { institution, contributions, loading, createContribution } = useInstitution();
  const [modalOpen, setModalOpen] = useState(false);
  const [actionLoading, setActionLoading] = useState(false);
  const [contribForm, setContribForm] = useState({ amount: '', tags: '' });

  const handleContribute = async (e) => {
    e.preventDefault();
    if (!institution) return;
    setActionLoading(true);
    try {
      const tagsArray = contribForm.tags.split(',').map(t => t.trim()).filter(Boolean);
      await createContribution(institution.id, {
        amount: Number(contribForm.amount),
        currency: 'USD',
        category_tags: tagsArray
      });

      setModalOpen(false);
      setContribForm({ amount: '', tags: '' });
    } catch (err) {
      console.error(err);
    } finally {
      setActionLoading(false);
    }
  };

  const columns = [
    { key: 'amount', label: 'Amount', render: c => formatCurrency(c.amount) },
    { key: 'tags', label: 'Tags', render: c => (
      <div style={{ display: 'flex', gap: 'var(--space-2)' }}>
        {(c.category_tags || []).map(t => <Badge key={t} variant="info">{t}</Badge>)}
      </div>
    )},
    { key: 'status', label: 'Status', render: c => <Badge variant={c.status === 'completed' ? 'success' : 'warning'}>{c.status}</Badge> },
    { key: 'date', label: 'Date', render: c => formatDate(c.created_at) },
  ];

  if (loading) {
    return <div style={{ padding: 'var(--space-6)' }}><div style={{ height: '300px', animation: 'shimmer 1.5s infinite', background: 'var(--bg-surface)' }}></div></div>;
  }

  return (
    <div style={{ padding: 'var(--space-6)', maxWidth: '1200px', margin: '0 auto' }}>
      <PageHeader
        title="Institution Dashboard"
        subtitle="Manage your institutional impact and directed pool contributions"
        action={<Button variant="primary" onClick={() => setModalOpen(true)}>New Contribution</Button>}
      />

      <Card style={{ marginBottom: 'var(--space-8)', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <div>
          <div style={{ display: 'flex', alignItems: 'center', gap: 'var(--space-3)', marginBottom: 'var(--space-2)' }}>
            <h2 style={{ fontSize: 'var(--font-xl)', color: 'var(--text-primary)' }}>{institution?.name || 'Unnamed Institution'}</h2>
            {institution?.verified && <Badge variant="success">Verified</Badge>}
          </div>
          <div style={{ color: 'var(--text-secondary)' }}>
            Type: <span style={{ color: 'var(--text-primary)', textTransform: 'capitalize' }}>{institution?.institution_type || 'N/A'}</span>
          </div>
          <div style={{ color: 'var(--text-secondary)' }}>
            ESG Rating: <span style={{ color: 'var(--text-primary)' }}>{institution?.esg_rating || 'Pending'}</span>
          </div>
        </div>
        <div style={{ textAlign: 'right' }}>
          <div style={{ fontSize: 'var(--font-sm)', color: 'var(--text-secondary)' }}>Total Contributed</div>
          <div style={{ fontFamily: 'var(--font-display)', fontSize: 'var(--font-2xl)', color: 'var(--solar-300)', letterSpacing: '-0.03em', textShadow: '0 0 20px rgba(242,201,76,0.2)' }}>
            {formatCurrency(contributions.reduce((acc, c) => acc + c.amount, 0))}
          </div>
        </div>
      </Card>

      <div>
        <h3 style={{ fontSize: 'var(--font-lg)', marginBottom: 'var(--space-4)', color: 'var(--text-primary)' }}>Contribution History</h3>
        <Table columns={columns} data={contributions} />
      </div>

      <Modal open={modalOpen} onClose={() => setModalOpen(false)} title="New Directed Contribution">
        <form onSubmit={handleContribute} style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-4)' }}>
          <Input type="number" placeholder="Amount ($)" value={contribForm.amount} onChange={e => setContribForm(p => ({ ...p, amount: e.target.value }))} required min="1" />
          <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-2)' }}>
            <label style={{ fontSize: 'var(--font-sm)', color: 'var(--text-secondary)' }}>Tags (comma separated)</label>
            <Input placeholder="e.g. medical, housing, new-york" value={contribForm.tags} onChange={e => setContribForm(p => ({ ...p, tags: e.target.value }))} />
            <span style={{ fontSize: 'var(--font-xs)', color: 'var(--text-secondary)' }}>Tags help direct funds to specific categories or regions if supported by the pool.</span>
          </div>
          <Button type="submit" variant="primary" disabled={actionLoading}>Initiate Contribution</Button>
        </form>
      </Modal>
    </div>
  );
}

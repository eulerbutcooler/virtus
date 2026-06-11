import { useState } from 'react';
import { loadStripe } from '@stripe/stripe-js';
import { Elements } from '@stripe/react-stripe-js';
import { useContributions } from '../../../hooks/useContributions';
import PageHeader from '../../../components/layout/PageHeader';
import Table from '../../../components/ui/Table';
import Badge from '../../../components/ui/Badge';
import Button from '../../../components/ui/Button';
import Modal from '../../../components/ui/Modal';
import Input from '../../../components/ui/Input';
import EmptyState from '../../../components/ui/EmptyState';
import ContributionForm from '../../../components/forms/ContributionForm';
import { formatCurrency, formatDate } from '../../../lib/formatters';

const stripePromise = loadStripe(import.meta.env.VITE_STRIPE_PUBLIC_KEY || 'pk_test_51MockStripeKey1234567890');

export default function ContributionsPage() {
  const { contributions, total, loading, createContribution, refresh, optimisticUpdate } = useContributions();
  const [modalOpen, setModalOpen] = useState(false);
  const [amount, setAmount] = useState('');
  const [clientSecret, setClientSecret] = useState(null);
  const [initLoading, setInitLoading] = useState(false);
  const [error, setError] = useState(null);

  const handleInitiate = async (e) => {
    e.preventDefault();
    setError(null);
    const numAmount = Number(amount);
    if (numAmount <= 0) {
      setError('Amount must be greater than zero');
      return;
    }

    setInitLoading(true);
    try {
      const data = await createContribution(numAmount, 'USD');
      if (data.client_secret) {
        setClientSecret(data.client_secret);

        optimisticUpdate({
          id: `temp-${Date.now()}`,
          amount: numAmount,
          currency: 'USD',
          status: 'processing',
          created_at: new Date().toISOString()
        });
      } else {

        setModalOpen(false);
        setAmount('');
      }
    } catch (err) {
      setError(err.message || 'Failed to initiate contribution');
    } finally {
      setInitLoading(false);
    }
  };

  const handleSuccess = () => {
    setModalOpen(false);
    setClientSecret(null);
    setAmount('');
    refresh(); // fetch latest status (should be completed)
  };

  const handleCancel = () => {
    setModalOpen(false);
    setClientSecret(null);
    setAmount('');
    refresh(); // refresh to correct the optimistic update if canceled
  };

  const getStatusBadge = (status) => {
    const map = {
      pending: 'warning',
      processing: 'info',
      completed: 'success',
      failed: 'error'
    };
    return <Badge variant={map[status] || 'neutral'}>{status}</Badge>;
  };

  const columns = [
    { key: 'amount', label: 'Amount', render: c => formatCurrency(c.amount, c.currency) },
    { key: 'status', label: 'Status', render: c => getStatusBadge(c.status) },
    { key: 'payment_ref', label: 'Ref', render: c => c.payment_ref ? <span style={{ fontFamily: 'monospace', color: 'var(--text-secondary)' }}>{c.payment_ref.substring(0, 12)}...</span> : '-' },
    { key: 'created_at', label: 'Date', render: c => formatDate(c.created_at) }
  ];

  return (
    <div style={{ padding: 'var(--space-6)', maxWidth: '1200px', margin: '0 auto' }}>
      <PageHeader
        title="My Contributions"
        subtitle="Your impact on the community pool"
        action={
          <Button variant="primary" onClick={() => setModalOpen(true)}>Contribute</Button>
        }
      />

      <div style={{ marginBottom: 'var(--space-8)' }}>
        <h2 style={{ fontSize: 'var(--font-sm)', color: 'var(--text-secondary)', textTransform: 'uppercase', marginBottom: 'var(--space-2)' }}>Total Contributed</h2>
        <div style={{ fontFamily: 'var(--font-display)', fontSize: 'var(--font-2xl)', color: 'var(--solar-300)', letterSpacing: '-0.03em', textShadow: '0 0 20px rgba(242,201,76,0.2)' }}>
          {formatCurrency(total?.amount || 0)}
        </div>
      </div>

      {loading && contributions.length === 0 ? (
        <div style={{ height: '300px', animation: 'shimmer 1.5s infinite', background: 'var(--bg-surface)', borderRadius: 'var(--radius-lg)' }}></div>
      ) : contributions.length === 0 ? (
        <EmptyState
          title="No contributions yet"
          description="Your first contribution helps fulfill a community need."
          action={<Button variant="primary" onClick={() => setModalOpen(true)}>Make a Contribution</Button>}
        />
      ) : (
        <Table columns={columns} data={contributions} />
      )}

      <Modal open={modalOpen} onClose={handleCancel} title="Make a Contribution">
        {!clientSecret ? (
          <form onSubmit={handleInitiate} style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-4)' }}>
            <p style={{ color: 'var(--text-secondary)' }}>Enter the amount you would like to contribute to the community pool.</p>
            {error && <div style={{ color: 'var(--error)', fontSize: 'var(--font-sm)' }}>{error}</div>}
            <Input
              type="number"
              placeholder="Amount ($)"
              value={amount}
              onChange={(e) => setAmount(e.target.value)}
              min="1"
              required
            />
            <div style={{ display: 'flex', gap: 'var(--space-4)', marginTop: 'var(--space-4)' }}>
              <Button type="submit" variant="primary" disabled={initLoading}>
                {initLoading ? 'Processing...' : 'Continue to Payment'}
              </Button>
              <Button type="button" variant="ghost" onClick={handleCancel} disabled={initLoading}>
                Cancel
              </Button>
            </div>
          </form>
        ) : (
          <Elements stripe={stripePromise} options={{ clientSecret, appearance: { theme: 'night', variables: { colorPrimary: '#5CAB6E', colorBackground: '#182B1C', colorText: '#DFF0E2' } } }}>
            <ContributionForm
              amount={amount}
              clientSecret={clientSecret}
              onSuccess={handleSuccess}
              onCancel={handleCancel}
            />
          </Elements>
        )}
      </Modal>
    </div>
  );
}

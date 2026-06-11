import { useState } from 'react';
import { useImpact } from '../../../hooks/useImpact';
import { useRequests } from '../../../hooks/useRequests';
import PageHeader from '../../../components/layout/PageHeader';
import Card from '../../../components/ui/Card';
import Button from '../../../components/ui/Button';
import Modal from '../../../components/ui/Modal';
import ImpactForm from '../../../components/forms/ImpactForm';
import EmptyState from '../../../components/ui/EmptyState';
import { formatDate } from '../../../lib/formatters';
import { HugeiconsIcon } from '@hugeicons/react';
import { StarIcon } from '@hugeicons/core-free-icons';

export default function ImpactPage() {
  const { impacts, loading, createImpact } = useImpact();
  const { requests } = useRequests();
  const [modalOpen, setModalOpen] = useState(false);
  const [actionLoading, setActionLoading] = useState(false);

  const deliveredRequests = requests?.filter(r => r.status === 'delivered') || [];
  const mockDeliveries = deliveredRequests.map(r => ({
    id: `del-${r.id}`,
    fulfillment: { request: r }
  }));

  const handleRecord = async (formData) => {
    setActionLoading(true);
    try {
      await createImpact(formData);
      setModalOpen(false);
    } catch (err) {
      console.error(err);
    } finally {
      setActionLoading(false);
    }
  };

  const renderStars = (score) => {
    return Array.from({ length: 5 }).map((_, i) => (
      <HugeiconsIcon
        key={i}
        icon={StarIcon}
        size={16}
        color={i < score ? 'var(--solar-300)' : 'var(--border-default)'}
        strokeWidth={2}
        style={{ fill: i < score ? 'var(--solar-300)' : 'transparent' }}
      />
    ));
  };

  return (
    <div style={{ padding: 'var(--space-6)', maxWidth: '1000px', margin: '0 auto' }}>
      <PageHeader
        title="My Impact"
        subtitle="Share how the community pool has helped you"
        action={
          <Button variant="primary" onClick={() => setModalOpen(true)} disabled={mockDeliveries.length === 0}>
            + Record Impact
          </Button>
        }
      />

      {loading ? (
        <div style={{ height: '300px', animation: 'shimmer 1.5s infinite', background: 'var(--bg-surface)', borderRadius: 'var(--radius-lg)' }}></div>
      ) : impacts.length === 0 ? (
        <EmptyState
          title="No impact records"
          description={mockDeliveries.length === 0 ? "You'll be able to record impact once your requests are delivered." : "Record your first impact update to inspire the community!"}
        />
      ) : (
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(300px, 1fr))', gap: 'var(--space-6)' }}>
          {impacts.map(impact => (
            <Card key={impact.id} style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-4)' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
                <div>
                  <div style={{ fontSize: 'var(--font-sm)', color: 'var(--text-secondary)', textTransform: 'capitalize' }}>
                    {impact.interval.replace('_', ' ')}
                  </div>
                  <div style={{ fontSize: 'var(--font-md)', color: 'var(--text-primary)', fontWeight: 'bold' }}>
                    {impact.outcome}
                  </div>
                </div>
                <div style={{ display: 'flex', gap: '2px' }}>
                  {renderStars(impact.satisfaction_score)}
                </div>
              </div>
              <p style={{ color: 'var(--text-secondary)', fontSize: 'var(--font-sm)', lineHeight: '1.5' }}>
                "{impact.description}"
              </p>
              <div style={{ fontSize: 'var(--font-xs)', color: 'var(--text-secondary)', marginTop: 'auto' }}>
                Recorded on {formatDate(impact.created_at)}
              </div>
            </Card>
          ))}
        </div>
      )}

      <Modal open={modalOpen} onClose={() => setModalOpen(false)} title="Record Impact">
        <ImpactForm
          deliveries={mockDeliveries}
          onSubmit={handleRecord}
          onCancel={() => setModalOpen(false)}
          isLoading={actionLoading}
        />
      </Modal>
    </div>
  );
}

import { useState, useEffect } from 'react';
import { useRequests } from '../../../hooks/useRequests';
import { api } from '../../../lib/api';
import PageHeader from '../../../components/layout/PageHeader';
import Table from '../../../components/ui/Table';
import ProgressBar from '../../../components/ui/ProgressBar';
import EmptyState from '../../../components/ui/EmptyState';
import Button from '../../../components/ui/Button';
import { Link } from 'react-router-dom';

export default function QueuePage() {
  const { requests, loading: reqLoading } = useRequests();
  const [queueEntries, setQueueEntries] = useState([]);
  const [queueLoading, setQueueLoading] = useState(true);

  useEffect(() => {
    async function fetchMyQueue() {
      if (reqLoading) return;

      const activeRequests = requests.filter(r => r.status === 'pending' || r.status === 'queued' || r.status === 'funded');

      if (activeRequests.length === 0) {
        setQueueLoading(false);
        return;
      }

      try {
        setQueueLoading(true);

        const entries = await Promise.all(
          activeRequests.map(async (req) => {
            try {
              const data = await api.get(`/queue/${req.id}`);
              return { ...data, request: req };
            } catch {
              return null;
            }
          })
        );
        setQueueEntries(entries.filter(Boolean));
      } catch (error) {
        console.error('Failed to fetch queue entries', error);
      } finally {
        setQueueLoading(false);
      }
    }
    fetchMyQueue();
  }, [requests, reqLoading]);

  const columns = [
    { key: 'position', label: 'Position', render: q => <strong style={{ fontFamily: 'var(--font-display)', color: 'var(--leaf-300)', fontSize: 'var(--font-md)' }}>#{q.position}</strong> },
    { key: 'item_name', label: 'Item', render: q => <span style={{ color: 'var(--text-primary)' }}>{q.request?.item_name}</span> },
    {
      key: 'funding',
      label: 'Funding Progress',
      render: q => (
        <div style={{ minWidth: '150px' }}>
          <ProgressBar value={q.funding_progress || 0} label={`${Math.floor(q.funding_progress || 0)}%`} />
        </div>
      )
    },
    { key: 'eta', label: 'ETA', render: q => <span style={{ color: 'var(--text-secondary)' }}>{q.estimated_fulfillment || 'Pending'}</span> }
  ];

  const loading = reqLoading || queueLoading;

  return (
    <div style={{ padding: 'var(--space-6)', maxWidth: '1000px', margin: '0 auto' }}>
      <PageHeader
        title="My Queue"
        subtitle="Track the status of your active requests in the community queue."
      />

      {loading ? (
        <div style={{ height: '300px', animation: 'shimmer 1.5s infinite', background: 'var(--bg-surface)', borderRadius: 'var(--radius-lg)' }}></div>
      ) : queueEntries.length === 0 ? (
        <EmptyState
          title="No active queue entries"
          description="You don't have any requests currently waiting in the queue."
          action={<Link to="/dashboard/requests/new"><Button variant="primary">Submit Request</Button></Link>}
        />
      ) : (
        <Table columns={columns} data={queueEntries} />
      )}
    </div>
  );
}

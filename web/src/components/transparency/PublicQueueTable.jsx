import Table from '../ui/Table';
import ProgressBar from '../ui/ProgressBar';

export default function PublicQueueTable({ queue, loading }) {
  if (loading) {
    return (
      <div style={{ height: '300px', animation: 'shimmer 1.5s infinite', background: 'var(--bg-surface)', borderRadius: 'var(--radius-lg)' }}></div>
    );
  }

  if (!queue || queue.length === 0) {
    return (
      <div style={{ padding: 'var(--space-8)', textAlign: 'center', background: 'var(--bg-elevated)', border: '1px solid var(--border-default)', borderRadius: 'var(--radius-lg)' }}>
        <div style={{ color: 'var(--text-secondary)' }}>The community queue is currently empty.</div>
      </div>
    );
  }

  const columns = [
    { key: 'position', label: 'Position', render: q => <strong style={{ color: 'var(--beige-200)' }}>#{q.position}</strong> },
    { key: 'item_name', label: 'Requested Item', render: q => <span style={{ color: 'var(--text-primary)' }}>{q.request?.item_name || 'Anonymous Item'}</span> },
    { 
      key: 'funding', 
      label: 'Funding Progress', 
      render: q => (
        <div style={{ minWidth: '150px' }}>
          <ProgressBar value={q.funding_progress || 0} label={`${Math.floor(q.funding_progress || 0)}%`} />
        </div>
      ) 
    },
    { key: 'eta', label: 'Estimated Fulfillment', render: q => <span style={{ color: 'var(--text-secondary)' }}>{q.estimated_fulfillment || 'Pending'}</span> }
  ];

  return <Table columns={columns} data={queue} />;
}

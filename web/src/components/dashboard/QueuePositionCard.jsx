import Card from '../ui/Card';
import ProgressBar from '../ui/ProgressBar';
import Button from '../ui/Button';
import { Link } from 'react-router-dom';

export default function QueuePositionCard({ queueEntry, loading }) {
  if (loading) {
    return (
      <Card style={{ height: '100%' }}>
        <div style={{ height: '100px', animation: 'shimmer 1.5s infinite', background: 'var(--bg-surface)', borderRadius: 'var(--radius-sm)' }}></div>
      </Card>
    );
  }

  if (!queueEntry) {
    return (
      <Card style={{ height: '100%', display: 'flex', flexDirection: 'column', justifyContent: 'center', alignItems: 'center', padding: 'var(--space-8)' }}>
        <div style={{ color: 'var(--text-secondary)', marginBottom: 'var(--space-4)', textAlign: 'center' }}>Submit a request to join the queue</div>
        <Link to="/dashboard/requests/new">
          <Button variant="primary">New Request</Button>
        </Link>
      </Card>
    );
  }

  return (
    <Card style={{ height: '100%', display: 'flex', flexDirection: 'column', justifyContent: 'space-between' }}>
      <div>
        <h2 style={{ fontSize: 'var(--font-xs)', color: 'var(--text-secondary)', textTransform: 'uppercase', letterSpacing: '0.5px', marginBottom: 'var(--space-4)' }}>My Queue Position</h2>
        <div style={{ fontSize: 'var(--font-2xl)', color: 'var(--beige-200)', marginBottom: 'var(--space-4)' }}>
          #{queueEntry.position}
        </div>
        <ProgressBar value={queueEntry.funding_progress || 0} max={100} label={`${queueEntry.funding_progress || 0}% Funded`} />
      </div>
      <div style={{ fontSize: 'var(--font-sm)', color: 'var(--text-secondary)', marginTop: 'var(--space-4)' }}>
        ETA: <span style={{ color: 'var(--text-primary)' }}>{queueEntry.estimated_fulfillment || 'Pending'}</span>
      </div>
    </Card>
  );
}

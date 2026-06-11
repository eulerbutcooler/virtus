import Card from '../ui/Card';
import ProgressBar from '../ui/ProgressBar';
import Button from '../ui/Button';
import { Link } from 'react-router-dom';

export default function QueuePositionCard({ queueEntry, loading }) {
  if (loading) {
    return (
      <Card style={{ height: '100%' }}>
        <div style={{
          height: '100px',
          background: 'linear-gradient(90deg, var(--bg-surface) 25%, var(--bg-hover) 50%, var(--bg-surface) 75%)',
          backgroundSize: '200% 100%',
          animation: 'shimmer 1.5s linear infinite',
          borderRadius: 'var(--radius-sm)'
        }} />
      </Card>
    );
  }

  if (!queueEntry) {
    return (
      <Card style={{ height: '100%', display: 'flex', flexDirection: 'column', justifyContent: 'center', alignItems: 'center', padding: 'var(--space-8)' }}>
        <div style={{
          fontFamily: 'var(--font-display)',
          fontSize: 'var(--font-xl)',
          color: 'var(--text-disabled)',
          marginBottom: 'var(--space-2)',
          textAlign: 'center',
          lineHeight: 1.2
        }}>
          No active request
        </div>
        <div style={{ color: 'var(--text-secondary)', marginBottom: 'var(--space-5)', textAlign: 'center', fontSize: 'var(--font-sm)' }}>
          Submit a request to join the queue — your first one starts the cycle
        </div>
        <Link to="/dashboard/requests/new">
          <Button variant="primary">Submit a request</Button>
        </Link>
      </Card>
    );
  }

  return (
    <Card style={{ height: '100%', display: 'flex', flexDirection: 'column', justifyContent: 'space-between' }}>
      <div>
        <div style={{
          fontSize: 'var(--font-xs)',
          color: 'var(--text-secondary)',
          textTransform: 'uppercase',
          letterSpacing: '0.06em',
          fontWeight: 'var(--font-weight-bold)',
          marginBottom: 'var(--space-4)'
        }}>
          Queue Position
        </div>

        <div style={{
          fontFamily: 'var(--font-display)',
          fontSize: 'var(--font-3xl)',
          fontWeight: 600,
          letterSpacing: '-0.04em',
          color: 'var(--leaf-300)',
          lineHeight: 1,
          marginBottom: 'var(--space-5)',
          textShadow: '0 0 20px rgba(92, 171, 110, 0.25)'
        }}>
          #{queueEntry.position}
        </div>
        <ProgressBar value={queueEntry.funding_progress || 0} max={100} label={`${queueEntry.funding_progress || 0}% funded`} />
      </div>
      <div style={{
        fontSize: 'var(--font-xs)',
        color: 'var(--text-secondary)',
        marginTop: 'var(--space-4)',
        letterSpacing: '0.03em'
      }}>
        ETA: <span style={{ color: 'var(--text-primary)', fontWeight: 'var(--font-weight-medium)' }}>
          {queueEntry.estimated_fulfillment || 'Pending review'}
        </span>
      </div>
    </Card>
  );
}

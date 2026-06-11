import Card from '../ui/Card';
import Badge from '../ui/Badge';
import { formatRelativeTime } from '../../lib/formatters';
import { Link } from 'react-router-dom';

export default function RecentActivityFeed({ requests, contributions, loading }) {
  if (loading) {
    return (
      <Card>
        {[1,2,3].map(i => (
          <div key={i} style={{
            height: '48px',
            marginBottom: i < 3 ? 'var(--space-3)' : 0,
            background: 'linear-gradient(90deg, var(--bg-surface) 25%, var(--bg-hover) 50%, var(--bg-surface) 75%)',
            backgroundSize: '200% 100%',
            animation: `shimmer 1.5s linear infinite`,
            animationDelay: `${i * 0.1}s`,
            borderRadius: 'var(--radius-sm)'
          }} />
        ))}
      </Card>
    );
  }

  const reqActivities = (requests || []).map(r => ({
    id: `req-${r.id}`,
    type: 'request',
    date: r.updated_at || r.created_at,
    desc: `Request for ${r.item_name}`,
    status: r.status
  }));

  const contribActivities = (contributions || []).map(c => ({
    id: `con-${c.id}`,
    type: 'contribution',
    date: c.created_at,
    desc: `Contributed ${c.amount} ${c.currency}`,
    status: c.status
  }));

  const activities = [...reqActivities, ...contribActivities]
    .sort((a, b) => new Date(b.date) - new Date(a.date))
    .slice(0, 8);

  const getStatusVariant = (status) => {
    switch (status) {
      case 'completed':
      case 'verified':
      case 'delivered':
      case 'funded':
        return 'success';
      case 'pending':
      case 'queued':
        return 'warning';
      case 'rejected':
      case 'failed':
        return 'error';
      default:
        return 'neutral';
    }
  };

  const getDotColor = (type) => type === 'contribution' ? 'var(--solar-300)' : 'var(--leaf-300)';

  return (
    <Card>
      <div style={{
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'baseline',
        marginBottom: 'var(--space-5)'
      }}>
        <h3 style={{
          fontFamily: 'var(--font-display)',
          fontSize: 'var(--font-lg)',
          fontWeight: 600,
          letterSpacing: '-0.02em',
          color: 'var(--text-primary)'
        }}>
          Recent Activity
        </h3>
        {activities.length > 0 && (
          <Link to="/dashboard/requests" style={{ fontSize: 'var(--font-xs)' }}>
            View all →
          </Link>
        )}
      </div>

      {activities.length === 0 ? (
        <div style={{
          color: 'var(--text-secondary)',
          padding: 'var(--space-6) 0',
          fontSize: 'var(--font-sm)',
          textAlign: 'center',
          letterSpacing: '0.01em'
        }}>
          No contributions yet — your first one starts the cycle
        </div>
      ) : (
        <div style={{ display: 'flex', flexDirection: 'column', position: 'relative' }}>
          {activities.map((activity, idx) => (
            <div key={activity.id} style={{
              display: 'flex',
              alignItems: 'center',
              gap: 'var(--space-4)',
              padding: 'var(--space-3) 0',
              borderBottom: idx !== activities.length - 1 ? '1px solid var(--border-default)' : 'none',
            }}>

              <div style={{
                width: 8,
                height: 8,
                borderRadius: '50%',
                background: getDotColor(activity.type),
                flexShrink: 0,
                boxShadow: `0 0 6px ${getDotColor(activity.type)}`
              }} />

              <div style={{ flex: 1, minWidth: 0 }}>
                <div style={{
                  color: 'var(--text-primary)',
                  fontSize: 'var(--font-sm)',
                  marginBottom: '2px',
                  overflow: 'hidden',
                  textOverflow: 'ellipsis',
                  whiteSpace: 'nowrap'
                }}>
                  {activity.desc}
                </div>
                <div style={{
                  color: 'var(--text-disabled)',
                  fontSize: 'var(--font-xs)',
                  letterSpacing: '0.02em'
                }}>
                  {formatRelativeTime(activity.date)}
                </div>
              </div>

              <Badge variant={getStatusVariant(activity.status)}>
                {activity.status}
              </Badge>
            </div>
          ))}
        </div>
      )}
    </Card>
  );
}

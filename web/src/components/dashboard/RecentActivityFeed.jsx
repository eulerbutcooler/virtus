import Card from '../ui/Card';
import Badge from '../ui/Badge';
import { formatRelativeTime } from '../../lib/formatters';

export default function RecentActivityFeed({ requests, contributions, loading }) {
  if (loading) {
    return (
      <Card>
        <div style={{ height: '200px', animation: 'shimmer 1.5s infinite', background: 'var(--bg-surface)', borderRadius: 'var(--radius-sm)' }}></div>
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

  if (activities.length === 0) {
    return (
      <Card>
        <h2 style={{ fontSize: 'var(--font-md)', marginBottom: 'var(--space-4)', color: 'var(--text-primary)' }}>Recent Activity</h2>
        <div style={{ color: 'var(--text-secondary)', padding: 'var(--space-4) 0', fontSize: 'var(--font-sm)' }}>No recent activity to show.</div>
      </Card>
    );
  }

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

  return (
    <Card>
      <h2 style={{ fontSize: 'var(--font-md)', marginBottom: 'var(--space-4)', color: 'var(--text-primary)' }}>Recent Activity</h2>
      <div style={{ display: 'flex', flexDirection: 'column' }}>
        {activities.map((activity, idx) => (
          <div key={activity.id} style={{ 
            display: 'flex', 
            justifyContent: 'space-between', 
            alignItems: 'center', 
            padding: 'var(--space-3) 0',
            borderBottom: idx !== activities.length - 1 ? '1px solid var(--border-default)' : 'none'
          }}>
            <div>
              <div style={{ color: 'var(--text-primary)', fontSize: 'var(--font-sm)', marginBottom: 'var(--space-1)' }}>{activity.desc}</div>
              <div style={{ color: 'var(--text-secondary)', fontSize: 'var(--font-xs)' }}>{formatRelativeTime(activity.date)}</div>
            </div>
            <Badge variant={getStatusVariant(activity.status)}>
              {activity.status}
            </Badge>
          </div>
        ))}
      </div>
    </Card>
  );
}

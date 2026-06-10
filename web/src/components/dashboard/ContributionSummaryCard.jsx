import Card from '../ui/Card';
import { formatCurrency } from '../../lib/formatters';

export default function ContributionSummaryCard({ total, loading }) {
  if (loading) {
    return (
      <Card style={{ height: '100%' }}>
        <div style={{ height: '100px', animation: 'shimmer 1.5s infinite', background: 'var(--bg-surface)', borderRadius: 'var(--radius-sm)' }}></div>
      </Card>
    );
  }

  return (
    <Card style={{ height: '100%' }}>
      <h2 style={{ fontSize: 'var(--font-xs)', color: 'var(--text-secondary)', textTransform: 'uppercase', letterSpacing: '0.5px', marginBottom: 'var(--space-4)' }}>My Contributions</h2>
      <div style={{ fontSize: 'var(--font-2xl)', color: 'var(--beige-200)', marginBottom: 'var(--space-2)' }}>
        {formatCurrency(total?.amount || 0)}
      </div>
      <div style={{ fontSize: 'var(--font-sm)', color: 'var(--text-secondary)' }}>
        {total?.count || 0} total contributions
      </div>
    </Card>
  );
}

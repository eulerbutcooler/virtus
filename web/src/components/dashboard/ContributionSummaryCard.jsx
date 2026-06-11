import Card from '../ui/Card';
import { formatCurrency } from '../../lib/formatters';
import { Link } from 'react-router-dom';

export default function ContributionSummaryCard({ total, loading }) {
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

  const count = total?.count || 0;

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
          My Contributions
        </div>
        <div style={{
          fontFamily: 'var(--font-display)',
          fontSize: 'var(--font-2xl)',
          fontWeight: 600,
          letterSpacing: '-0.03em',
          color: 'var(--solar-300)',
          lineHeight: 1,
          marginBottom: 'var(--space-2)',
          textShadow: '0 0 20px rgba(242, 201, 76, 0.2)'
        }}>
          {formatCurrency(total?.amount || 0)}
        </div>
        <div style={{ fontSize: 'var(--font-xs)', color: 'var(--text-secondary)' }}>
          {count === 0
            ? 'No contributions yet — yours starts the cycle'
            : `${count} contribution${count !== 1 ? 's' : ''} to the pool`}
        </div>
      </div>
      <Link to="/dashboard/contributions" style={{
        fontSize: 'var(--font-xs)',
        color: 'var(--leaf-300)',
        fontWeight: 'var(--font-weight-medium)',
        letterSpacing: '0.02em',
        marginTop: 'var(--space-4)',
        display: 'inline-block'
      }}>
        Contribute more →
      </Link>
    </Card>
  );
}

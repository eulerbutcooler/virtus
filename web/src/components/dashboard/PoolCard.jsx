import Card from '../ui/Card';
import { formatCurrency } from '../../lib/formatters';

export default function PoolCard({ pool, loading }) {
  if (loading) {
    return (
      <Card style={{ marginBottom: 'var(--space-6)' }}>
        <div style={{ height: '80px', animation: 'shimmer 1.5s infinite', background: 'var(--bg-surface)', borderRadius: 'var(--radius-sm)' }}></div>
      </Card>
    );
  }

  const balance = pool?.balance || 0;
  const totalIn = pool?.total_in || 0;
  const totalOut = pool?.total_out || 0;
  
  const healthRatio = totalIn > 0 ? (totalOut / totalIn) * 100 : 0;

  return (
    <Card style={{ marginBottom: 'var(--space-6)' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-end', marginBottom: 'var(--space-4)' }}>
        <div>
          <h2 style={{ fontSize: 'var(--font-xs)', color: 'var(--text-secondary)', textTransform: 'uppercase', letterSpacing: '0.5px', marginBottom: 'var(--space-2)' }}>Community Pool</h2>
          <div style={{ fontSize: 'var(--font-2xl)', color: 'var(--beige-200)', lineHeight: '1' }}>
            {formatCurrency(balance)}
          </div>
        </div>
        <div style={{ textAlign: 'right', fontSize: 'var(--font-sm)', color: 'var(--text-secondary)' }}>
          <div>Total In: <span style={{ color: 'var(--text-primary)' }}>{formatCurrency(totalIn)}</span></div>
          <div>Total Out: <span style={{ color: 'var(--text-primary)' }}>{formatCurrency(totalOut)}</span></div>
        </div>
      </div>
      <div style={{ width: '100%', height: '2px', background: 'var(--bg-surface)', marginTop: 'var(--space-4)', position: 'relative' }}>
        <div style={{ position: 'absolute', left: 0, top: 0, bottom: 0, width: `${healthRatio}%`, background: 'var(--beige-300)' }}></div>
      </div>
    </Card>
  );
}

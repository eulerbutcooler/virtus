import { usePool } from '../../hooks/usePool';
import { formatCurrency } from '../../lib/formatters';

export default function PoolOverview() {
  const { pool, loading } = usePool();

  if (loading) {
    return (
      <div style={{ height: '140px', animation: 'shimmer 1.5s infinite', background: 'var(--bg-surface)', borderRadius: 'var(--radius-lg)', marginBottom: 'var(--space-6)' }}></div>
    );
  }

  const balance = pool?.balance || 0;
  const totalIn = pool?.total_in || 0;
  const totalOut = pool?.total_out || 0;
  const healthRatio = totalIn > 0 ? (totalOut / totalIn) * 100 : 0;

  return (
    <div style={{ background: 'var(--bg-elevated)', border: '1px solid var(--border-default)', borderRadius: 'var(--radius-lg)', padding: 'var(--space-8)', marginBottom: 'var(--space-8)', display: 'flex', flexDirection: 'column', gap: 'var(--space-4)' }}>
      <div>
        <h2 style={{ fontSize: 'var(--font-sm)', color: 'var(--text-secondary)', textTransform: 'uppercase', letterSpacing: '0.5px', marginBottom: 'var(--space-2)' }}>Community Pool Balance</h2>
        <div style={{ fontFamily: 'var(--font-display)', fontSize: 'var(--font-2xl)', color: 'var(--leaf-200)', lineHeight: '1', letterSpacing: '-0.03em', textShadow: '0 0 24px rgba(92,171,110,0.2)' }}>
          {formatCurrency(balance)}
        </div>
      </div>

      <div style={{ width: '100%', height: '4px', background: 'var(--bg-surface)', borderRadius: 'var(--radius-sm)', position: 'relative', overflow: 'hidden' }}>
          <div style={{ position: 'absolute', left: 0, top: 0, bottom: 0, width: `${healthRatio}%`, background: 'linear-gradient(90deg, var(--leaf-400), var(--leaf-300) 60%, #A8E6B0)', backgroundSize: '200% 100%', animation: 'leafShimmer 4s linear infinite', boxShadow: '0 0 8px rgba(92,171,110,0.4)', borderRadius: '999px', transition: 'width 800ms ease-out' }}></div>
      </div>

      <div style={{ display: 'flex', gap: 'var(--space-8)' }}>
        <div>
          <div style={{ fontSize: 'var(--font-xs)', color: 'var(--text-secondary)', textTransform: 'uppercase' }}>Total Distributed</div>
          <div style={{ fontSize: 'var(--font-md)', color: 'var(--text-primary)' }}>{formatCurrency(totalOut)}</div>
        </div>
        <div>
          <div style={{ fontSize: 'var(--font-xs)', color: 'var(--text-secondary)', textTransform: 'uppercase' }}>Total Contributed</div>
          <div style={{ fontSize: 'var(--font-md)', color: 'var(--text-primary)' }}>{formatCurrency(totalIn)}</div>
        </div>
      </div>
    </div>
  );
}

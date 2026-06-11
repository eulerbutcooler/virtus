import Card from '../ui/Card';
import { formatCurrency } from '../../lib/formatters';

export default function PoolCard({ pool, loading }) {
  if (loading) {
    return (
      <Card style={{ marginBottom: 'var(--space-6)' }}>
        <div style={{
          height: '80px',
          background: 'linear-gradient(90deg, var(--bg-surface) 25%, var(--bg-hover) 50%, var(--bg-surface) 75%)',
          backgroundSize: '200% 100%',
          animation: 'shimmer 1.5s linear infinite',
          borderRadius: 'var(--radius-sm)'
        }} />
      </Card>
    );
  }

  const balance  = pool?.balance   || 0;
  const totalIn  = pool?.total_in  || 0;
  const totalOut = pool?.total_out || 0;
  const healthRatio = totalIn > 0 ? Math.min((totalOut / totalIn) * 100, 100) : 0;

  return (
    <Card style={{ marginBottom: 'var(--space-6)' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-end', marginBottom: 'var(--space-5)' }}>
        <div>

          <div style={{
            fontSize: 'var(--font-xs)',
            color: 'var(--text-secondary)',
            textTransform: 'uppercase',
            letterSpacing: '0.06em',
            fontWeight: 'var(--font-weight-bold)',
            marginBottom: 'var(--space-2)'
          }}>
            Community Pool
          </div>

          <div style={{
            fontFamily: 'var(--font-display)',
            fontSize: 'var(--font-2xl)',
            fontWeight: 600,
            letterSpacing: '-0.03em',
            color: 'var(--leaf-200)',
            lineHeight: 1,
            textShadow: '0 0 24px rgba(92, 171, 110, 0.2)'
          }}>
            {formatCurrency(balance)}
          </div>
        </div>


        <div style={{ textAlign: 'right', fontSize: 'var(--font-xs)', color: 'var(--text-secondary)', lineHeight: 1.8 }}>
          <div>
            <span style={{ marginRight: '4px' }}>↑</span>
            <span style={{ color: 'var(--leaf-300)', fontWeight: 'var(--font-weight-medium)' }}>{formatCurrency(totalIn)}</span>
            <span style={{ marginLeft: '4px' }}>in</span>
          </div>
          <div>
            <span style={{ marginRight: '4px' }}>↓</span>
            <span style={{ color: 'var(--solar-300)', fontWeight: 'var(--font-weight-medium)' }}>{formatCurrency(totalOut)}</span>
            <span style={{ marginLeft: '4px' }}>out</span>
          </div>
        </div>
      </div>


      <div style={{
        width: '100%',
        height: '3px',
        background: 'var(--bg-surface)',
        borderRadius: '999px',
        overflow: 'hidden',
        position: 'relative'
      }}>
        <div style={{
          position: 'absolute',
          left: 0, top: 0, bottom: 0,
          width: `${healthRatio}%`,
          background: 'linear-gradient(90deg, var(--leaf-400), var(--leaf-300) 60%, #A8E6B0)',
          backgroundSize: '200% 100%',
          animation: 'leafShimmer 4s linear infinite',
          borderRadius: '999px',
          transition: 'width 800ms ease-out',
          boxShadow: '0 0 8px rgba(92, 171, 110, 0.4)'
        }} />
      </div>

      <div style={{
        marginTop: 'var(--space-2)',
        fontSize: 'var(--font-xs)',
        color: 'var(--text-disabled)',
        letterSpacing: '0.03em'
      }}>
        {healthRatio.toFixed(1)}% distributed
      </div>
    </Card>
  );
}

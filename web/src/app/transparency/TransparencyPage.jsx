import { useFullQueue } from '../../hooks/useQueue';
import PageHeader from '../../components/layout/PageHeader';
import PoolOverview from '../../components/transparency/PoolOverview';
import PublicQueueTable from '../../components/transparency/PublicQueueTable';
import ThemeToggle from '../../components/ui/ThemeToggle.jsx';
import { Link } from 'react-router-dom';
import Button from '../../components/ui/Button';

export default function TransparencyPage() {
  const { queue, loading } = useFullQueue(50, 0);

  return (
    <div style={{ minHeight: '100vh', background: 'var(--bg-primary)' }}>

      <header style={{ borderBottom: '1px solid var(--border-default)', background: 'var(--bg-elevated)', padding: 'var(--space-4) var(--space-6)', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Link
          to="/"
          className="transparency-wordmark"
          style={{
            fontFamily: 'var(--font-display)',
            fontSize: '1.35rem',
            fontWeight: 600,
            letterSpacing: '-0.02em',
            color: 'var(--leaf-300)',
            textDecoration: 'none',
            textShadow: '0 0 20px rgba(92,171,110,0.35)',
            transition: 'color 150ms ease, text-shadow 150ms ease',
          }}
          onMouseEnter={e => { e.currentTarget.style.color = 'var(--leaf-200)'; e.currentTarget.style.textShadow = '0 0 28px rgba(92,171,110,0.55)'; }}
          onMouseLeave={e => { e.currentTarget.style.color = 'var(--leaf-300)'; e.currentTarget.style.textShadow = '0 0 20px rgba(92,171,110,0.35)'; }}
        >
          Virtus
        </Link>
        <div style={{ display: 'flex', gap: 'var(--space-4)', alignItems: 'center' }}>
          <ThemeToggle />
          <Link to="/login"><Button variant="ghost">Sign In</Button></Link>
          <Link to="/register"><Button variant="primary">Join Community</Button></Link>
        </div>
      </header>

      <main style={{ padding: 'var(--space-8) var(--space-6)', maxWidth: '1000px', margin: '0 auto' }}>
        <PageHeader
          title="Community Transparency"
          subtitle="Real-time visibility into the community pool and fulfillment queue."
        />

        <PoolOverview />

        <div>
          <h2 style={{ fontSize: 'var(--font-lg)', color: 'var(--text-primary)', marginBottom: 'var(--space-4)' }}>Live Fulfillment Queue</h2>
          <PublicQueueTable queue={queue} loading={loading} />
        </div>
      </main>
    </div>
  );
}

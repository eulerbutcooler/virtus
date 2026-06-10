import { useFullQueue } from '../../hooks/useQueue';
import PageHeader from '../../components/layout/PageHeader';
import PoolOverview from '../../components/transparency/PoolOverview';
import PublicQueueTable from '../../components/transparency/PublicQueueTable';
import { Link } from 'react-router-dom';
import Button from '../../components/ui/Button';

export default function TransparencyPage() {
  const { queue, loading } = useFullQueue(50, 0);

  return (
    <div style={{ minHeight: '100vh', background: 'var(--bg-primary)' }}>
      {/* Public Topbar */}
      <header style={{ borderBottom: '1px solid var(--border-default)', background: 'var(--bg-elevated)', padding: 'var(--space-4) var(--space-6)', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h1 style={{ fontSize: 'var(--font-lg)', color: 'var(--beige-300)' }}>Virtus</h1>
        <div style={{ display: 'flex', gap: 'var(--space-4)' }}>
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

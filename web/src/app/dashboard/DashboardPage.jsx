import { usePool } from '../../hooks/usePool';
import { useQueue } from '../../hooks/useQueue';
import { useContributions } from '../../hooks/useContributions';
import { useRequests } from '../../hooks/useRequests';
import PoolCard from '../../components/dashboard/PoolCard';
import QueuePositionCard from '../../components/dashboard/QueuePositionCard';
import ContributionSummaryCard from '../../components/dashboard/ContributionSummaryCard';
import RecentActivityFeed from '../../components/dashboard/RecentActivityFeed';

export default function DashboardPage() {
  const { pool, loading: poolLoading } = usePool();
  const { requests, loading: reqLoading } = useRequests();
  const activeRequest = requests?.find(r =>
    r.status === 'pending' || r.status === 'funded' || r.status === 'queued'
  );
  const { queueEntry, loading: queueLoading } = useQueue(activeRequest?.id);
  const { total, contributions, loading: contribLoading } = useContributions(5);

  return (
    <div style={{ padding: 'var(--space-6)', maxWidth: '1000px', margin: '0 auto' }}>

      <div style={{ marginBottom: 'var(--space-6)' }}>
        <h1 style={{
          fontFamily: 'var(--font-display)',
          fontSize: 'var(--font-2xl)',
          fontWeight: 600,
          letterSpacing: '-0.03em',
          color: 'var(--text-primary)',
          lineHeight: 1.15,
          marginBottom: 'var(--space-1)'
        }}>
          Dashboard
        </h1>
        <p style={{
          fontSize: 'var(--font-sm)',
          color: 'var(--text-secondary)'
        }}>
          Your community at a glance
        </p>
      </div>


      <PoolCard pool={pool} loading={poolLoading} />


      <div style={{
        display: 'flex',
        gap: 'var(--space-5)',
        marginBottom: 'var(--space-5)',
        flexWrap: 'wrap'
      }}>
        <div style={{ flex: '1 1 55%', minWidth: '280px' }}>
          <QueuePositionCard queueEntry={queueEntry} loading={reqLoading || queueLoading} />
        </div>
        <div style={{ flex: '1 1 38%', minWidth: '240px' }}>
          <ContributionSummaryCard total={total} loading={contribLoading} />
        </div>
      </div>


      <RecentActivityFeed
        requests={requests}
        contributions={contributions}
        loading={reqLoading || contribLoading}
      />
    </div>
  );
}

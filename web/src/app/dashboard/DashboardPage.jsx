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
  const activeRequest = requests?.find(r => r.status === 'pending' || r.status === 'funded' || r.status === 'queued');
  
  const { queueEntry, loading: queueLoading } = useQueue(activeRequest?.id);
  const { total, contributions, loading: contribLoading } = useContributions(5);

  return (
    <div style={{ padding: 'var(--space-6)', maxWidth: '1000px', margin: '0 auto' }}>
      <h1 style={{ fontSize: 'var(--font-xl)', marginBottom: 'var(--space-6)', color: 'var(--text-primary)' }}>Dashboard</h1>
      
      <PoolCard pool={pool} loading={poolLoading} />

      <div style={{ display: 'flex', gap: 'var(--space-6)', marginBottom: 'var(--space-6)', flexWrap: 'wrap' }}>
        <div style={{ flex: '1 1 55%', minWidth: '300px' }}>
          <QueuePositionCard queueEntry={queueEntry} loading={reqLoading || queueLoading} />
        </div>
        <div style={{ flex: '1 1 40%', minWidth: '250px' }}>
          <ContributionSummaryCard total={total} loading={contribLoading} />
        </div>
      </div>

      <RecentActivityFeed requests={requests} contributions={contributions} loading={reqLoading || contribLoading} />
    </div>
  );
}

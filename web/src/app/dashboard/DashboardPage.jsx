import PageHeader from '../../components/layout/PageHeader.jsx'

export default function DashboardPage() {
  return (
    <div>
      <PageHeader
        title="Dashboard"
        subtitle="Overview of your activity, queue position, and contributions."
      />
      <div style={{ display: 'grid', gap: 'var(--space-5)', gridTemplateColumns: '1fr', maxWidth: '800px' }}>
        <div style={{ padding: 'var(--space-6)', background: 'var(--bg-elevated)', borderRadius: 'var(--radius-lg)', border: '1px solid var(--border-default)' }}>
          <h2 style={{ fontSize: 'var(--font-lg)', color: 'var(--beige-300)', marginBottom: 'var(--space-3)' }}>The Pool</h2>
          <p style={{ color: 'var(--text-secondary)' }}>Pool balance and health metrics will appear here.</p>
        </div>
        <div style={{ display: 'grid', gap: 'var(--space-5)', gridTemplateColumns: '1fr 1fr' }}>
          <div style={{ padding: 'var(--space-6)', background: 'var(--bg-elevated)', borderRadius: 'var(--radius-lg)', border: '1px solid var(--border-default)' }}>
            <h2 style={{ fontSize: 'var(--font-lg)', color: 'var(--beige-300)', marginBottom: 'var(--space-3)' }}>Queue Position</h2>
            <p style={{ color: 'var(--text-secondary)' }}>Your current position in the queue.</p>
          </div>
          <div style={{ padding: 'var(--space-6)', background: 'var(--bg-elevated)', borderRadius: 'var(--radius-lg)', border: '1px solid var(--border-default)' }}>
            <h2 style={{ fontSize: 'var(--font-lg)', color: 'var(--beige-300)', marginBottom: 'var(--space-3)' }}>Contributions</h2>
            <p style={{ color: 'var(--text-secondary)' }}>Your total contributions.</p>
          </div>
        </div>
        <div style={{ padding: 'var(--space-6)', background: 'var(--bg-elevated)', borderRadius: 'var(--radius-lg)', border: '1px solid var(--border-default)' }}>
          <h2 style={{ fontSize: 'var(--font-lg)', color: 'var(--beige-300)', marginBottom: 'var(--space-3)' }}>Recent Activity</h2>
          <p style={{ color: 'var(--text-secondary)' }}>Recent requests and contributions will appear here.</p>
        </div>
      </div>
    </div>
  )
}

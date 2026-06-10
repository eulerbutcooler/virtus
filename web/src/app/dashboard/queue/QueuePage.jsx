import PageHeader from '../../../components/layout/PageHeader.jsx'

export default function QueuePage() {
  return (
    <div>
      <PageHeader
        title="Queue"
        subtitle="View your position in the fulfillment queue."
      />
      <p style={{ color: 'var(--text-secondary)' }}>Queue position and funding progress will appear here.</p>
    </div>
  )
}

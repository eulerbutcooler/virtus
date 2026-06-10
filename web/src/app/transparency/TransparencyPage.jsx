import PageHeader from '../../components/layout/PageHeader.jsx'

export default function TransparencyPage() {
  return (
    <div>
      <PageHeader
        title="Transparency"
        subtitle="Public-facing view of the pool and queue."
      />
      <p style={{ color: 'var(--text-secondary)' }}>Public pool overview and queue will appear here.</p>
    </div>
  )
}

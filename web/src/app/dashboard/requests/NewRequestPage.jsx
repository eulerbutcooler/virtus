import PageHeader from '../../../components/layout/PageHeader.jsx'

export default function NewRequestPage() {
  return (
    <div>
      <PageHeader
        title="New Request"
        subtitle="Submit a new item request to join the queue."
      />
      <p style={{ color: 'var(--text-secondary)' }}>Request form will appear here.</p>
    </div>
  )
}

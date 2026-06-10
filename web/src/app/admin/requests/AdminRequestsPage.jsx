import PageHeader from '../../../components/layout/PageHeader.jsx'

export default function AdminRequestsPage() {
  return (
    <div>
      <PageHeader
        title="Admin Requests"
        subtitle="Review and approve pending item requests."
      />
      <p style={{ color: 'var(--text-secondary)' }}>Admin request management table will appear here.</p>
    </div>
  )
}

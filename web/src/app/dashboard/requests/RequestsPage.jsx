import PageHeader from '../../../components/layout/PageHeader.jsx'
import Button from '../../../components/ui/Button.jsx'

export default function RequestsPage() {
  return (
    <div>
      <PageHeader
        title="My Requests"
        subtitle="Manage your item requests and track their status."
        actions={<Button variant="primary" size="sm">+ New Request</Button>}
      />
      <p style={{ color: 'var(--text-secondary)' }}>Your requests will appear here in a table.</p>
    </div>
  )
}

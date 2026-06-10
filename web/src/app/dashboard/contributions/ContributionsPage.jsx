import PageHeader from '../../../components/layout/PageHeader.jsx'
import Button from '../../../components/ui/Button.jsx'

export default function ContributionsPage() {
  return (
    <div>
      <PageHeader
        title="Contributions"
        subtitle="View your contribution history and make new contributions."
        actions={<Button variant="primary" size="sm">Contribute</Button>}
      />
      <p style={{ color: 'var(--text-secondary)' }}>Your contributions will appear here in a table.</p>
    </div>
  )
}

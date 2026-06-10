import PageHeader from '../../../components/layout/PageHeader.jsx'
import Button from '../../../components/ui/Button.jsx'

export default function ImpactPage() {
  return (
    <div>
      <PageHeader
        title="Impact"
        subtitle="Record and view impact outcomes from fulfilled requests."
        actions={<Button variant="primary" size="sm">Record Impact</Button>}
      />
      <p style={{ color: 'var(--text-secondary)' }}>Impact records will appear here.</p>
    </div>
  )
}

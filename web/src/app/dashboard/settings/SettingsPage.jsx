import PageHeader from '../../../components/layout/PageHeader.jsx'

export default function SettingsPage() {
  return (
    <div>
      <PageHeader
        title="Settings"
        subtitle="Manage your profile and preferences."
      />
      <p style={{ color: 'var(--text-secondary)' }}>Settings form will appear here.</p>
    </div>
  )
}

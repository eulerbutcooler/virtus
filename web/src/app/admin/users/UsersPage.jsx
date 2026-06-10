import PageHeader from '../../../components/layout/PageHeader.jsx'

export default function UsersPage() {
  return (
    <div>
      <PageHeader
        title="Users"
        subtitle="Manage community members and verify accounts."
      />
      <p style={{ color: 'var(--text-secondary)' }}>User management table will appear here.</p>
    </div>
  )
}

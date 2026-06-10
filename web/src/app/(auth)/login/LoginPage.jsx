export default function LoginPage() {
  return (
    <div style={{ minHeight: '100vh', display: 'flex', alignItems: 'center', justifyContent: 'center', background: 'var(--bg-primary)' }}>
      <div style={{ width: '100%', maxWidth: '480px', padding: 'var(--space-8)' }}>
        <h1 style={{ fontSize: 'var(--font-2xl)', color: 'var(--beige-300)', marginBottom: 'var(--space-6)', textAlign: 'center' }}>
          Virtus
        </h1>
        <div style={{ background: 'var(--bg-elevated)', borderRadius: 'var(--radius-lg)', border: '1px solid var(--border-default)', padding: 'var(--space-6)' }}>
          <h2 style={{ fontSize: 'var(--font-lg)', marginBottom: 'var(--space-4)', textAlign: 'center' }}>Sign In</h2>
          <p style={{ color: 'var(--text-secondary)', textAlign: 'center' }}>Login form will appear here.</p>
        </div>
      </div>
    </div>
  )
}

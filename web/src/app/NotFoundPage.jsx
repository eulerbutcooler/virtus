import PageHeader from '../../components/layout/PageHeader';
import Button from '../../components/ui/Button';
import { Link } from 'react-router-dom';

export default function NotFoundPage() {
  return (
    <div style={{ padding: 'var(--space-8)', maxWidth: '800px', margin: '0 auto', textAlign: 'center' }}>
      <h1 style={{ fontFamily: 'var(--font-display)', fontSize: 'var(--font-3xl)', color: 'var(--leaf-300)', marginBottom: 'var(--space-4)', letterSpacing: '-0.04em', textShadow: '0 0 24px rgba(92,171,110,0.3)' }}>404</h1>
      <PageHeader title="Page Not Found" subtitle="The page you are looking for doesn't exist or has been moved." />
      <div style={{ marginTop: 'var(--space-6)' }}>
        <Link to="/dashboard"><Button variant="primary">Return Home</Button></Link>
      </div>
    </div>
  );
}

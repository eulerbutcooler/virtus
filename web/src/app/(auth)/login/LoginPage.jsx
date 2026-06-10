import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../../../hooks/useAuth.jsx';
import Button from '../../../components/ui/Button.jsx';
import Input from '../../../components/ui/Input.jsx';

export default function LoginPage() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const { login } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      await login(email, password);
      navigate('/dashboard');
    } catch (err) {
      setError(err.message || 'Failed to login');
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <h2 className="auth-form-title">Sign In</h2>
      <p className="auth-form-subtitle">Welcome back to Virtus</p>
      
      {error && <div style={{ color: 'var(--error)', marginBottom: 'var(--space-4)', fontSize: 'var(--font-sm)' }}>{error}</div>}

      <form className="auth-form" onSubmit={handleSubmit}>
        <Input 
          type="email" 
          placeholder="Email address" 
          value={email} 
          onChange={(e) => setEmail(e.target.value)} 
          required 
        />
        <Input 
          type="password" 
          placeholder="Password" 
          value={password} 
          onChange={(e) => setPassword(e.target.value)} 
          required 
        />
        <Button type="submit" variant="primary" disabled={loading} style={{ marginTop: 'var(--space-2)' }}>
          {loading ? 'Signing in...' : 'Sign In'}
        </Button>
      </form>
      
      <p style={{ marginTop: 'var(--space-6)', fontSize: 'var(--font-sm)', color: 'var(--text-secondary)' }}>
        Don't have an account? <Link to="/register">Create one</Link>
      </p>
    </>
  );
}

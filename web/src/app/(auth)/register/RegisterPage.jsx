import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../../../hooks/useAuth.jsx';
import Button from '../../../components/ui/Button.jsx';
import Input from '../../../components/ui/Input.jsx';

export default function RegisterPage() {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const { register } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      await register(name, email, password);
      navigate('/dashboard');
    } catch (err) {
      setError(err.message || 'Failed to register');
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <h2 className="auth-form-title">Create Account</h2>
      <p className="auth-form-subtitle">Join the Virtus community</p>
      
      {error && <div style={{ color: 'var(--error)', marginBottom: 'var(--space-4)', fontSize: 'var(--font-sm)' }}>{error}</div>}

      <form className="auth-form" onSubmit={handleSubmit}>
        <Input 
          type="text" 
          placeholder="Full name" 
          value={name} 
          onChange={(e) => setName(e.target.value)} 
          required 
        />
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
          {loading ? 'Creating account...' : 'Create Account'}
        </Button>
      </form>
      
      <p style={{ marginTop: 'var(--space-6)', fontSize: 'var(--font-sm)', color: 'var(--text-secondary)' }}>
        Already a member? <Link to="/login">Sign in</Link>
      </p>
    </>
  );
}

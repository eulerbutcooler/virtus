import { useState } from 'react';
import { api } from '../../../lib/api';
import PageHeader from '../../../components/layout/PageHeader';
import Card from '../../../components/ui/Card';
import Input from '../../../components/ui/Input';
import Button from '../../../components/ui/Button';

export default function SettingsPage() {
  const [profileLoading, setProfileLoading] = useState(false);
  const [pwdLoading, setPwdLoading] = useState(false);
  const [toast, setToast] = useState(null);

  const [name, setName] = useState('');
  const [pwdForm, setPwdForm] = useState({ current_password: '', new_password: '' });

  const showToast = (msg, type = 'success') => {
    setToast({ msg, type });
    setTimeout(() => setToast(null), 3000);
  };

  const handleUpdateProfile = async (e) => {
    e.preventDefault();
    setProfileLoading(true);
    try {
      await api.patch('/me', { name });
      showToast('Profile updated successfully');
    } catch (err) {
      showToast(err.message || 'Failed to update profile', 'error');
    } finally {
      setProfileLoading(false);
    }
  };

  const handleUpdatePassword = async (e) => {
    e.preventDefault();
    setPwdLoading(true);
    try {
      await api.post('/me/password', pwdForm);
      showToast('Password changed successfully');
      setPwdForm({ current_password: '', new_password: '' });
    } catch (err) {
      showToast(err.message || 'Failed to change password', 'error');
    } finally {
      setPwdLoading(false);
    }
  };

  return (
    <div style={{ padding: 'var(--space-6)', maxWidth: '800px', margin: '0 auto', position: 'relative' }}>
      <PageHeader title="Settings" subtitle="Manage your account preferences" />

      {toast && (
        <div style={{
          position: 'fixed', bottom: 'var(--space-6)', right: 'var(--space-6)',
          background: toast.type === 'error' ? 'var(--error)' : 'var(--success)',
          color: '#fff', padding: 'var(--space-3) var(--space-4)', borderRadius: 'var(--radius-md)',
          boxShadow: '0 4px 12px rgba(0,0,0,0.1)', zIndex: 100
        }}>
          {toast.msg}
        </div>
      )}

      <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-8)' }}>
        <Card>
          <h2 style={{ fontSize: 'var(--font-lg)', color: 'var(--text-primary)', marginBottom: 'var(--space-4)' }}>Profile Settings</h2>
          <form onSubmit={handleUpdateProfile} style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-4)' }}>
            <Input placeholder="Full Name" value={name} onChange={e => setName(e.target.value)} required />
            <div>
              <Button type="submit" variant="primary" disabled={profileLoading}>Save Changes</Button>
            </div>
          </form>
        </Card>

        <Card>
          <h2 style={{ fontSize: 'var(--font-lg)', color: 'var(--text-primary)', marginBottom: 'var(--space-4)' }}>Security</h2>
          <form onSubmit={handleUpdatePassword} style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-4)' }}>
            <Input type="password" placeholder="Current Password" value={pwdForm.current_password} onChange={e => setPwdForm(p => ({ ...p, current_password: e.target.value }))} required />
            <Input type="password" placeholder="New Password" value={pwdForm.new_password} onChange={e => setPwdForm(p => ({ ...p, new_password: e.target.value }))} required />
            <div>
              <Button type="submit" variant="primary" disabled={pwdLoading}>Change Password</Button>
            </div>
          </form>
        </Card>
      </div>
    </div>
  );
}

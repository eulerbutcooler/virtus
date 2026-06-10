import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useRequests } from '../../../hooks/useRequests';
import RequestForm from '../../../components/forms/RequestForm';
import PageHeader from '../../../components/layout/PageHeader';
import Card from '../../../components/ui/Card';

export default function NewRequestPage() {
  const navigate = useNavigate();
  const { createRequest } = useRequests();
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (formData) => {
    try {
      setLoading(true);
      setError(null);
      await createRequest(formData);
      navigate('/dashboard/requests');
    } catch (err) {
      setError(err.message || 'Failed to create request');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ padding: 'var(--space-6)', maxWidth: '800px', margin: '0 auto' }}>
      <PageHeader 
        title="New Request" 
        subtitle="Submit a new fulfillment request to the community" 
      />
      
      <Card>
        {error && <div style={{ color: 'var(--error)', marginBottom: 'var(--space-4)' }}>{error}</div>}
        <RequestForm onSubmit={handleSubmit} isLoading={loading} />
      </Card>
    </div>
  );
}

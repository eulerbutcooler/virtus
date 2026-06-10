import { useState } from 'react';
import Input from '../ui/Input';
import Button from '../ui/Button';

export default function ImpactForm({ deliveries, onSubmit, onCancel, isLoading }) {
  const [formData, setFormData] = useState({
    delivery_id: '',
    interval: '30_days',
    outcome: '',
    satisfaction_score: 5,
    description: ''
  });

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: name === 'satisfaction_score' ? Number(value) : value }));
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    onSubmit(formData);
  };

  return (
    <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-4)' }}>
      <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-2)' }}>
        <label style={{ fontSize: 'var(--font-sm)', color: 'var(--text-secondary)' }}>Delivered Item</label>
        <select 
          name="delivery_id" 
          value={formData.delivery_id} 
          onChange={handleChange}
          required
          style={{ padding: 'var(--space-2)', background: 'var(--bg-surface)', border: '1px solid var(--border-default)', color: 'var(--text-primary)', borderRadius: 'var(--radius-sm)' }}
        >
          <option value="">Select a delivered item...</option>
          {deliveries.map(d => (
            <option key={d.id} value={d.id}>
              {d.fulfillment?.request?.item_name || `Delivery ${d.id.substring(0, 8)}`}
            </option>
          ))}
        </select>
      </div>

      <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-2)' }}>
        <label style={{ fontSize: 'var(--font-sm)', color: 'var(--text-secondary)' }}>Time Interval</label>
        <select 
          name="interval" 
          value={formData.interval} 
          onChange={handleChange}
          style={{ padding: 'var(--space-2)', background: 'var(--bg-surface)', border: '1px solid var(--border-default)', color: 'var(--text-primary)', borderRadius: 'var(--radius-sm)' }}
        >
          <option value="immediate">Immediate</option>
          <option value="30_days">30 Days</option>
          <option value="90_days">90 Days</option>
          <option value="180_days">180+ Days</option>
        </select>
      </div>

      <Input name="outcome" placeholder="Primary Outcome (e.g. Returned to work)" value={formData.outcome} onChange={handleChange} required />

      <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-2)' }}>
        <label style={{ fontSize: 'var(--font-sm)', color: 'var(--text-secondary)' }}>Satisfaction Score ({formData.satisfaction_score}/5)</label>
        <input 
          type="range" 
          name="satisfaction_score" 
          min="1" max="5" 
          value={formData.satisfaction_score} 
          onChange={handleChange} 
        />
      </div>

      <textarea 
        name="description" 
        placeholder="How has this item impacted your situation?" 
        value={formData.description} 
        onChange={handleChange}
        style={{ padding: 'var(--space-2)', background: 'var(--bg-surface)', border: '1px solid var(--border-default)', color: 'var(--text-primary)', borderRadius: 'var(--radius-sm)', minHeight: '80px', fontFamily: 'inherit' }}
        required
      />

      <div style={{ display: 'flex', gap: 'var(--space-4)', marginTop: 'var(--space-4)' }}>
        <Button type="submit" variant="primary" disabled={isLoading}>
          {isLoading ? 'Saving...' : 'Record Impact'}
        </Button>
        <Button type="button" variant="ghost" onClick={onCancel} disabled={isLoading}>
          Cancel
        </Button>
      </div>
    </form>
  );
}

import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import Input from '../ui/Input';
import Button from '../ui/Button';

const CATEGORIES = ['housing', 'medical', 'transportation', 'education', 'food', 'utilities', 'other'];

export default function RequestForm({ initialData = null, onSubmit, isLoading = false }) {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    item_name: '',
    item_category: 'housing',
    description: '',
    urgency: 'standard',
    estimated_cost: '',
    justification: ''
  });

  useEffect(() => {
    if (initialData) {
      setFormData({
        item_name: initialData.item_name || '',
        item_category: initialData.item_category || 'housing',
        description: initialData.description || '',
        urgency: initialData.urgency || 'standard',
        estimated_cost: initialData.estimated_cost || '',
        justification: initialData.justification || ''
      });
    }
  }, [initialData]);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: name === 'estimated_cost' ? Number(value) : value }));
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    onSubmit(formData);
  };

  return (
    <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-4)' }}>
      <Input
        name="item_name"
        placeholder="Item Name (e.g. Winter Coat, Medical Bill)"
        value={formData.item_name}
        onChange={handleChange}
        required
      />

      <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-2)' }}>
        <label style={{ fontSize: 'var(--font-sm)', color: 'var(--text-secondary)' }}>Category</label>
        <select 
          name="item_category" 
          value={formData.item_category} 
          onChange={handleChange}
          style={{ 
            height: '36px', 
            background: 'var(--bg-surface)', 
            border: '1px solid var(--border-default)', 
            borderRadius: 'var(--radius-sm)', 
            color: 'var(--text-primary)',
            padding: '0 var(--space-2)'
          }}
          required
        >
          {CATEGORIES.map(cat => (
            <option key={cat} value={cat}>{cat.charAt(0).toUpperCase() + cat.slice(1)}</option>
          ))}
        </select>
      </div>

      <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-2)' }}>
        <label style={{ fontSize: 'var(--font-sm)', color: 'var(--text-secondary)' }}>Urgency</label>
        <div style={{ display: 'flex', gap: 'var(--space-4)' }}>
          {['low', 'standard', 'high', 'critical'].map(level => (
            <label key={level} style={{ display: 'flex', alignItems: 'center', gap: 'var(--space-2)', fontSize: 'var(--font-sm)', color: 'var(--text-primary)' }}>
              <input 
                type="radio" 
                name="urgency" 
                value={level} 
                checked={formData.urgency === level} 
                onChange={handleChange} 
              />
              {level.charAt(0).toUpperCase() + level.slice(1)}
            </label>
          ))}
        </div>
      </div>

      <Input
        name="estimated_cost"
        type="number"
        placeholder="Estimated Cost ($)"
        value={formData.estimated_cost}
        onChange={handleChange}
        min="1"
        required
      />

      <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-2)' }}>
        <label style={{ fontSize: 'var(--font-sm)', color: 'var(--text-secondary)' }}>Description</label>
        <textarea
          name="description"
          placeholder="Detailed description of what you need..."
          value={formData.description}
          onChange={handleChange}
          required
          style={{
            background: 'var(--bg-surface)',
            border: '1px solid var(--border-default)',
            borderRadius: 'var(--radius-sm)',
            color: 'var(--text-primary)',
            padding: 'var(--space-2)',
            minHeight: '80px',
            fontFamily: 'inherit'
          }}
        />
      </div>

      <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-2)' }}>
        <label style={{ fontSize: 'var(--font-sm)', color: 'var(--text-secondary)' }}>Justification</label>
        <textarea
          name="justification"
          placeholder="Why is this needed? How will it help?"
          value={formData.justification}
          onChange={handleChange}
          required
          style={{
            background: 'var(--bg-surface)',
            border: '1px solid var(--border-default)',
            borderRadius: 'var(--radius-sm)',
            color: 'var(--text-primary)',
            padding: 'var(--space-2)',
            minHeight: '80px',
            fontFamily: 'inherit'
          }}
        />
      </div>

      <div style={{ display: 'flex', gap: 'var(--space-4)', marginTop: 'var(--space-4)' }}>
        <Button type="submit" variant="primary" disabled={isLoading}>
          {isLoading ? 'Saving...' : initialData ? 'Save Changes' : 'Submit Request'}
        </Button>
        <Button type="button" variant="ghost" onClick={() => navigate('/dashboard/requests')} disabled={isLoading}>
          Cancel
        </Button>
      </div>
    </form>
  );
}

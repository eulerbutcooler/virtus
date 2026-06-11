import { useState, useEffect } from 'react';
import { useStripe, useElements, PaymentElement } from '@stripe/react-stripe-js';
import Button from '../ui/Button';

export default function ContributionForm({ amount, clientSecret, onSuccess, onCancel }) {
  const stripe = useStripe();
  const elements = useElements();
  const [error, setError] = useState(null);
  const [isProcessing, setIsProcessing] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!stripe || !elements) {
      return; // Stripe.js hasn't loaded yet
    }

    setIsProcessing(true);
    setError(null);

    const { error: submitError } = await elements.submit();
    if (submitError) {
      setError(submitError.message);
      setIsProcessing(false);
      return;
    }

    const { error: confirmError } = await stripe.confirmPayment({
      elements,
      clientSecret,
      confirmParams: {

        return_url: window.location.origin + '/dashboard/contributions?success=true',
      },
      redirect: 'if_required', // Avoid full page redirect if possible
    });

    if (confirmError) {
      setError(confirmError.message);
      setIsProcessing(false);
    } else {

      onSuccess();
    }
  };

  return (
    <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: 'var(--space-4)' }}>
      <PaymentElement />

      {error && <div style={{ color: 'var(--error)', fontSize: 'var(--font-sm)' }}>{error}</div>}

      <div style={{ display: 'flex', gap: 'var(--space-4)', marginTop: 'var(--space-4)' }}>
        <Button type="submit" variant="primary" disabled={!stripe || isProcessing}>
          {isProcessing ? 'Processing...' : `Contribute $${amount}`}
        </Button>
        <Button type="button" variant="ghost" onClick={onCancel} disabled={isProcessing}>
          Cancel
        </Button>
      </div>
    </form>
  );
}

import { Outlet } from 'react-router-dom';
import './AuthLayout.css';

export default function AuthLayout() {
  return (
    <div className="auth-layout">
      <div className="auth-form-side">
        <h1 className="auth-logo">Virtus</h1>
        <div className="auth-form-container">
          <Outlet />
        </div>
      </div>
      <div className="auth-context-side">
        <div className="auth-stat-card">
          <div className="auth-stat-number">412</div>
          <div className="auth-stat-label">items fulfilled this month</div>
        </div>
        <div className="auth-attribution">A transparent community fulfillment ecosystem</div>
      </div>
    </div>
  );
}

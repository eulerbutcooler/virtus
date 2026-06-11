import { Link } from 'react-router-dom';
import { Outlet } from 'react-router-dom';
import './AuthLayout.css';

export default function AuthLayout() {
  return (
    <div className="auth-layout">
      <div className="auth-form-side">

        <div className="auth-form-header">
          <Link to="/" className="auth-logo">Virtus</Link>
        </div>
        <div className="auth-form-container">
          <Outlet />
        </div>
      </div>
      <div className="auth-context-side">
        <div className="auth-stat-card">
          <div className="auth-stat-number">412</div>
          <div className="auth-stat-divider" />
          <div className="auth-stat-label">
            items fulfilled this month by community members pooling together
          </div>
        </div>
        <div className="auth-attribution">A transparent community fulfillment ecosystem</div>
      </div>
    </div>
  );
}

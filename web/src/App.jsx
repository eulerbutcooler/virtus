import { BrowserRouter, Routes, Route, Navigate, Outlet } from 'react-router-dom'
import { AuthProvider, useAuth } from './hooks/useAuth.jsx'
import AppShell from './components/layout/AppShell.jsx'
import AuthLayout from './app/(auth)/AuthLayout.jsx'
import LoginPage from './app/(auth)/login/LoginPage.jsx'
import RegisterPage from './app/(auth)/register/RegisterPage.jsx'
import DashboardPage from './app/dashboard/DashboardPage.jsx'
import RequestsPage from './app/dashboard/requests/RequestsPage.jsx'
import NewRequestPage from './app/dashboard/requests/NewRequestPage.jsx'
import ContributionsPage from './app/dashboard/contributions/ContributionsPage.jsx'
import QueuePage from './app/dashboard/queue/QueuePage.jsx'
import ImpactPage from './app/dashboard/impact/ImpactPage.jsx'
import SettingsPage from './app/dashboard/settings/SettingsPage.jsx'
import UsersPage from './app/admin/users/UsersPage.jsx'
import AdminRequestsPage from './app/admin/requests/AdminRequestsPage.jsx'
import FulfillmentsPage from './app/admin/fulfillments/FulfillmentsPage.jsx'
import TransparencyPage from './app/transparency/TransparencyPage.jsx'
import InstitutionDashPage from './app/institutions/dashboard/InstitutionDashPage.jsx'

function ProtectedRoute({ children, requiredRole }) {
  const { isAuthenticated, loading, user } = useAuth();
  
  if (loading) {
    return <div style={{ minHeight: '100vh', display: 'flex', alignItems: 'center', justifyContent: 'center', background: 'var(--bg-primary)', color: 'var(--text-primary)' }}>Loading...</div>;
  }
  
  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  if (requiredRole && user?.role !== requiredRole) {
    return <Navigate to="/dashboard" replace />;
  }
  
  return children || <Outlet />;
}

function AdminRoute({ children }) {
  return <ProtectedRoute requiredRole="admin">{children || <Outlet />}</ProtectedRoute>;
}

function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          {/* Public routes */}
          <Route element={<AuthLayout />}>
            <Route path="/login" element={<LoginPage />} />
            <Route path="/register" element={<RegisterPage />} />
          </Route>
          
          <Route path="/transparency" element={<TransparencyPage />} />

          {/* Authenticated routes inside AppShell */}
          <Route element={<AppShell />}>
            <Route element={<ProtectedRoute />}>
              {/* Member routes */}
              <Route path="/dashboard" element={<DashboardPage />} />
              <Route path="/dashboard/requests" element={<RequestsPage />} />
              <Route path="/dashboard/requests/new" element={<NewRequestPage />} />
              <Route path="/dashboard/contributions" element={<ContributionsPage />} />
              <Route path="/dashboard/queue" element={<QueuePage />} />
              <Route path="/dashboard/impact" element={<ImpactPage />} />
              <Route path="/dashboard/settings" element={<SettingsPage />} />
            </Route>

            {/* Admin routes */}
            <Route element={<AdminRoute />}>
              <Route path="/admin/users" element={<UsersPage />} />
              <Route path="/admin/requests" element={<AdminRequestsPage />} />
              <Route path="/admin/fulfillments" element={<FulfillmentsPage />} />
            </Route>
          </Route>

          {/* Institution dashboard */}
          <Route path="/institutions/dashboard" element={<AppShell />}>
            <Route element={<ProtectedRoute requiredRole="institution" />}>
              <Route index element={<InstitutionDashPage />} />
            </Route>
          </Route>

          {/* Root redirect */}
          <Route path="/" element={<Navigate to="/dashboard" replace />} />
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  )
}

export default App

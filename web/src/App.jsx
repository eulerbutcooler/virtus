import { BrowserRouter, Routes, Route, Navigate, Outlet } from 'react-router-dom'
import AppShell from './components/layout/AppShell.jsx'
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

// Simple protected route wrapper (no auth yet, just redirects)
function ProtectedRoute({ children }) {
  // TODO: check auth state
  return children || <Outlet />
}

function AdminRoute({ children }) {
  // TODO: check admin role
  return children || <Outlet />
}

function App() {
  return (
    <BrowserRouter>
      <Routes>
        {/* Public routes */}
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
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

        {/* Institution dashboard placeholder */}
        <Route path="/institutions/dashboard" element={<AppShell />}>
          <Route index element={<div style={{ padding: 'var(--space-8)' }}><h1>Institution Dashboard</h1><p>Placeholder for institution dashboard.</p></div>} />
        </Route>

        {/* Root redirect */}
        <Route path="/" element={<Navigate to="/dashboard" replace />} />
      </Routes>
    </BrowserRouter>
  )
}

export default App

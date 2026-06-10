import { lazy, Suspense } from 'react'
import { BrowserRouter, Routes, Route, Navigate, Outlet } from 'react-router-dom'
import { AuthProvider, useAuth } from './hooks/useAuth.jsx'
import AppShell from './components/layout/AppShell.jsx'
import AuthLayout from './app/(auth)/AuthLayout.jsx'
import ErrorBoundary from './components/layout/ErrorBoundary.jsx'

// Lazy loading components
const LoginPage = lazy(() => import('./app/(auth)/login/LoginPage.jsx'))
const RegisterPage = lazy(() => import('./app/(auth)/register/RegisterPage.jsx'))
const DashboardPage = lazy(() => import('./app/dashboard/DashboardPage.jsx'))
const RequestsPage = lazy(() => import('./app/dashboard/requests/RequestsPage.jsx'))
const NewRequestPage = lazy(() => import('./app/dashboard/requests/NewRequestPage.jsx'))
const ContributionsPage = lazy(() => import('./app/dashboard/contributions/ContributionsPage.jsx'))
const QueuePage = lazy(() => import('./app/dashboard/queue/QueuePage.jsx'))
const ImpactPage = lazy(() => import('./app/dashboard/impact/ImpactPage.jsx'))
const SettingsPage = lazy(() => import('./app/dashboard/settings/SettingsPage.jsx'))
const UsersPage = lazy(() => import('./app/admin/users/UsersPage.jsx'))
const AdminRequestsPage = lazy(() => import('./app/admin/requests/AdminRequestsPage.jsx'))
const FulfillmentsPage = lazy(() => import('./app/admin/fulfillments/FulfillmentsPage.jsx'))
const TransparencyPage = lazy(() => import('./app/transparency/TransparencyPage.jsx'))
const InstitutionDashPage = lazy(() => import('./app/institutions/dashboard/InstitutionDashPage.jsx'))
const NotFoundPage = lazy(() => import('./app/NotFoundPage.jsx'))

const PageLoader = () => (
  <div style={{ height: 'calc(100vh - 100px)', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
    <div style={{ width: '40px', height: '40px', borderRadius: '50%', border: '2px solid var(--border-default)', borderTopColor: 'var(--beige-200)', animation: 'spin 1s linear infinite' }}></div>
  </div>
);

function ProtectedRoute({ children, requiredRole }) {
  const { isAuthenticated, loading, user } = useAuth();
  
  if (loading) {
    return <PageLoader />;
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
    <ErrorBoundary>
      <AuthProvider>
        <BrowserRouter>
          <Suspense fallback={<PageLoader />}>
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
              
              {/* 404 */}
              <Route path="*" element={<NotFoundPage />} />
            </Routes>
          </Suspense>
        </BrowserRouter>
      </AuthProvider>
    </ErrorBoundary>
  )
}

export default App

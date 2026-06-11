import { useState, useEffect } from 'react'
import { Outlet, useLocation } from 'react-router-dom'
import Sidebar from './Sidebar.jsx'
import Topbar from './Topbar.jsx'
import { useAuth } from '../../hooks/useAuth.jsx'
import './AppShell.css'

export default function AppShell() {
  const [collapsed, setCollapsed] = useState(false)
  const [mobileOpen, setMobileOpen] = useState(false)
  const [isMobile, setIsMobile] = useState(false)
  const location = useLocation()

  useEffect(() => {
    const handleResize = () => {
      const width = window.innerWidth
      setIsMobile(width <= 768)
      if (width >= 769 && width <= 1023) {
        setCollapsed(true)
      } else if (width > 1023) {
        setCollapsed(false)
      }
    }
    handleResize()
    window.addEventListener('resize', handleResize)
    return () => window.removeEventListener('resize', handleResize)
  }, [])

  useEffect(() => {
    setMobileOpen(false)
  }, [location.pathname])

  const handleToggle = () => {
    if (isMobile) {
      setMobileOpen(prev => !prev)
    } else {
      setCollapsed(prev => !prev)
    }
  }

  const { user, logout } = useAuth()

  const getPageTitle = () => {
    const titles = {
      '/dashboard': 'Dashboard',
      '/dashboard/requests': 'My Requests',
      '/dashboard/requests/new': 'New Request',
      '/dashboard/contributions': 'Contributions',
      '/dashboard/queue': 'Queue',
      '/dashboard/impact': 'Impact',
      '/dashboard/settings': 'Settings',
      '/admin/users': 'Users',
      '/admin/requests': 'Admin Requests',
      '/admin/fulfillments': 'Fulfillments',
      '/institutions/dashboard': 'Institution Dashboard',
      '/transparency': 'Transparency',
    }
    return titles[location.pathname] || ''
  }

  return (
    <div className="appshell">
      <Sidebar
        collapsed={collapsed}
        onToggle={handleToggle}
        mobileOpen={mobileOpen}
        onMobileClose={() => setMobileOpen(false)}
        isMobile={isMobile}
      />
      <Topbar
        pageTitle={getPageTitle()}
        onMenuToggle={handleToggle}
        user={user}
        onLogout={logout}
        collapsed={collapsed}
      />
      <div className={[
        'appshell__content',
        collapsed ? 'appshell__content--collapsed' : '',
        isMobile ? 'appshell__content--mobile' : '',
      ].join(' ')}>
        <main className="appshell__main">
          <Outlet />
        </main>
      </div>
    </div>
  )
}
